package article

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func testRouter(repository Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router, NewHandler(NewService(repository), slog.New(slog.NewTextHandler(io.Discard, nil))))
	return router
}

func TestCreateHandler(t *testing.T) {
	body, _ := json.Marshal(validInput())
	request := httptest.NewRequest(http.MethodPost, "/article/", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	testRouter(newFakeRepository()).ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestInvalidPaginationHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/article/nope/0", nil)
	response := httptest.NewRecorder()
	testRouter(newFakeRepository()).ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestNotFoundHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/article/999", nil)
	response := httptest.NewRecorder()
	testRouter(newFakeRepository()).ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
