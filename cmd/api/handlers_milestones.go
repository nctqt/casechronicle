package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nctqt/casechronicle/internal/database"
	"github.com/nctqt/casechronicle/internal/jsonhelp"
)

type AddMilestoneRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	EventDate   time.Time `json:"event_date"`
}

func (cfg *apiConfig) handlerAddMilestone(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("case_id")
	caseUUID, err := uuid.Parse(caseID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	var req AddMilestoneRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Title == "" {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Title is required", nil)
		return
	}

	now := time.Now().UTC()
	newMilestone, err := cfg.queries.AddMilestone(r.Context(), database.AddMilestoneParams{
		ID:          uuid.New(),
		CaseID:      caseUUID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not create milestone", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusCreated, newMilestone)
}

func (cfg *apiConfig) handlerListMilestonesByCase(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("case_id")
	caseUUID, err := uuid.Parse(caseID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	milestones, err := cfg.queries.ListMilestonesByCase(r.Context(), caseUUID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not list milestones", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, milestones)
}

func (cfg *apiConfig) handlerGetMilestoneByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("milestone_id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	milestone, err := cfg.queries.GetMilestoneByID(r.Context(), UUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonhelp.RespondWithError(w, http.StatusNotFound, "Milestone not found", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve milestone", err)
		return
	}
	jsonhelp.RespondWithJSON(w, http.StatusOK, milestone)
}

func (cfg *apiConfig) handlerDeleteMilestone(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("milestone_id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	err = cfg.queries.DeleteMilestone(r.Context(), UUID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not delete milestone", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateMilestoneTitleRequest struct {
	Title string `json:"title"`
}

func (cfg *apiConfig) handlerUpdateMilestoneTitle(w http.ResponseWriter, r *http.Request) {
	var req UpdateMilestoneTitleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
	}

	id := r.PathValue("milestone_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	now := time.Now().UTC()
	err = cfg.queries.UpdateMilestoneTitle(r.Context(), database.UpdateMilestoneTitleParams{
		ID:        userUUID,
		Title:     req.Title,
		UpdatedAt: now,
	})

	w.WriteHeader(http.StatusNoContent)
}

type UpdateMilestoneDescriptionRequest struct {
	Description *string `json:"description"`
}

func (cfg *apiConfig) handlerUpdateMilestoneDescription(w http.ResponseWriter, r *http.Request) {
	var req UpdateMilestoneDescriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
	}

	id := r.PathValue("milestone_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	now := time.Now().UTC()
	err = cfg.queries.UpdateMilestoneDescription(r.Context(), database.UpdateMilestoneDescriptionParams{
		ID:          userUUID,
		Description: req.Description,
		UpdatedAt:   now,
	})

	w.WriteHeader(http.StatusNoContent)
}

type UpdateEventDateRequest struct {
	EventDate time.Time `json:"event_date"`
}

func (cfg *apiConfig) handlerUpdateMilestoneEventDate(w http.ResponseWriter, r *http.Request) {
	var req UpdateEventDateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could not parse body", err)
	}

	id := r.PathValue("milestone_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	now := time.Now().UTC()
	err = cfg.queries.UpdateEventDate(r.Context(), database.UpdateEventDateParams{
		ID:        userUUID,
		EventDate: req.EventDate,
		UpdatedAt: now,
	})

	w.WriteHeader(http.StatusNoContent)
}
