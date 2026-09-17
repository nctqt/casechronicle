-- name: AddVideo :one
INSERT INTO videos (
    id,
    milestone_id,
    youtube_video_id,
    title,
    channel_name,
    status,
    category,
    ai_summary,
    raw_transcript,
    estimated_event_date,
    enriched_date,
    summary_source,
    published_at,
    created_at,
    updated_at
) 
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: GetVideoByID :one
SELECT *
FROM videos
WHERE id = $1
LIMIT 1;

-- name: ListVideos :many
SELECT *
FROM videos
ORDER BY created_at DESC;

-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = $1;

-- name: ListVideosByMilestone :many
SELECT * 
FROM videos
WHERE milestone_id = ANY(sqlc.slice('milestone_id')::uuid[])
  AND status IN ('analyzed', 'approved')
ORDER BY COALESCE(estimated_event_date, published_at) ASC;

-- name: ListUnlinkedVideos :many
SELECT * 
FROM videos
WHERE milestone_id IS NULL
ORDER BY published_at DESC;

-- name: ListVideosByStatus :many
SELECT *
FROM videos
WHERE status = $1
ORDER BY created_at DESC;

-- name: ListVideosNotEnriched :many
SELECT *
FROM videos
WHERE enriched_date IS NULL -- tie to ai_summary
ORDER BY created_at DESC;

-- name: ListVideosBySummarySource :many
SELECT *
FROM videos
WHERE summary_source = $1
ORDER BY created_at DESC;

-- name: ListVideosMissingTranscripts :many
SELECT * 
FROM videos
WHERE raw_transcript IS NULL 
ORDER BY created_at ASC;

-- name: LinkVideoToMilestone :exec
UPDATE videos
SET 
    milestone_id = $1,
    updated_at = $2
WHERE id = $3;

-- name: UnlinkVideoFromMilestone :exec
UPDATE videos
SET milestone_id = NULL,
    updated_at = $2
WHERE id = $1;

-- name: UpdateVideoStatus :exec
UPDATE videos
SET status = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateVideoTranscript :exec
UPDATE videos
SET raw_transcript = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateVideoCategory :exec
UPDATE videos
SET category = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateVideoEstimatedEventDate :exec
UPDATE videos
SET estimated_event_date = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateVideoSummary :one
UPDATE videos
SET 
    ai_summary = $2,
    estimated_event_date = $3,
    summary_source = $4,
    status = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;
