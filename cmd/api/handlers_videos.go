package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nctqt/casechronicle/internal/database"
	"github.com/nctqt/casechronicle/internal/jsonhelp"
	"github.com/nctqt/casechronicle/internal/worker"
)

type AddVideoRequest struct {
	URL         string     `json:"url"`
	MilestoneID *uuid.UUID `json:"milestone_id,omitempty"`
}

func (cfg *apiConfig) handlerAddVideo(w http.ResponseWriter, r *http.Request) {
	var req AddVideoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "URL is required", errors.New("missing url"))
		return
	}

	// map pointer to uuid.NullUUID for database insertion
	var milestoneID uuid.NullUUID
	if req.MilestoneID != nil {
		milestoneID = uuid.NullUUID{
			UUID:  *req.MilestoneID,
			Valid: true,
		}
	}

	// we have a url, get video meta data
	ytMeta, err := cfg.yt.FetchVideoMetaData(r.Context(), req.URL)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Failed to fetch yt metadata", err)
		return
	}

	var eventDate = sql.NullTime{
		Time:  ytMeta.PublishedAt,
		Valid: true,
	}

	// push to db
	now := time.Now().UTC()
	newVideo, err := cfg.queries.AddVideo(r.Context(), database.AddVideoParams{
		ID:                 uuid.New(),
		MilestoneID:        milestoneID,
		YoutubeVideoID:     ytMeta.ID,
		Title:              ytMeta.Title,
		ChannelName:        ytMeta.ChannelName,
		Description:        &ytMeta.Description,
		PublishedAt:        ytMeta.PublishedAt,
		CreatedAt:          now,
		UpdatedAt:          now,
		Category:           "uncategorized",
		Status:             "pending",
		EstimatedEventDate: eventDate,
	})
	if err != nil {
		// check if the error is a pgx unique constraint violation (SQLSTATE 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			jsonhelp.RespondWithError(w, http.StatusConflict, "A video with this YouTube ID already exists", err)
			return
		}

		log.Printf("Error creating video record: %v", err)
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not create video record", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusCreated, newVideo)
}

func (cfg *apiConfig) handlerListUnlinkedVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := cfg.queries.ListUnlinkedVideos(r.Context())
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not gather unlinked videos", err)
		return
	}
	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)
}

func (cfg *apiConfig) handlerListVideosByMilestone(w http.ResponseWriter, r *http.Request) {
	milestoneID := r.PathValue("milestone_id")
	milestoneUUID, err := uuid.Parse(milestoneID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid milestone id format", err)
		return
	}

	// Pass a slice of uuid.UUID containing just this one ID
	videos, err := cfg.queries.ListVideosByMilestone(r.Context(), []uuid.UUID{milestoneUUID})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not gather videos by milestone", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)
}

func (cfg *apiConfig) handlerListVideosMissingTranscripts(w http.ResponseWriter, r *http.Request) {
	videos, err := cfg.queries.ListVideosMissingTranscripts(r.Context())
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not gather videos without transcripts", err)
		return
	}
	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)
}

func (cfg *apiConfig) handlerListVideosNotEnriched(w http.ResponseWriter, r *http.Request) {
	videos, err := cfg.queries.ListVideosNotEnriched(r.Context())
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not gather videos without enrichment", err)
		return
	}
	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)
}

type ListVideosBySummarySourceRequest struct {
	Source string `json:"source"`
}

func (cfg *apiConfig) handlerListVideosBySummarySource(w http.ResponseWriter, r *http.Request) {
	var req ListVideosBySummarySourceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
		return
	}

	videos, err := cfg.queries.ListVideosBySummarySource(r.Context(), req.Source)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not gather videos by summary source", err)
		return
	}
	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)

	w.WriteHeader(http.StatusNoContent)
}

type ListVideosByStatusRequest struct {
	Status string `json:"status"`
}

func (cfg *apiConfig) handlerListVideosByStatus(w http.ResponseWriter, r *http.Request) {
	var req ListVideosByStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
		return
	}

	videos, err := cfg.queries.ListVideosByStatus(r.Context(), req.Status)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not gather videos by status", err)
		return
	}
	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerGetVideoByID(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("video_id")
	videoUUID, err := uuid.Parse(videoID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid video id format", err)
		return
	}

	video, err := cfg.queries.GetVideoByID(r.Context(), videoUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonhelp.RespondWithError(w, http.StatusNotFound, "Video not found", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not get video by id", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, video)
}

func (cfg *apiConfig) handlerLinkVideoToMilestone(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("video_id")
	videoUUID, err := uuid.Parse(videoID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid video id format", err)
		return
	}

	milestoneID := r.PathValue("milestone_id")
	milestoneUUID, err := uuid.Parse(milestoneID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	nullMilestoneID := uuid.NullUUID{
		UUID:  milestoneUUID,
		Valid: true,
	}

	err = cfg.queries.LinkVideoToMilestone(r.Context(), database.LinkVideoToMilestoneParams{
		MilestoneID: nullMilestoneID,
		ID:          videoUUID,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not link video", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUnlinkVideoFromMilestone(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("video_id")
	videoUUID, err := uuid.Parse(videoID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid video id format", err)
		return
	}

	err = cfg.queries.UnlinkVideoFromMilestone(r.Context(), database.UnlinkVideoFromMilestoneParams{
		ID:        videoUUID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not unlink video", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateVideoStatusParams struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func (cfg *apiConfig) handlerUpdateVideoStatus(w http.ResponseWriter, r *http.Request) { // take a look at this for functionality
	videoID := r.PathValue("video_id")
	videoUUID, err := uuid.Parse(videoID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid video id format", err)
		return
	}

	var req UpdateVideoStatusParams
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// return the exact go decoding error to see what failed
		jsonhelp.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	err = cfg.queries.UpdateVideoStatus(r.Context(), database.UpdateVideoStatusParams{
		ID:        videoUUID,
		Status:    req.Status,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Failed to update video status", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Status updated successfully"}`))
}

type UpdateVideoCategoryParams struct {
	ID       uuid.UUID `json:"id"`
	Category string    `json:"category"`
}

func (cfg *apiConfig) handlerUpdateVideoCategory(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("video_id")
	videoUUID, err := uuid.Parse(videoID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid video id format", err)
		return
	}

	var req UpdateVideoCategoryParams
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// return the exact go decoding error to see what failed
		jsonhelp.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	err = cfg.queries.UpdateVideoCategory(r.Context(), database.UpdateVideoCategoryParams{
		ID:        videoUUID,
		Category:  req.Category,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Failed to update video category", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Category updated successfully"}`))
}

type UpdateVideoEstimatedEventDateRequest struct {
	EventDate string `json:"event_date"`
}

func (cfg *apiConfig) handlerUpdateVideoEstimatedEventDate(w http.ResponseWriter, r *http.Request) {
	var req UpdateVideoEstimatedEventDateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
		return
	}

	id := r.PathValue("video_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	// Parse the incoming string into sql.NullTime
	var nullTime sql.NullTime
	if req.EventDate != "" && req.EventDate != "null" {
		parsedTime, err := time.Parse("2006-01-02", req.EventDate) // HTML date inputs send "YYYY-MM-DD"
		if err != nil {
			jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD", err)
			return
		}
		nullTime = sql.NullTime{
			Time:  parsedTime,
			Valid: true,
		}
	} else {
		nullTime = sql.NullTime{Valid: false} // Clears the date if empty
	}

	now := time.Now().UTC()
	err = cfg.queries.UpdateVideoEstimatedEventDate(r.Context(), database.UpdateVideoEstimatedEventDateParams{
		ID:                 userUUID,
		EstimatedEventDate: nullTime,
		UpdatedAt:          now,
	})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not update estimated event date", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateVideoTranscriptRequest struct {
	RawTranscript *string `json:"event_date"`
}

func (cfg *apiConfig) handlerUpdateVideoTranscript(w http.ResponseWriter, r *http.Request) {
	var req UpdateVideoTranscriptRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
		return
	}

	id := r.PathValue("video_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	now := time.Now().UTC()
	err = cfg.queries.UpdateVideoTranscript(r.Context(), database.UpdateVideoTranscriptParams{
		ID:            userUUID,
		RawTranscript: req.RawTranscript,
		UpdatedAt:     now,
	})

	w.WriteHeader(http.StatusNoContent)
}

type UpdateVideoStatusRequest struct {
	EventdDate sql.NullTime `json:"event_date"`
}

func (cfg *apiConfig) handlerListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := cfg.queries.ListVideos(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonhelp.RespondWithError(w, http.StatusNotFound, "Videos not found", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not get videos", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, videos)
}

func (cfg *apiConfig) handlerDeleteVideo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("video_id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	err = cfg.queries.DeleteVideo(r.Context(), UUID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not delete video", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerEnrichVideo(w http.ResponseWriter, r *http.Request) {
	// get id from path
	videoIDStr := r.PathValue("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid video ID format", err)
		return
	}

	// optional quick check: ensure video exists before queuing
	_, err = cfg.queries.GetVideoByID(r.Context(), videoID)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonhelp.RespondWithError(w, http.StatusNotFound, "Video not found", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Database error retrieving video", err)
		return
	}

	// enqueue the job for the worker pool
	enqueued := cfg.wp.Enqueue(worker.Task{VideoID: videoID})
	if !enqueued {
		jsonhelp.RespondWithError(w, http.StatusServiceUnavailable, "Enrichment queue is full", nil)
		return
	}

	// immediate 202 response
	jsonhelp.RespondWithJSON(w, http.StatusAccepted, map[string]string{
		"message":  "Video enrichment enqueued successfully",
		"video_id": videoID.String(),
		"status":   "pending review",
	})
}
