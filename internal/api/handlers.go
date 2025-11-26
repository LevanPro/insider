package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LevanPro/insider/internal/infra/database"
)

// SchedulerStatus godoc
// @Summary      Get scheduler status
// @Description  Returns scheduler state
// @Tags         scheduler
// @Success      200  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/scheduler/status [get]
func (app *App) SchedulerStatus(w http.ResponseWriter, r *http.Request) {
	status := app.scheduler.IsRunning()
	statusCode := http.StatusOK

	data := struct {
		Status bool `json:"running"`
	}{
		Status: status,
	}

	if err := response(w, statusCode, data); err != nil {
		app.log.Errorw("SchedulerStatus", "ERROR", err)
	}
}

// StartScheduler godoc
// @Summary      Start automatic message sending
// @Description  Starts background job that every 2 minutes sends 2 unsent messages
// @Tags         scheduler
// @Success      200  {object} map[string]string
// @Failure      409  {object} map[string]string
// @Router       /api/v1/scheduler/start [post]
func (app *App) StartScheduler(w http.ResponseWriter, r *http.Request) {
	err := app.scheduler.Start()

	message := "scheduler has started"
	statusCode := http.StatusOK

	// Need checking not expose internal errors
	if err != nil {
		message = err.Error()
		statusCode = http.StatusConflict
	}

	data := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	if err := response(w, statusCode, data); err != nil {
		app.log.Errorw("StartScheduler", "ERROR", err)
	}
}

// StopScheduler godoc
// @Summary      Stop automatic message sending
// @Description  Stop background job that sends every 2 minutes sends 2 unsent messages
// @Tags         scheduler
// @Success      200  {object} map[string]string
// @Failure      409  {object} map[string]string
// @Router       /api/v1/scheduler/stop [post]
func (app *App) StopScheduler(w http.ResponseWriter, r *http.Request) {
	err := app.scheduler.Stop()

	message := "scheduler has stopped"
	statusCode := http.StatusOK

	// Need checking not expose internal errors
	if err != nil {
		message = err.Error()
		statusCode = http.StatusConflict
	}

	data := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	if err := response(w, statusCode, data); err != nil {
		app.log.Errorw("StopScheduler", "ERROR", err)
	}
}

// GetSentMessages godoc
// @Summary      List sent messages
// @Description  Returns a paginated list of messages with status = sent
// @Tags         messages
// @Param        limit   query   int   false  "Limit (default 50)"
// @Param        offset  query   int   false  "Offset (default 0)"
// @Success      200  {array}  domain.Message
// @Failure      500  {object} map[string]string
// @Router       /messages/sent [get]
func (app *App) GetSentMessages(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 50)
	offset := parseIntQuery(r, "offset", 0)

	msgs, err := app.service.ListSent(r.Context(), limit, offset)
	if err != nil {
		err = response(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "something went wrong",
		})
		return
	}

	if err := response(w, http.StatusOK, map[string]any{
		"data": msgs,
	}); err != nil {
		app.log.Errorw("liveness", "ERROR", err)
	}
}

func (app *App) Liveness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status,omitempty"`
	}{
		Status: "up",
	}

	statusCode := http.StatusOK

	if err := response(w, statusCode, data); err != nil {
		app.log.Errorw("liveness", "ERROR", err)
	}
}

func (app *App) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	if err := database.StatusCheck(ctx, app.db); err != nil {
		status = "db not ready yet"
		statusCode = http.StatusInternalServerError
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	if err := response(w, statusCode, data); err != nil {
		app.log.Errorw("readiness", "ERROR", err)
	}
}

func response(w http.ResponseWriter, statusCode int, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func parseIntQuery(r *http.Request, key string, def int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return def
	}
	return n
}
