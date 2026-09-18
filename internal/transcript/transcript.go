package transcript

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoTranscriptFound = errors.New("no transcript available for this video")

	// Regex matching WebVTT timestamp lines (e.g. 00:01:20.500 --> 00:01:23.000)
	vttTimestampRegex = regexp.MustCompile(`(?m)^\d{2}:\d{2}(:\d{2})?\.\d{3}\s+-->\s+\d{2}:\d{2}(:\d{2})?\.\d{3}.*$\n?`)
	// Regex matching WebVTT header tags, positioning rules, or XML tags like <c> text </c>
	vttTagRegex = regexp.MustCompile(`<[^>]*>`)
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fetchTranscript attempts to retrieve and clean captions for a given YouTube video ID.
// relies on yt-dlp command line tool
func (c *Client) FetchTranscript(ctx context.Context, videoID string) (string, error) {
	if strings.TrimSpace(videoID) == "" {
		return "", errors.New("videoID cannot be empty")
	}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	// yt-dlp to download
	cmd := exec.Command("yt-dlp",
		"--write-auto-sub",
		"--skip-download",
		"-o", "transcript.%(ext)s",
		url,
	)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to download subtitle: %w", err)
	}

	// find file
	matches, err := filepath.Glob("transcript.*.vtt")
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("subtitle file not found")
	}

	subFile := matches[0]

	contentBytes, err := os.ReadFile(subFile)
	if err != nil {
		return "", fmt.Errorf("failed to read subtitle file: %w", err)
	}

	cleanText := CleanVTT(string(contentBytes))
	if strings.TrimSpace(cleanText) == "" {
		return "", ErrNoTranscriptFound
	}

	return cleanText, nil
}

// CleanVTT strips WebVTT headers, timecodes, XML tags, numeric cue IDs, and deduplicates repeating lines.
func CleanVTT(vtt string) string {
	lines := strings.Split(vtt, "\n")
	var filteredLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" ||
			strings.HasPrefix(trimmed, "WEBVTT") ||
			strings.HasPrefix(trimmed, "Kind:") ||
			strings.HasPrefix(trimmed, "Language:") ||
			strings.HasPrefix(trimmed, "NOTE") {
			continue
		}
		filteredLines = append(filteredLines, line)
	}

	cleaned := strings.Join(filteredLines, "\n")

	// Strip timecodes (e.g. 00:00:01.000 --> 00:00:04.000)
	cleaned = vttTimestampRegex.ReplaceAllString(cleaned, "")

	// Strip XML inline tags (e.g. <c>text</c> or <00:00:01.500>)
	cleaned = vttTagRegex.ReplaceAllString(cleaned, "")

	rawLines := strings.Split(cleaned, "\n")
	var finalLines []string
	var lastLine string

	for _, line := range rawLines {
		t := strings.TrimSpace(line)

		// Skip empty lines or pure integer sequence markers
		if t == "" || isInteger(t) || t == lastLine {
			continue
		}

		finalLines = append(finalLines, t)
		lastLine = t
	}

	fullTranscript := strings.Join(finalLines, " ")

	// Truncate to maximum ~18,000 characters (~4,000 tokens) for LLM context safety
	const maxChars = 18000
	if len(fullTranscript) > maxChars {
		fullTranscript = fullTranscript[:maxChars] + "..."
	}

	return fullTranscript
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
