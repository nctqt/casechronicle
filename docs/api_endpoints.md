# API Documentation

## Overview

This API provides system status endpoints plus versioned resources for users, cases, milestones, and videos. It is implemented with Go's `net/http` route patterns.

All versioned application endpoints use the `/api/v1` prefix. Path parameters use braces, such as `{case_id}`.

### Access levels

| Access | Meaning |
|---|---|
| Public | The route has no explicit middleware in the registration. |
| Optional authentication | The route uses `middlewareOptionalAuth`; unauthenticated access is allowed, while authenticated callers may receive different data. |
| Admin required | The route uses `adminMiddleware`. |

Middleware expects bearer tokens; Administrative requests use:

```http
Authorization: Bearer <admin-access-token>
```

```bash
# Public request
curl "$API_BASE_URL/api/healthz"

# Administrative request
curl \
  -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/users"
```

## Endpoint summary

| Method | Path | Access | Handler |
|---|---|---|---|
| GET | `/` | Public | `handlerRoot` |
| GET | `/api/healthz` | Public | `handlerHealthz` |
| POST | `/api/v1/users/register` | Public | `handlerAddUser` |
| POST | `/api/v1/users/login` | Public | `handlerLoginUser` |
| GET | `/api/v1/users` | Admin required | `handlerListUsers` |
| GET | `/api/v1/users/{user_id}` | Admin required | `handlerGetUserByID` |
| DELETE | `/api/v1/users/{user_id}` | Admin required | `handlerDeleteUser` |
| GET | `/api/v1/cases` | Public | `handlerListCases` |
| POST | `/api/v1/cases` | Admin required | `handlerAddCase` |
| GET | `/api/v1/cases/{case_id}` | Public | `handlerGetCaseByID` |
| DELETE | `/api/v1/cases/{case_id}` | Admin required | `handlerDeleteCase` |
| PATCH | `/api/v1/cases/{case_id}/title` | Admin required | `handlerUpdateCaseTitle` |
| PATCH | `/api/v1/cases/{case_id}/description` | Admin required | `handlerUpdateCaseDescription` |
| GET | `/api/v1/cases/{case_id}/timeline` | Optional authentication | `handlerGetCaseTimeline` |
| GET | `/api/v1/cases/{case_id}/milestones` | Public | `handlerListMilestonesByCase` |
| POST | `/api/v1/cases/{case_id}/milestones` | Admin required | `handlerAddMilestone` |
| GET | `/api/v1/milestones/{milestone_id}` | Admin required | `handlerGetMilestoneByID` |
| DELETE | `/api/v1/milestones/{milestone_id}` | Admin required | `handlerDeleteMilestone` |
| PATCH | `/api/v1/milestones/{milestone_id}/title` | Admin required | `handlerUpdateMilestoneTitle` |
| PATCH | `/api/v1/milestones/{milestone_id}/description` | Admin required | `handlerUpdateMilestoneDescription` |
| PATCH | `/api/v1/milestones/{milestone_id}/event_date` | Admin required | `handlerUpdateMilestoneEventDate` |
| GET | `/api/v1/videos` | Admin required | `handlerListVideos` |
| POST | `/api/v1/videos` | Admin required | `handlerAddVideo` |
| GET | `/api/v1/videos/unlinked` | Admin required | `handlerListUnlinkedVideos` |
| GET | `/api/v1/videos/status` | Admin required | `handlerListVideosByStatus` |
| GET | `/api/v1/videos/summary_source` | Admin required | `handlerListVideosBySummarySource` |
| GET | `/api/v1/videos/not_enriched` | Admin required | `handlerListVideosNotEnriched` |
| GET | `/api/v1/videos/missing_transcripts` | Admin required | `handlerListVideosMissingTranscripts` |
| GET | `/api/v1/videos/{video_id}` | Admin required | `handlerGetVideoByID` |
| DELETE | `/api/v1/videos/{video_id}` | Admin required | `handlerDeleteVideo` |
| PATCH | `/api/v1/videos/{video_id}/enrich` | Admin required | `handlerEnrichVideo` |
| PATCH | `/api/v1/videos/{video_id}/category` | Admin required | `handlerUpdateVideoCategory` |
| PATCH | `/api/v1/videos/{video_id}/status` | Admin required | `handlerUpdateVideoStatus` |
| PATCH | `/api/v1/videos/{video_id}/transcript` | Admin required | `handlerUpdateVideoTranscript` |
| PATCH | `/api/v1/videos/{video_id}/event_date` | Admin required | `handlerUpdateVideoEstimatedEventDate` |
| GET | `/api/v1/milestones/{milestone_id}/videos` | Public | `handlerListVideosByMilestone` |
| PUT | `/api/v1/milestones/{milestone_id}/videos/{video_id}` | Admin required | `handlerLinkVideoToMilestone` |
| DELETE | `/api/v1/milestones/{milestone_id}/videos/{video_id}` | Admin required | `handlerUnlinkVideoFromMilestone` |

## System

### `GET /`

Returns the root service response.

- Access: Public
- Handler: `handlerRoot`
- Response: TBD

```bash
curl "$API_BASE_URL/"
```

### `GET /api/healthz`

Returns the service health status.

- Access: Public
- Handler: `handlerHealthz`
- Response: TBD

```bash
curl "$API_BASE_URL/api/healthz"
```

## Users

### `POST /api/v1/users/register`

Registers a user.

- Access: Public
- Handler: `handlerAddUser`
- Content type: Confirm from the handler implementation; likely `application/json`
- Request body: TBD
- Success response: TBD

```bash
curl -X POST "$API_BASE_URL/api/v1/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "...": "See handlerAddUser for accepted registration fields"
  }'
```

### `POST /api/v1/users/login`

Authenticates a user.

- Access: Public
- Handler: `handlerLoginUser`
- Request body: TBD
- Success response: TBD; confirm whether it returns a token, session, cookie, or another credential.

```bash
curl -X POST "$API_BASE_URL/api/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{
    "...": "See handlerLoginUser for credential fields"
  }'
```

### `GET /api/v1/users`

Lists users.

- Access: Admin required
- Handler: `handlerListUsers`
- Query parameters: TBD
- Response: TBD

```bash
curl \
  -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/users"
```

### `GET /api/v1/users/{user_id}`

Retrieves a user by ID.

| Parameter | Location |
|---|---|
| `user_id` | Path |

- Access: Admin required
- Handler: `handlerGetUserByID`
- Response: TBD

```bash
curl -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/users/<user_id>"
```

### `DELETE /api/v1/users/{user_id}`

Deletes a user by ID.

| Parameter | Location |
|---|---|
| `user_id` | Path |

- Access: Admin required
- Handler: `handlerDeleteUser`
- Response: TBD

```bash
curl -X DELETE -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/users/<user_id>"
```

## Cases and timeline

### `GET /api/v1/cases`

Lists cases.

- Access: Public
- Handler: `handlerListCases`
- Query parameters: TBD
- Response: TBD

```bash
curl "$API_BASE_URL/api/v1/cases"
```

### `POST /api/v1/cases`

Creates a case.

- Access: Admin required
- Handler: `handlerAddCase`
- Request body: TBD
- Response: TBD

```bash
curl -X POST "$API_BASE_URL/api/v1/cases" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "...": "See handlerAddCase for supported fields"
  }'
```

### `GET /api/v1/cases/{case_id}`

Retrieves a case by ID.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Public
- Handler: `handlerGetCaseByID`
- Response: TBD

```bash
curl "$API_BASE_URL/api/v1/cases/<case_id>"
```

### `DELETE /api/v1/cases/{case_id}`

Deletes a case by ID.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Admin required
- Handler: `handlerDeleteCase`
- Response: TBD

```bash
curl -X DELETE -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/cases/<case_id>"
```

### `PATCH /api/v1/cases/{case_id}/title`

Updates a case title.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Admin required
- Handler: `handlerUpdateCaseTitle`
- Request body: Confirm the exact field contract in the handler.
- Response: TBD

```bash
curl -X PATCH "$API_BASE_URL/api/v1/cases/<case_id>/title" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated case title"}'
```

### `PATCH /api/v1/cases/{case_id}/description`

Updates a case description.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Admin required
- Handler: `handlerUpdateCaseDescription`
- Request body: Confirm the exact field contract in the handler.
- Response: TBD

```bash
curl -X PATCH "$API_BASE_URL/api/v1/cases/<case_id>/description" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"description":"Updated case description"}'
```

### `GET /api/v1/cases/{case_id}/timeline`

Returns a case timeline.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Optional authentication
- Handler: `handlerGetCaseTimeline`
- Response: TBD; verify whether authentication changes the returned data.

```bash
curl "$API_BASE_URL/api/v1/cases/<case_id>/timeline"
```

## Milestones

### `GET /api/v1/cases/{case_id}/milestones`

Lists milestones belonging to a case.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Public
- Handler: `handlerListMilestonesByCase`
- Query parameters: TBD
- Response: TBD

```bash
curl "$API_BASE_URL/api/v1/cases/<case_id>/milestones"
```

### `POST /api/v1/cases/{case_id}/milestones`

Creates a milestone under a case.

| Parameter | Location |
|---|---|
| `case_id` | Path |

- Access: Admin required
- Handler: `handlerAddMilestone`
- Request body: TBD
- Response: TBD

```bash
curl -X POST "$API_BASE_URL/api/v1/cases/<case_id>/milestones" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "...": "See handlerAddMilestone for supported fields"
  }'
```

### `GET /api/v1/milestones/{milestone_id}`

Retrieves a milestone by ID.

| Parameter | Location |
|---|---|
| `milestone_id` | Path |

- Access: Admin required
- Handler: `handlerGetMilestoneByID`
- Response: TBD

```bash
curl -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/milestones/<milestone_id>"
```

### `DELETE /api/v1/milestones/{milestone_id}`

Deletes a milestone by ID.

| Parameter | Location |
|---|---|
| `milestone_id` | Path |

- Access: Admin required
- Handler: `handlerDeleteMilestone`
- Response: TBD

```bash
curl -X DELETE -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/milestones/<milestone_id>"
```

### Milestone updates

All milestone update endpoints require administrator access and use `milestone_id` as a path parameter.

| Method | Path | Handler | Operation |
|---|---|---|---|
| PATCH | `/api/v1/milestones/{milestone_id}/title` | `handlerUpdateMilestoneTitle` | Updates the title. |
| PATCH | `/api/v1/milestones/{milestone_id}/description` | `handlerUpdateMilestoneDescription` | Updates the description. |
| PATCH | `/api/v1/milestones/{milestone_id}/event_date` | `handlerUpdateMilestoneEventDate` | Updates the event date. |

```bash
curl -X PATCH "$API_BASE_URL/api/v1/milestones/<milestone_id>/title" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated milestone title"}'

curl -X PATCH "$API_BASE_URL/api/v1/milestones/<milestone_id>/description" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"description":"Updated milestone description"}'

curl -X PATCH "$API_BASE_URL/api/v1/milestones/<milestone_id>/event_date" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"event_date":"YYYY-MM-DD"}'
```

## Videos

### `GET /api/v1/videos`

Lists videos.

- Access: Admin required
- Handler: `handlerListVideos`
- Query parameters: TBD
- Response: TBD

```bash
curl -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/videos"
```

### `POST /api/v1/videos`

Creates a video.

- Access: Admin required
- Handler: `handlerAddVideo`
- Request body: TBD
- Response: TBD

```bash
curl -X POST "$API_BASE_URL/api/v1/videos" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "...": "See handlerAddVideo for supported fields"
  }'
```

### Video filters

All video filter endpoints require administrator access. Confirm the supported query parameters in the handler implementations.

| Method | Path | Handler | Operation |
|---|---|---|---|
| GET | `/api/v1/videos/unlinked` | `handlerListUnlinkedVideos` | Lists videos not linked to a milestone. |
| GET | `/api/v1/videos/status` | `handlerListVideosByStatus` | Lists videos filtered by status. |
| GET | `/api/v1/videos/summary_source` | `handlerListVideosBySummarySource` | Lists videos filtered by summary source. |
| GET | `/api/v1/videos/not_enriched` | `handlerListVideosNotEnriched` | Lists videos that are not enriched. |
| GET | `/api/v1/videos/missing_transcripts` | `handlerListVideosMissingTranscripts` | Lists videos without transcripts. |

```bash
curl -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/videos/unlinked"

curl -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/videos/status?<filter>=<value>"
```

### `GET /api/v1/videos/{video_id}`

Retrieves a video by ID.

| Parameter | Location |
|---|---|
| `video_id` | Path |

- Access: Admin required
- Handler: `handlerGetVideoByID`
- Response: TBD

```bash
curl -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/videos/<video_id>"
```

### `DELETE /api/v1/videos/{video_id}`

Deletes a video by ID.

| Parameter | Location |
|---|---|
| `video_id` | Path |

- Access: Admin required
- Handler: `handlerDeleteVideo`
- Response: TBD

```bash
curl -X DELETE -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/videos/<video_id>"
```

### Video updates and enrichment

All endpoints below require administrator access and use `video_id` as a path parameter. Confirm exact bodies and response formats in the relevant handlers.

| Method | Path | Handler | Operation |
|---|---|---|---|
| PATCH | `/api/v1/videos/{video_id}/enrich` | `handlerEnrichVideo` | Enriches a video. |
| PATCH | `/api/v1/videos/{video_id}/category` | `handlerUpdateVideoCategory` | Updates the category. |
| PATCH | `/api/v1/videos/{video_id}/status` | `handlerUpdateVideoStatus` | Updates the status. |
| PATCH | `/api/v1/videos/{video_id}/transcript` | `handlerUpdateVideoTranscript` | Updates the transcript. |
| PATCH | `/api/v1/videos/{video_id}/event_date` | `handlerUpdateVideoEstimatedEventDate` | Updates the estimated event date. |

```bash
curl -X PATCH "$API_BASE_URL/api/v1/videos/<video_id>/enrich" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{}'

curl -X PATCH "$API_BASE_URL/api/v1/videos/<video_id>/category" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"category":"<category>"}'

curl -X PATCH "$API_BASE_URL/api/v1/videos/<video_id>/status" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"<status>"}'

curl -X PATCH "$API_BASE_URL/api/v1/videos/<video_id>/transcript" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"transcript":"<transcript text>"}'

curl -X PATCH "$API_BASE_URL/api/v1/videos/<video_id>/event_date" \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{"event_date":"YYYY-MM-DD"}'
```

## Milestone-video links

### `GET /api/v1/milestones/{milestone_id}/videos`

Lists videos linked to a milestone.

| Parameter | Location |
|---|---|
| `milestone_id` | Path |

- Access: Public
- Handler: `handlerListVideosByMilestone`
- Query parameters: TBD
- Response: TBD

```bash
curl "$API_BASE_URL/api/v1/milestones/<milestone_id>/videos"
```

### `PUT /api/v1/milestones/{milestone_id}/videos/{video_id}`

Links a video to a milestone.

| Parameter | Location |
|---|---|
| `milestone_id` | Path |
| `video_id` | Path |

- Access: Admin required
- Handler: `handlerLinkVideoToMilestone`
- Request body: Confirm handler behavior; no body is implied by the route.
- Response: TBD

```bash
curl -X PUT -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/milestones/<milestone_id>/videos/<video_id>"
```

### `DELETE /api/v1/milestones/{milestone_id}/videos/{video_id}`

Unlinks a video from a milestone.

| Parameter | Location |
|---|---|
| `milestone_id` | Path |
| `video_id` | Path |

- Access: Admin required
- Handler: `handlerUnlinkVideoFromMilestone`
- Response: TBD

```bash
curl -X DELETE -H "Authorization: Bearer <admin-access-token>" \
  "$API_BASE_URL/api/v1/milestones/<milestone_id>/videos/<video_id>"
```
