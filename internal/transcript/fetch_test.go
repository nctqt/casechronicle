package transcript

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestFetchTranscript_RealSite(t *testing.T) {
	// Skip this test unless an environment variable is explicitly set.
	// Run with: INTEGRATION=true go test -v -run TestFetchTranscript_RealSite
	// if os.Getenv("INTEGRATION") != "true" {
	// t.Skip("Skipping integration test against real YouTube. Set INTEGRATION=true to run.")
	// }
	// Use a reliable, long-standing public video known to have captions
	// (e.g., Google's "Me at the zoo" or a stable tech trailer/announcement).
	// Let's use "Me at the zoo" (jNQXAC9IVRw) as a classic fallback, or any active video ID.
	testVideoID := "lbZim2SUJpw"

	// Create a real client with a reasonable timeout so the test doesn't hang indefinitely
	client := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	transcript, err := client.FetchTranscript(ctx, testVideoID)
	if err != nil {
		t.Fatalf("Failed to fetch transcript from real site for video %s: %v", testVideoID, err)
	}

	// Basic assertions to ensure we got actual data back
	if strings.TrimSpace(transcript) == "" {
		t.Error("Expected transcript text to be non-empty, but got empty string")
	}

	// Print a snippet of the result so you can visually verify it in the test logs
	snippet := transcript
	if len(snippet) > 100 {
		snippet = snippet[:100] + "..."
	}
	t.Logf("Successfully fetched real transcript! Snippet: %q", snippet)
}
