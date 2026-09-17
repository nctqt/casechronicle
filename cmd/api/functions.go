package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nctqt/casechronicle/internal/database"
	"github.com/nctqt/casechronicle/internal/jsonhelp"
	"github.com/nctqt/casechronicle/internal/openrouter"
	"golang.org/x/crypto/bcrypt"
)

// enforce a category enum
type Category string

const (
	CategoryPodcast         Category = "podcast"
	CategoryNewsBroadcast   Category = "news broadcast"
	CategoryGeneral         Category = "general"
	CategoryInterview       Category = "interview"
	CategoryCourtroom       Category = "courtroom"
	CategoryCreatorAnalysis Category = "creator analysis"
	CategoryPoliceBodycam   Category = "police / bodycam"
)

func (c Category) isValid() bool {
	switch c {
	case CategoryPodcast, CategoryNewsBroadcast, CategoryGeneral, CategoryInterview, CategoryCourtroom, CategoryCreatorAnalysis, CategoryPoliceBodycam:
		return true
	}
	return false
}

func (cfg *apiConfig) processVideoEnrichment(ctx context.Context, videoID uuid.UUID) (err error) {
	// create an independent context that outlives short test/request timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// immediate update video status to analyzing
	err = cfg.queries.UpdateVideoStatus(ctx, database.UpdateVideoStatusParams{
		ID:        videoID,
		Status:    "analyzing",
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("[Worker] Failed to set status='analyzing' for video %s: %v", videoID, err)
	}

	// if there is an error during enrichment, mark video as failed
	defer func() {
		if err != nil {
			log.Printf("[Worker] Marking video %s as failed due to error: %v", videoID, err)
			failErr := cfg.queries.UpdateVideoStatus(ctx, database.UpdateVideoStatusParams{
				ID:        videoID,
				Status:    "failed",
				UpdatedAt: time.Now().UTC(),
			})
			if failErr != nil {
				log.Printf("[Worker] Failed to set status='failed' for video %s: %v", videoID, failErr)
			}
		}
	}()

	// fetch video object from db
	video, err := cfg.queries.GetVideoByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("get video by id: %w", err)
	}

	// resolve transcript (use DB transcript if present, else fetch live via transcript client)
	transcriptText := ""
	if video.RawTranscript != nil && *video.RawTranscript != "" {
		transcriptText = *video.RawTranscript
	} else if video.YoutubeVideoID != "" {
		fetchedText, fetchErr := cfg.transcript.FetchTranscript(ctx, video.YoutubeVideoID)
		if fetchErr != nil {
			log.Printf("[Worker] Could not fetch transcript for video %s (%s): %v. Falling back to metadata.", videoID, video.YoutubeVideoID, fetchErr)
		} else {
			transcriptText = fetchedText

			// persist raw transcript to DB for future runs
			updateErr := cfg.queries.UpdateVideoTranscript(ctx, database.UpdateVideoTranscriptParams{
				ID:            video.ID,
				RawTranscript: &transcriptText,
				UpdatedAt:     time.Now().UTC(),
			})
			if updateErr != nil {
				log.Printf("[Worker] Warning: failed to save fetched transcript to DB for video %s: %v", videoID, updateErr)
			}
		}
	}

	// cap transcript length to avoid blowing through token limits on long videos
	const maxTranscriptChars = 15000
	trimmedTranscript := transcriptText
	if len(trimmedTranscript) > maxTranscriptChars {
		log.Printf("[Worker] Truncating transcript for video %s from %d to %d chars", videoID, len(trimmedTranscript), maxTranscriptChars)
		trimmedTranscript = trimmedTranscript[:maxTranscriptChars] + "\n\n[Transcript truncated due to length...]"
	}

	// prepare OpenRouter request
	analysisReq := openrouter.AnalysisRequest{
		Title:         video.Title,
		ChannelName:   video.ChannelName,
		RawTranscript: trimmedTranscript,
	}

	// call OpenRouter client
	aiResult, err := cfg.openrouter.AnalyzeVideo(ctx, analysisReq)
	if err != nil {
		return fmt.Errorf("openrouter analyze failed: %w", err)
	}

	// parse estimated event date
	var estimatedEventDate sql.NullTime
	if aiResult.EstimatedEventDate != "" {
		parsedDate, err := time.Parse("2006-01-02", aiResult.EstimatedEventDate)
		if err == nil {
			estimatedEventDate = sql.NullTime{Time: parsedDate, Valid: true}
		}
	}

	// set category
	category := Category(aiResult.Category)
	if !category.isValid() {
		category = CategoryGeneral
	}

	err = cfg.queries.UpdateVideoCategory(ctx, database.UpdateVideoCategoryParams{
		ID:       videoID,
		Category: string(category),
	})
	if err != nil {
		return fmt.Errorf("failed to update video category: %w", err)
	}

	// update video data in db
	now := time.Now().UTC()
	_, err = cfg.queries.UpdateVideoSummary(ctx, database.UpdateVideoSummaryParams{
		ID:                 video.ID,
		AiSummary:          &aiResult.Summary,
		EstimatedEventDate: estimatedEventDate,
		SummarySource:      aiResult.SummarySource,
		Status:             "analyzed",
		UpdatedAt:          now,
	})
	if err != nil {
		return fmt.Errorf("failed to update video analysis: %w", err)
	}

	log.Printf("[Worker] Successfully enriched video %s using source: %s", video.ID, aiResult.SummarySource)
	return nil
}

func (cfg *apiConfig) updatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not update password", err)
	}

	id := r.PathValue("user_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	if len(req.Password) < 8 {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long", nil)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	now := time.Now().UTC()
	err = cfg.queries.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		ID:             userUUID,
		HashedPassword: string(hashedPassword),
		UpdatedAt:      now,
	})

	w.WriteHeader(http.StatusNoContent)
}
