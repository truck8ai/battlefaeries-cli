package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/truck8ai/battlefaeries-cli/internal/config"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logEnabled bool
}

type APIResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Raw     json.RawMessage `json:"-"`
}

type LogEntry struct {
	Timestamp  string          `json:"ts"`
	Method     string          `json:"method"`
	Path       string          `json:"path"`
	Status     int             `json:"status"`
	DurationMs int64           `json:"duration_ms"`
	Success    bool            `json:"success"`
	Error      string          `json:"error,omitempty"`
	Request    json.RawMessage `json:"request,omitempty"`
	Response   json.RawMessage `json:"response"`
}

func New() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("not logged in. Run: bf login")
	}

	return &Client{
		baseURL: cfg.APIURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logEnabled: cfg.LogEnabled,
	}, nil
}

func NewWithKey(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetLogEnabled overrides the config-based log setting (for --log flag)
func (c *Client) SetLogEnabled(enabled bool) {
	c.logEnabled = enabled
}

func (c *Client) Get(path string) (json.RawMessage, error) {
	return c.do("GET", path, nil)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.do("POST", path, body)
}

func (c *Client) Delete(path string, body interface{}) (json.RawMessage, error) {
	return c.do("DELETE", path, body)
}

func (c *Client) do(method, path string, body interface{}) (json.RawMessage, error) {
	url := c.baseURL + "/api/agent" + path

	var bodyReader io.Reader
	var bodyJSON json.RawMessage
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
		bodyJSON = data
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "battlefaeries-cli/1.0")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	durationMs := time.Since(start).Milliseconds()

	if err != nil {
		c.writeLog(method, path, 0, durationMs, false, err.Error(), bodyJSON, nil)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.writeLog(method, path, resp.StatusCode, durationMs, false, "read error", bodyJSON, nil)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		c.writeLog(method, path, resp.StatusCode, durationMs, false, "parse error", bodyJSON, data)
		return nil, fmt.Errorf("invalid response: %s", string(data))
	}

	if !apiResp.Success {
		errMsg := apiResp.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("request failed with status %d", resp.StatusCode)
		}
		c.writeLog(method, path, resp.StatusCode, durationMs, false, errMsg, bodyJSON, data)
		return nil, fmt.Errorf("%s", errMsg)
	}

	c.writeLog(method, path, resp.StatusCode, durationMs, true, "", bodyJSON, data)
	return data, nil
}

func (c *Client) writeLog(method, path string, status int, durationMs int64, success bool, errMsg string, reqBody, respBody json.RawMessage) {
	if !c.logEnabled {
		return
	}

	entry := LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Method:     method,
		Path:       path,
		Status:     status,
		DurationMs: durationMs,
		Success:    success,
		Error:      errMsg,
		Request:    reqBody,
		Response:   respBody,
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return
	}

	logDir := config.LogDir()
	if err := os.MkdirAll(logDir, 0700); err != nil {
		return
	}

	f, err := os.OpenFile(config.LogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(line)
	f.Write([]byte("\n"))
}
