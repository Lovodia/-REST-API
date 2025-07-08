package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lovodia/restapi/internal/models"
	"github.com/Lovodia/restapi/internal/storage"
	"github.com/labstack/echo/v4"
)

func TestPostHandler_Success(t *testing.T) {

	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	store := storage.NewResultStore()
	handler := PostHandler(logger, store)

	reqBody := `{"token": "abc123", "values": [2.0, 3.0]}`
	req := httptest.NewRequest(http.MethodPost, "/sum", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("ожидали 200, а получили %d", rec.Code)
	}

	var resp models.SumResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal("не удалось распарсить ответ JSON:", err)
	}

	expected := 5.0
	if resp.Sum != expected {
		t.Errorf("ожидали сумму %.1f, а получили %.1f", expected, resp.Sum)
	}
}

func TestMultiplyHandler_Success(t *testing.T) {

	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	store := storage.NewResultStore()
	handler := MultiplyHandler(logger, store)

	reqBody := `{"token": "abc123", "values": [2.0, 3.0]}`
	req := httptest.NewRequest(http.MethodPost, "/multiply", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("ожидали 200, а получили %d", rec.Code)
	}

	var resp models.MultiplyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal("не удалось распарсить ответ JSON:", err)
	}

	expected := 6.0
	if resp.Multiply != expected {
		t.Errorf("ожидали произведение %.1f, а получили %.1f", expected, resp.Multiply)
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
		t.Errorf("ожидали 200, а получили %d", rec.Code)
	}

	var results map[string]float64
	if err := json.Unmarshal(rec.Body.Bytes(), &results); err != nil {
		t.Fatal("не удалось распарсить JSON:", err)
	}

	if len(results) != 2 {
		t.Errorf("ожидали 2 результата, а получили %d", len(results))
	}

	if results["sum_1"] != 5.0 || results["mul_1"] != 6.0 {
		t.Errorf("неожиданные значения: %v", results)
	}
}
