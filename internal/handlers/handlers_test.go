package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lovodia/restapi/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSumAndMultiplyHandlers(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name      string
		path      string
		body      string
		wantCode  int
		wantKey   string
		wantValue float64
	}{
		{
			name:      "Sum success",
			path:      "/sum",
			body:      `{"token":"abc123","values":[2,3,5]}`,
			wantCode:  http.StatusOK,
			wantKey:   "sum",
			wantValue: 10,
		},
		{
			name:      "Multiply success",
			path:      "/multiply",
			body:      `{"token":"abc123","values":[2,3,5]}`,
			wantCode:  http.StatusOK,
			wantKey:   "multiply",
			wantValue: 30,
		},
		{
			name:     "Sum missing token",
			path:     "/sum",
			body:     `{"token":"","values":[1,2]}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Multiply missing token",
			path:     "/multiply",
			body:     `{"token":"","values":[1,2]}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Sum invalid JSON",
			path:     "/sum",
			body:     `{"token":"abc", "values":`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Multiply invalid JSON",
			path:     "/multiply",
			body:     `{"token":"abc", "values":`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		store := storage.NewResultStore()
		var handler echo.HandlerFunc
		switch tt.path {
		case "/sum":
			handler = SumHandler(logger, store)
		case "/multiply":
			handler = MultiplyHandler(logger, store)
		}

		req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		require.Equal(t, tt.wantCode, rec.Code, "%s: HTTP status", tt.name)

		if tt.wantCode != http.StatusOK {
			require.Error(t, err, "%s: expected error", tt.name)
			continue
		}

		require.NoError(t, err, "%s: unexpected handler error", tt.name)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp), "%s: JSON unmarshal", tt.name)

		val, ok := resp[tt.wantKey]
		require.True(t, ok, "%s: missing key %q", tt.name, tt.wantKey)
		assert.Equal(t, tt.wantValue, val.(float64), "%s: %s value", tt.name, tt.wantKey)
	}
}
func TestGetAllResultsByTokenHandler_Success(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	store := storage.NewResultStore()

	token := "abc123"
	store.Save(token, "sum_1", 5.0)
	store.Save(token, "mul_1", 6.0)

	handler := GetAllResultsByTokenHandler(logger, store)

	req := httptest.NewRequest(http.MethodGet, "/results?token="+token, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, received %d", rec.Code)
	}

	var results map[string]float64
	if err := json.Unmarshal(rec.Body.Bytes(), &results); err != nil {
		t.Fatal("failed to parse JSON:", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, received %d", len(results))
	}

	if results["sum_1"] != 5.0 || results["mul_1"] != 6.0 {
		t.Errorf("unexpected values: %v", results)
	}
}
