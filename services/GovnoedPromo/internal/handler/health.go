package handler

import (
	"database/sql"
	"net/http"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(
	db *sql.DB,
) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

type healthResponse struct {
	Status string `json:"status"`
}

func (h *HealthHandler) Handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := h.db.Ping()

	if err != nil {
		writeJSON(
			w,
			http.StatusInternalServerError,
			healthResponse{
				Status: "error",
			},
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		healthResponse{
			Status: "ok",
		},
	)
}
