package main

import (
	"net/http"
)

func registerEndpoints(apiCfg *apiConfig) {
	// http router
	apiCfg.mux = http.NewServeMux()

	// system
	apiCfg.mux.HandleFunc("GET /", handlerRoot)               // public
	apiCfg.mux.HandleFunc("GET /api/healthz", handlerHealthz) // public

	// timeline
	apiCfg.mux.HandleFunc("GET /api/v1/cases/{case_id}/timeline", apiCfg.middlewareOptionalAuth(apiCfg.handlerGetCaseTimeline)) // public

	// cases
	apiCfg.mux.HandleFunc("POST /api/v1/cases", apiCfg.adminMiddleware(apiCfg.handlerAddCase))
	apiCfg.mux.HandleFunc("GET /api/v1/cases", apiCfg.handlerListCases)             // public
	apiCfg.mux.HandleFunc("GET /api/v1/cases/{case_id}", apiCfg.handlerGetCaseByID) // public
	apiCfg.mux.HandleFunc("DELETE /api/v1/cases/{case_id}", apiCfg.adminMiddleware(apiCfg.handlerDeleteCase))
	apiCfg.mux.HandleFunc("PATCH /api/v1/cases/{case_id}/description", apiCfg.adminMiddleware(apiCfg.handlerUpdateCaseDescription))
	apiCfg.mux.HandleFunc("PATCH /api/v1/cases/{case_id}/title", apiCfg.adminMiddleware(apiCfg.handlerUpdateCaseTitle))

	// users
	apiCfg.mux.HandleFunc("POST /api/v1/users/register", apiCfg.handlerAddUser) // public
	apiCfg.mux.HandleFunc("POST /api/v1/users/login", apiCfg.handlerLoginUser)  // public
	apiCfg.mux.HandleFunc("GET /api/v1/users", apiCfg.adminMiddleware(apiCfg.handlerListUsers))
	apiCfg.mux.HandleFunc("GET /api/v1/users/{user_id}", apiCfg.adminMiddleware(apiCfg.handlerGetUserByID))
	apiCfg.mux.HandleFunc("DELETE /api/v1/users/{user_id}", apiCfg.adminMiddleware(apiCfg.handlerDeleteUser))

	// milestones
	apiCfg.mux.HandleFunc("POST /api/v1/cases/{case_id}/milestones", apiCfg.adminMiddleware(apiCfg.handlerAddMilestone))
	apiCfg.mux.HandleFunc("GET /api/v1/milestones/{case_id}", apiCfg.handlerListMilestonesByCase) // public
	apiCfg.mux.HandleFunc("GET /api/v1/milestones/{milestone_id}", apiCfg.adminMiddleware(apiCfg.handlerGetMilestoneByID))
	apiCfg.mux.HandleFunc("DELETE /api/v1/milestones/{milestone_id}", apiCfg.adminMiddleware(apiCfg.handlerDeleteMilestone))
	apiCfg.mux.HandleFunc("PATCH /api/v1/cases/{case_id}/title", apiCfg.adminMiddleware(apiCfg.handlerUpdateMilestoneTitle))
	apiCfg.mux.HandleFunc("PATCH /api/v1/cases/{case_id}/description", apiCfg.adminMiddleware(apiCfg.handlerUpdateMilestoneDescription))
	apiCfg.mux.HandleFunc("PATCH /api/v1/cases/{case_id}/event_date", apiCfg.adminMiddleware(apiCfg.handlerUpdateMilestoneEventDate))

	// videos
	apiCfg.mux.HandleFunc("POST /api/v1/videos", apiCfg.adminMiddleware(apiCfg.handlerAddVideo))
	apiCfg.mux.HandleFunc("POST /api/v1/videos/{video_id}/enrich", apiCfg.adminMiddleware(apiCfg.handlerEnrichVideo))
	apiCfg.mux.HandleFunc("GET /api/v1/videos", apiCfg.adminMiddleware(apiCfg.handlerListVideos))
	apiCfg.mux.HandleFunc("GET /api/v1/videos/{video_id}", apiCfg.adminMiddleware(apiCfg.handlerGetVideoByID))
	apiCfg.mux.HandleFunc("GET /api/v1/videos/unlinked", apiCfg.adminMiddleware(apiCfg.handlerListUnlinkedVideos))
	apiCfg.mux.HandleFunc("GET /api/v1/videos/status", apiCfg.adminMiddleware(apiCfg.handlerListVideosByStatus))
	apiCfg.mux.HandleFunc("GET /api/v1/videos/summary_source", apiCfg.adminMiddleware(apiCfg.handlerListVideosBySummarySource))
	apiCfg.mux.HandleFunc("GET /api/v1/videos/not_enriched", apiCfg.adminMiddleware(apiCfg.handlerListVideosNotEnriched))
	apiCfg.mux.HandleFunc("GET /api/v1/videos/missing_transcripts", apiCfg.adminMiddleware(apiCfg.handlerListVideosMissingTranscripts))
	apiCfg.mux.HandleFunc("GET /api/v1/milestones/{milestone_id}/videos", apiCfg.handlerListVideosByMilestone) // public
	apiCfg.mux.HandleFunc("PATCH /api/v1/videos/{video_id}/category", apiCfg.adminMiddleware(apiCfg.handlerUpdateVideoCategory))
	apiCfg.mux.HandleFunc("PATCH /api/v1/videos/{video_id}/status", apiCfg.adminMiddleware(apiCfg.handlerUpdateVideoStatus))
	apiCfg.mux.HandleFunc("PATCH /api/v1/videos/{video_id}/transcript", apiCfg.adminMiddleware(apiCfg.handlerUpdateVideoTranscript))
	apiCfg.mux.HandleFunc("PATCH /api/v1/videos/{video_id}/event_date", apiCfg.adminMiddleware(apiCfg.handlerUpdateVideoEstimatedEventDate))
	apiCfg.mux.HandleFunc("DELETE /api/v1/videos/{video_id}", apiCfg.adminMiddleware(apiCfg.handlerDeleteMilestone))
	apiCfg.mux.HandleFunc("PUT /api/v1/milestones/{milestone_id}/videos/{video_id}", apiCfg.adminMiddleware(apiCfg.handlerLinkVideoToMilestone))
	apiCfg.mux.HandleFunc("DELETE /api/v1/milestones/{milestone_id}/videos/{video_id}", apiCfg.adminMiddleware(apiCfg.handlerUnlinkVideoFromMilestone))
}
