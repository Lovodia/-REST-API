package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Lovodia/restapi/internal/models"
	"github.com/Lovodia/restapi/internal/storage"
	"github.com/Lovodia/restapi/internal/usecase"

	"github.com/labstack/echo/v4"
)

// PostHandler godoc
// @Summary Сложение чисел
// @Description Принимает массив чисел и токен, возвращает сумму и сохраняет результат по токену
// @Tags math
// @Accept json
// @Produce json
// @Param input body models.Numbers true "Входные данные с токеном и числами"
// @Success 200 {object} models.SumResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /sum [post]
func PostHandler(logger *slog.Logger, store *storage.ResultStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		var nums models.Numbers
		if err := c.Bind(&nums); err != nil {
			logger.Error("Failed to bind request body", "error", err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data format")
		}

		if nums.Token == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Token is requived")
		}

		if nums.Values == nil {
			logger.Info("Received numbers", "values", "nil slice")
		} else {
			logger.Info("Received numbers", "values", nums.Values)
		}

		total := usecase.CalculateSum(nums.Values)

		resp := models.SumResponse{
			Token: nums.Token,
			Sum:   total,
		}

		logger.Info("Calculated sum", "sum", total)

		key := strconv.FormatInt(time.Now().UnixNano(), 10)
		store.Save(nums.Token, key, total)

		return c.JSON(http.StatusOK, resp)
	}
}

// MultiplyHandler godoc
// @Summary Умножение чисел
// @Description Принимает массив чисел и токен, возвращает произведение и сохраняет результат по токену
// @Tags math
// @Accept json
// @Produce json
// @Param input body models.Numbers true "Входные данные с токеном и числами"
// @Success 200 {object} models.MultiplyResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /multiply [post]
func MultiplyHandler(logger *slog.Logger, store *storage.ResultStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Numbers
		if err := c.Bind(&req); err != nil {
			logger.Error("Failed to bind multiply request", "error", err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data format")
		}
		if req.Token == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Token is requived")
		}
		if req.Values == nil {
			logger.Info("Received numbers", "values", "nil slice")
		} else {
			logger.Info("Received numbers", "values", req.Values)
		}

		multiply := usecase.CalculatedMultiply(req.Values)

		resp := models.MultiplyResponse{
			Token:    req.Token,
			Multiply: multiply,
		}
		logger.Info("Calculated multiply", "multiply", multiply)

		key := strconv.FormatInt(time.Now().UnixNano(), 10)
		store.Save(req.Token, key, multiply)

		return c.JSON(http.StatusOK, resp)
	}
}

// GetAllResultsByTokenHandler godoc
// @Summary Получение всех результатов по токену
// @Description Возвращает все сохранённые результаты (суммы и произведения) для заданного токена
// @Tags math
// @Produce json
// @Param token query string true "Токен пользователя"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} models.ErrorResponse
// @Router /results [get]
func GetAllResultsByTokenHandler(logger *slog.Logger, store *storage.ResultStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		if token == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Token query param is required")
		}
		results := store.GetAllByToken(token)
		if results == nil {
			results = map[string]float64{}
		}
		return c.JSON(http.StatusOK, results)
	}
}
