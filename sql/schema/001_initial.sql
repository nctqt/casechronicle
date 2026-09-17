-- +goose Up
CREATE TABLE cases (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE milestones (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_milestones_case_id ON milestones(case_id);

CREATE TABLE videos (
    id UUID PRIMARY KEY,
    milestone_id UUID REFERENCES milestones(id) ON DELETE SET NULL,
    youtube_video_id VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending_review',
    category VARCHAR(50) NOT NULL, -- 'news', 'speculation', 'podcast'
    ai_summary TEXT,
    raw_transcript TEXT,
    estimated_event_date TIMESTAMPTZ,
    enriched_date TIMESTAMPTZ,
    summary_source VARCHAR(20) NOT NULL DEFAULT 'metadata',
    published_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_videos_milestone_id ON videos(milestone_id);

CREATE INDEX idx_videos_status ON videos(status);

CREATE INDEX idx_videos_transcript_pending 
    ON videos(id) 
    WHERE transcript_processed_date IS NULL AND raw_transcript IS NOT NULL;

CREATE INDEX idx_videos_milestone_status_event_date 
ON videos (milestone_id, status, COALESCE(estimated_event_date, published_at) ASC);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_videos_milestone_status_event_date;
DROP INDEX IF EXISTS idx_videos_transcript_pending;
DROP INDEX IF EXISTS idx_videos_status;
DROP INDEX IF EXISTS idx_videos_milestone_id;
DROP TABLE IF EXISTS videos;
DROP INDEX IF EXISTS idx_milestones_case_id;
DROP TABLE IF EXISTS milestones;
DROP TABLE IF EXISTS cases;
