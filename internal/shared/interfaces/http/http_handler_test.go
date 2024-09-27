package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Mock the postgres package
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) CheckDBConnection() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) GetConnection() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

func (m *MockDatabase) Close() {
	m.Called()
}

func (m *MockDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(*sql.Row)
}

func (m *MockDatabase) AutoMigrate() error {
	args := m.Called()
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{})  {}
func (m *MockLogger) Error(args ...interface{}) {}
func (m *MockLogger) Warn(args ...interface{})  {}
func (m *MockLogger) Debug(args ...interface{}) {}

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	database := new(MockDatabase)
	logger := new(MockLogger)
	sharedHandler := NewSharedHandler(database, logger)

	router := gin.New()
	router.GET("/healthz", sharedHandler.Healthz)

	t.Run("success", func(t *testing.T) {
		database.On("CheckDBConnection", mock.Anything).Return(nil).Once()

		// new a context with a request
		request, _ := http.NewRequest("GET", "/healthz", strings.NewReader(""))
		writer := httptest.NewRecorder()

		router.ServeHTTP(writer, request)

		if writer.Code != 200 {
			t.Errorf("expected status 200, got %d", writer.Code)
		}

		// Check if the response body is empty
		if writer.Body.String() != "" {
			t.Errorf("expected empty response body, got %s", writer.Body.String())
		}
	})

	t.Run("error in database connection", func(t *testing.T) {
		database.On("CheckDBConnection", mock.Anything).Return(errors.New("db connection fail")).Once()

		// new a context with a request
		request, _ := http.NewRequest("GET", "/healthz", strings.NewReader(""))
		writer := httptest.NewRecorder()

		router.ServeHTTP(writer, request)

		if writer.Code != 503 {
			t.Errorf("expected status 503, got %d", writer.Code)
		}

		// Check if the response body is empty
		if writer.Body.String() != "" {
			t.Errorf("expected empty response body, got %s", writer.Body.String())
		}
	})

	t.Run("bad request", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/healthz", strings.NewReader("testing"))
		writer := httptest.NewRecorder()

		router.ServeHTTP(writer, request)

		if writer.Code != 400 {
			t.Errorf("expected status 400, got %d", writer.Code)
		}

		// Check if the response body is empty
		if writer.Body.String() != "" {
			t.Errorf("expected empty response body, got %s", writer.Body.String())
		}
	})
}
