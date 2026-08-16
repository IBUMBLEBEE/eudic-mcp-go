// Package eudic implements the Eudic / 欧路词典 OpenAPI study-list client.
package eudic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL         = "https://api.frdic.com/api/open/v1"
	DefaultLanguage = "en"
)

// Client calls the Eudic OpenAPI.
type Client struct {
	Token      string
	HTTPClient *http.Client
	BaseURL    string
}

// NewClientFromEnv builds a client using EUDIC_API_TOKEN (no "NIS " prefix).
func NewClientFromEnv() (*Client, error) {
	token := strings.TrimSpace(os.Getenv("EUDIC_API_TOKEN"))
	token = strings.TrimPrefix(token, "NIS ")
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("missing EUDIC_API_TOKEN (token only, without NIS prefix); get it from https://my.eudic.net/OpenAPI/Authorization")
	}
	return &Client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: BaseURL,
	}, nil
}

// APIError is a non-2xx Eudic response.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("eudic API %d: %s", e.StatusCode, e.Body)
}

func (c *Client) request(ctx context.Context, method, path string, query url.Values, payload any) (any, error) {
	base := c.BaseURL
	if base == "" {
		base = BaseURL
	}
	u, err := url.Parse(base + path)
	if err != nil {
		return nil, err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "eudic-mcp-go/1.0")
	req.Header.Set("Authorization", "NIS "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(raw)}
	}
	if resp.StatusCode == http.StatusNoContent || len(bytes.TrimSpace(raw)) == 0 {
		return map[string]any{"success": true}, nil
	}

	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{"raw": string(raw)}, nil
	}
	return out, nil
}

func langOrDefault(language string) string {
	if strings.TrimSpace(language) == "" {
		return DefaultLanguage
	}
	return language
}

// IDString normalizes category ids that may arrive as string or number.
func IDString(v any) string {
	switch t := v.(type) {
	case nil:
		return "0"
	case string:
		if t == "" {
			return "0"
		}
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case json.Number:
		return t.String()
	default:
		return fmt.Sprint(t)
	}
}

func (c *Client) ListCategories(ctx context.Context, language string) (any, error) {
	q := url.Values{"language": {langOrDefault(language)}}
	return c.request(ctx, http.MethodGet, "/studylist/category", q, nil)
}

func (c *Client) CreateCategory(ctx context.Context, name, language string) (any, error) {
	return c.request(ctx, http.MethodPost, "/studylist/category", nil, map[string]any{
		"language": langOrDefault(language),
		"name":     name,
	})
}

func (c *Client) RenameCategory(ctx context.Context, categoryID, name, language string) (any, error) {
	return c.request(ctx, http.MethodPatch, "/studylist/category", nil, map[string]any{
		"id":       categoryID,
		"language": langOrDefault(language),
		"name":     name,
	})
}

func (c *Client) DeleteCategory(ctx context.Context, categoryID, name, language string) (any, error) {
	return c.request(ctx, http.MethodDelete, "/studylist/category", nil, map[string]any{
		"id":       categoryID,
		"language": langOrDefault(language),
		"name":     name,
	})
}

func (c *Client) ListWords(ctx context.Context, categoryID, language string, page, pageSize int) (any, error) {
	if pageSize <= 0 {
		pageSize = 100
	}
	q := url.Values{
		"language":    {langOrDefault(language)},
		"category_id": {categoryID},
		"page":        {strconv.Itoa(page)},
		"page_size":   {strconv.Itoa(pageSize)},
	}
	return c.request(ctx, http.MethodGet, "/studylist/words", q, nil)
}

func (c *Client) AddWord(ctx context.Context, word, language string, categoryIDs []string, star int, contextLine string) (any, error) {
	if star <= 0 {
		star = 2
	}
	payload := map[string]any{
		"language": langOrDefault(language),
		"word":     word,
		"star":     star,
	}
	if len(categoryIDs) > 0 {
		payload["category_ids"] = categoryIDs
	}
	if contextLine != "" {
		payload["context_line"] = contextLine
	}
	return c.request(ctx, http.MethodPost, "/studylist/word", nil, payload)
}

func (c *Client) AddWordsBulk(ctx context.Context, words []string, categoryID, language string) (any, error) {
	return c.request(ctx, http.MethodPost, "/studylist/words", nil, map[string]any{
		"language":    langOrDefault(language),
		"category_id": categoryID,
		"words":       words,
	})
}

func (c *Client) DeleteWords(ctx context.Context, words []string, categoryID, language string) (any, error) {
	return c.request(ctx, http.MethodDelete, "/studylist/words", nil, map[string]any{
		"language":    langOrDefault(language),
		"category_id": categoryID,
		"words":       words,
	})
}

func (c *Client) GetWord(ctx context.Context, word, language string) (any, error) {
	q := url.Values{
		"language": {langOrDefault(language)},
		"word":     {word},
	}
	return c.request(ctx, http.MethodGet, "/studylist/word", q, nil)
}

func (c *Client) ListMasteredWords(ctx context.Context, language, categoryID string, page, pageSize int) (any, error) {
	if pageSize <= 0 {
		pageSize = 100
	}
	q := url.Values{
		"language":  {langOrDefault(language)},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}
	if categoryID != "" {
		q.Set("category_id", categoryID)
	}
	return c.request(ctx, http.MethodGet, "/studylist/mastered_words", q, nil)
}

func (c *Client) ListNotes(ctx context.Context, language string, page, pageSize int) (any, error) {
	if pageSize <= 0 {
		pageSize = 100
	}
	q := url.Values{
		"language":  {langOrDefault(language)},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}
	return c.request(ctx, http.MethodGet, "/studylist/notes", q, nil)
}

func (c *Client) GetNote(ctx context.Context, word, language string) (any, error) {
	q := url.Values{
		"language": {langOrDefault(language)},
		"word":     {word},
	}
	return c.request(ctx, http.MethodGet, "/studylist/note", q, nil)
}

func (c *Client) AddNote(ctx context.Context, word, note, language string) (any, error) {
	return c.request(ctx, http.MethodPost, "/studylist/note", nil, map[string]any{
		"language": langOrDefault(language),
		"word":     word,
		"note":     note,
	})
}

func (c *Client) DeleteNote(ctx context.Context, word, language string) (any, error) {
	return c.request(ctx, http.MethodDelete, "/studylist/note", nil, map[string]any{
		"language": langOrDefault(language),
		"word":     word,
	})
}

func (c *Client) ListSentences(ctx context.Context, language string, page, pageSize int) (any, error) {
	if pageSize <= 0 {
		pageSize = 100
	}
	q := url.Values{
		"language":  {langOrDefault(language)},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}
	return c.request(ctx, http.MethodGet, "/studylist/sentences", q, nil)
}
