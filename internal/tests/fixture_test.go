package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hesampakdaman/wallet-service/internal/bootstrap"
)

type Fixture struct {
	app    bootstrap.App
	pool   *pgxpool.Pool
	server *httptest.Server
	client *http.Client
}

func NewFixture(t *testing.T) Fixture {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	pool := db.TestPool(t)
	app := bootstrap.NewApp(logger, pool)

	srv := httptest.NewServer(app.Router)
	t.Cleanup(srv.Close)

	return Fixture{
		app:    app,
		pool:   pool,
		server: srv,
		client: srv.Client(),
	}
}

func (f Fixture) DoJSON(method, path string, body any, v any) (*http.Response, error) {
	resp, data, err := f.DoJSONRaw(method, path, body)
	if err != nil {
		return resp, err
	}

	if v == nil {
		return resp, nil
	}

	if err := json.Unmarshal(data, v); err != nil {
		return resp, fmt.Errorf(
			"DoJSON(%s %s) data: %s",
			method,
			path,
			data,
		)
	}

	return resp, nil
}

func (f Fixture) DoRequest(
	method, path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, []byte, error) {
	if f.server == nil || f.client == nil {
		return nil, nil, fmt.Errorf("fixture HTTP server not initialized")
	}

	req, err := http.NewRequest(method, f.server.URL+path, body)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return resp, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, data, nil
}

func (f Fixture) DoJSONRaw(method, path string, body any) (*http.Response, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, nil, fmt.Errorf("failed to encode request body: %w", err)
		}
		reqBody = &buf
	}

	headers := map[string]string{
		"Accept": "application/json",
	}

	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	return f.DoRequest(method, path, reqBody, headers)
}
