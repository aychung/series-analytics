package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	client           *http.Client
	userAgent        string
	maxResponseBytes int64
}

type FetcherResult struct {
	FinalURL    string
	StatusCode  int
	ContentType string
	Body        []byte
	Size        int64
}

const userAgent = "curl/8.11.0"
const timeout = 10 * time.Second
const maxResponseBytes = 2 * 1024 * 1024

func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent:        userAgent,
		maxResponseBytes: maxResponseBytes,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, rawURL string) (*FetcherResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	if !isHTMLContentType(contentType) {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	limited := io.LimitReader(resp.Body, f.maxResponseBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}

	return &FetcherResult{
		FinalURL:    resp.Request.URL.String(),
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		Body:        body,
		Size:        int64(len(body)),
	}, nil
}

func isHTMLContentType(contentType string) bool {
	return contentType == "" ||
		len(contentType) >= 9 && contentType[:9] == "text/html" ||
		len(contentType) >= 21 && contentType[:21] == "application/xhtml+xml"
}
