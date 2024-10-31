package test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go-template/internal/auth"
	"go-template/internal/config"
	sharedConfig "go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/database"
	sharedHttp "go-template/internal/shared/interfaces/http"
	"go-template/internal/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{})  {}
func (m *MockLogger) Error(args ...interface{}) {}
func (m *MockLogger) Debug(args ...interface{}) {}
func (m *MockLogger) Warn(args ...interface{})  {}

// Mock cloudwatch module
type MockCloudWatchModule struct {
	mock.Mock
}

func (m *MockCloudWatchModule) PublishMetric(namespace, metricName string, value float64, unit types.StandardUnit) {
}

func TestAuthAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rootPath, _ := utils.GetProjectRootPath()
	viper.AddConfigPath(rootPath)
	if err := sharedConfig.Load(&config.App); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mockLogger := new(MockLogger)
	mockCloudWatchModule := new(MockCloudWatchModule)
	database := initDatabase(mockCloudWatchModule)
	defer database.Close()

	// Setup server
	server := sharedHttp.NewServer()
	server.AddModules(
		auth.NewModule(database, mockLogger),
	)

	go func() {
		if err := server.Start(":3000"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for server to start
	if err := waitForServerStart(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	router := server.GetRouter()

	email := "test@example.com"
	password := "password"
	firstName := "John"
	lastName := "Doe"

	t.Run("TestCreateUser", func(t *testing.T) {
		// Create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"password": "password",
			"first_name": "John",
			"last_name": "Doe"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Check the response body
		assert.Contains(t, w.Body.String(), email)
		assert.Contains(t, w.Body.String(), firstName)
		assert.Contains(t, w.Body.String(), lastName)
		assert.NotContains(t, w.Body.String(), password)

		// Check if the user is created in the database
		var count int
		database.GetConnection().QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
		assert.Equal(t, 1, count)

		// Check if the password is hashed
		var hashedPassword string
		database.GetConnection().QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&hashedPassword)
		assert.NotEqual(t, password, hashedPassword)

		// 422
		req, _ = http.NewRequest("POST", "/v1/user", bytes.NewBuffer([]byte(`{
			"email": "test",
			"password": "password",
			"first_name": "John",
			"last_name": "Doe"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		// 400 - email already exists
		req, _ = http.NewRequest("POST", "/v1/user", bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"password": "password",
			"first_name": "John",
			"last_name": "Doe"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestUpdateUser", func(t *testing.T) {
		newFirstName := "Jane"
		newLastName := "Smith"
		newPassword := "newpassword"

		// Create a request to pass to our handler
		req, _ := http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer([]byte(
			fmt.Sprintf(`{
				"first_name": "%s",
				"last_name": "%s",
				"password": "%s"
			}`, newFirstName, newLastName, newPassword),
		)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic dGVzdEBleGFtcGxlLmNvbTpwYXNzd29yZA==")
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())

		// Check if the user is updated in the database
		var newPasswordHash string
		err := database.GetConnection().QueryRow("SELECT password FROM users WHERE email = $1 and first_name = $2 and last_name = $3", email, newFirstName, newLastName).Scan(&newPasswordHash)
		assert.NoError(t, err)
		assert.NotEqual(t, newPassword, newPasswordHash)

		// 422
		req, _ = http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer([]byte(`{
			"last_name": "Doe",
			"password": "password"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		base64Token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, newPassword)))

		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Token))
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		// 401 - Unauthorized
		// wrong token
		req, _ = http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"first_name": "John",
			"last_name": "Doe",
			"password": "password"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		base64Token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, password)))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Token))
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// empty token
		req, _ = http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"first_name": "John",
			"last_name": "Doe",
			"password": "password"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// 400 - extra json field
		req, _ = http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"first_name": "John",
			"last_name": "Doe",
			"password": "password"
		}`)))
		req.Header.Set("Content-Type", "application/json")
		base64Token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, newPassword)))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Token))
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		firstName = newFirstName
		lastName = newLastName
		password = newPassword
	})

	t.Run("TestGetUser", func(t *testing.T) {
		// Create a request to pass to our handler
		req, _ := http.NewRequest("GET", "/v1/user/self", nil)
		base64Token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, password)))

		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Token))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), email)
		assert.Contains(t, w.Body.String(), firstName)
		assert.Contains(t, w.Body.String(), lastName)
		assert.NotContains(t, w.Body.String(), password)

		// 401 - Unauthorized
		// wrong token
		req, _ = http.NewRequest("GET", "/v1/user/self", nil)
		base64Token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, "wrongpassword")))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Token))
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func initDatabase(mockCloudWatchModule *MockCloudWatchModule) database.BaseDatabase {
	postgres := database.NewPostgresDatabase(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			config.App.Database.TestUsername,
			config.App.Database.TestPassword,
			config.App.Database.TestHost,
			config.App.Database.TestPort,
			config.App.Database.TestName,
		),
		mockCloudWatchModule,
	)

	resetDatabase(postgres)

	if err := postgres.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return postgres
}

func waitForServerStart() error {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		_, err := http.Get("http://localhost:3000/v1/user")
		if err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not started after %d attempts", maxAttempts)
}

func resetDatabase(db database.BaseDatabase) {
	// Reset database
	db.GetConnection().Exec("TRUNCATE TABLE users")
}
