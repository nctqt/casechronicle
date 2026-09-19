# CaseChronicle

CaseChronicle is a full-stack true-crime timeline application that lets users explore a case as it unfolded. It organizes major events chronologically and links them with the news coverage, YouTube analysis, and commentary available at the time.

Built as a capstone project for the Boot.dev Backend Engineering course, it addresses the challenge of organizing fragmented case media into a navigable, event-based research experience.

# https://case-chronicle.com/

- future picture
##
- future .gif

## Video Enrichment

CaseChronicle uses an LLM-assisted enrichment pipeline to make imported videos
more useful than their platform metadata alone.

For each video, the backend can process available video descriptions and
subtitles/transcripts to generate:

- An enriched summary that better represents the video's content
- An estimated event date based on the events discussed in the video
- Structured context that helps distinguish the event date from the video's
  publication date
- Data that helps administrators associate a video with the most relevant
  case milestone

This allows a timeline to reflect when an event occurred, rather than only
when a video discussing that event was published.

## Features

- Browse cases and their major events in chronological order
- View videos associated with individual case milestones
- Compare early reporting and commentary with later developments
- Register and authenticate users
- Manage cases, milestones, videos, and video-to-milestone links through an admin API
- Identify unlinked videos, videos missing transcripts, and videos pending enrichment
- Store application data in PostgreSQL with versioned Goose migrations

## Tech Stack

### Backend

- Go
- PostgreSQL 16
- SQLC for type-safe database queries
- Goose for database migrations
- Go standard library HTTP routing
- YouTube Data API, OpenRouter API, and yt-dlp integrations

### Frontend

- React
- Vite
- TypeScript
- Tailwind CSS and CSS Modules
- Native Fetch API

## Architecture & Design

- **Type-Safe SQL:** Uses `sqlc` to generate type-safe Go code from raw SQL
  queries.
- **Versioned Schema Management:** Uses Goose migrations to evolve the
  PostgreSQL schema predictably.
- **Containerized Database:** Runs PostgreSQL locally in Docker using
  `postgres:16-alpine`.
- **REST API and SPA:** A Go REST API provides application data to a decoupled
  React single-page application.
- **LLM-Assisted Enrichment:** Processes video descriptions and available
  transcripts to generate enriched summaries and estimate the date of the
  event discussed, supporting more accurate timeline placement.

## Project Status

This is an active portfolio and capstone project. Features, data models, and integrations may change while development continues.

## Roadmap

- Improve event-date extraction and confidence scoring
- Automate video-to-milestone matching using enriched video metadata
- Support human approval before generated content changes public timelines
- Explore assisted case construction, where enriched video data can suggest
  milestones and chronology for administrator review
- Add additional source types beyond YouTube
