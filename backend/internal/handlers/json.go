package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (h *Handlers) logInternalServerError(r *http.Request, err error) {
	h.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (h *Handlers) readJSON(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}
	return nil
}

func (h *Handlers) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logInternalServerError(r, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handlers) successResponse(w http.ResponseWriter, r *http.Request, message string, data any) {
	h.writeJSON(w, r, http.StatusOK, response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *Handlers) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.writeJSON(w, r, http.StatusOK, response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	})
}

func (h *Handlers) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.logInternalServerError(r, err)
	h.errorResponse(w, r, errors.New("服务器内部错误"))
}
