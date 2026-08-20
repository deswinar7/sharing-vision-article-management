package article

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Request body must be valid JSON", nil)
		return
	}
	result, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) DispatchGet(c *gin.Context) {
	parts := strings.Split(strings.Trim(c.Param("path"), "/"), "/")
	switch len(parts) {
	case 1:
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil || id < 1 {
			writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Article id must be a positive integer", nil)
			return
		}
		h.getByID(c, id)
	case 2:
		limit, limitErr := strconv.Atoi(parts[0])
		offset, offsetErr := strconv.Atoi(parts[1])
		if limitErr != nil || offsetErr != nil {
			writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Pagination parameters must be integers", nil)
			return
		}
		h.list(c, limit, offset)
	default:
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Article path must contain an id or limit and offset", nil)
	}
}

func (h *Handler) list(c *gin.Context, limit, offset int) {
	filter := ListFilter{Limit: limit, Offset: offset}
	if value := c.Query("status"); value != "" {
		status := Status(value)
		filter.Status = &status
	}
	result, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) getByID(c *gin.Context, id int64) {
	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Request body must be valid JSON", nil)
		return
	}
	result, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	var validationError *ValidationError
	switch {
	case errors.As(err, &validationError):
		writeError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Request validation failed", validationError.Fields)
	case errors.Is(err, ErrNotFound):
		writeError(c, http.StatusNotFound, "ARTICLE_NOT_FOUND", "Article not found", nil)
	default:
		h.logger.Error("article request failed", "error", err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred", nil)
	}
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Article id must be a positive integer", nil)
		return 0, false
	}
	return id, true
}

func writeError(c *gin.Context, status int, code, message string, fields map[string]string) {
	c.JSON(status, errorEnvelope{Error: apiError{Code: code, Message: message, Fields: fields}})
}
