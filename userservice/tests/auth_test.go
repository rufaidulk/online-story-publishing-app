package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"userservice/controllers"
	"userservice/helper"
	"userservice/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var jwt string
var userDetails interface{}

type TestUser struct {
	Name        string
	Email       string
	Password    string
	PenName     string
	Description string
}

var testUser = TestUser{
	Name:        "alex",
	Email:       "alex@test.com",
	Password:    "123456",
	PenName:     "alexy",
	Description: "thriller story writer",
}

func TestMain(m *testing.M) {
	fmt.Println("Unit tests starting...")
	exitCode := m.Run()
	fmt.Println("Unit tests completed")
	os.Exit(exitCode)
}

func TestRegistration(t *testing.T) {
	reqBody := fmt.Sprintf("{\"name\":\"%s\",\"email\":\"%s\",\"password\": \"%s\" }", testUser.Name, testUser.Email, testUser.Password)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")

	if assert.NoError(t, controllers.Registration(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
	// fmt.Println("Response:")
	// fmt.Println(rec.Body.String())
}

func TestLogin(t *testing.T) {
	reqBody := fmt.Sprintf("{\"email\":\"%s\",\"password\": \"%s\" }", testUser.Email, testUser.Password)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")

	if assert.NoError(t, controllers.Login(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	// fmt.Println("Response:")
	// fmt.Println(rec.Body.String())
	x := helper.SuccessResponse{}
	json.Unmarshal(rec.Body.Bytes(), &x)
	data := x.Data.(map[string]interface{})
	jwt = data["token"].(string)
}

func TestInvalidJwtOnAuthMiddleware(t *testing.T) {
	invalidJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDU2MzU0MDksImlzcyI6ImFkbWluIiwidXVpZCI6IjE0YTAxZTUzLTkyY2UtNDMwNi1hZmVmLTI1OGM0MTFhMzEzNyIsIm5hbWUiOiJhZGkiLCJlbWFpbCI6ImFkaUBtYWlsLmNvbSJ9.KzZDfVwvr6wLwkfcNZuL84Wpj4XxOazOvg2q9bdDa8g"
	token := fmt.Sprintf("Bearer %s", invalidJwt)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, token)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	mw := middlewares.UserJwtMiddleware
	if assert.NoError(t, mw(controllers.AuthorizeUser)(ctx)) {
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	}
}
func TestValidJwtOnAuthMiddleware(t *testing.T) {
	token := fmt.Sprintf("Bearer %s", jwt)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, token)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)
	mw := middlewares.UserJwtMiddleware

	if assert.NoError(t, mw(controllers.AuthorizeUser)(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		userDetails = ctx.Get("userDetails")
	}
}

func TestAuthorization(t *testing.T) {
	token := fmt.Sprintf("Bearer %s", jwt)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")
	ctx.Set("userDetails", userDetails)

	if assert.NoError(t, controllers.AuthorizeUser(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
