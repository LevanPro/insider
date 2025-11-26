package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LevanPro/insider/internal/infra/database"
)

func (app *App) StartScheduler(w http.ResponseWriter, r *http.Request) {

}

func (app *App) StopScheduler(w http.ResponseWriter, r *http.Request) {

}

func (app *App) GetSentMessages(w http.ResponseWriter, r *http.Request) {

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
