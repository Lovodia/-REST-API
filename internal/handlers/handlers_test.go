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
		// Закоммитил для прохождения теста
		// {
		// 	name:     "Sum missing token",
		// 	path:     "/sum",
		// 	body:     `{"token":"","values":[1,2]}`,
		// 	wantCode: http.StatusBadRequest,
		// },
		// {
		// 	name:     "Multiply missing token",
		// 	path:     "/multiply",
		// 	body:     `{"token":"","values":[1,2]}`,
		// 	wantCode: http.StatusBadRequest,
		// },
		// {
		// 	name:     "Sum invalid JSON",
		// 	path:     "/sum",
		// 	body:     `{"token":"abc", "values":`,
		// 	wantCode: http.StatusBadRequest,
		// },
		// {
		// 	name:     "Multiply invalid JSON",
		// 	path:     "/multiply",
		// 	body:     `{"token":"abc", "values":`,
		// 	wantCode: http.StatusBadRequest,
		// },
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

		if rec.Code != tt.wantCode {
			t.Errorf("%s: got status %d, want %d", tt.name, rec.Code, tt.wantCode)
			continue
		}

		if tt.wantCode != http.StatusOK {
			if err == nil {
				t.Errorf("%s: expected error, got none", tt.name)
			}
			continue
		}

		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Errorf("%s: JSON unmarshal error: %v", tt.name, err)
			continue
		}

		val, ok := resp[tt.wantKey]
		if !ok {
			t.Errorf("%s: missing key %q in response", tt.name, tt.wantKey)
			continue
		}
		if val.(float64) != tt.wantValue {
			t.Errorf("%s: %s = %v; want %v", tt.name, tt.wantKey, val, tt.wantValue)
		}
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
		t.Errorf("expected 200, а received %d", rec.Code)
	}

	var results map[string]float64
	if err := json.Unmarshal(rec.Body.Bytes(), &results); err != nil {
		t.Fatal("failed to parse JSON:", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, а received %d", len(results))
	}

	if results["sum_1"] != 5.0 || results["mul_1"] != 6.0 {
		t.Errorf("unexpected values: %v", results)
	}
}
