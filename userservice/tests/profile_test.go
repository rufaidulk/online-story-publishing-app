package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"userservice/controllers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileFetchApi(t *testing.T) {
	token := fmt.Sprintf("Bearer %s", jwt)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")
	ctx.Set("userDetails", userDetails)

	if assert.NoError(t, controllers.ViewUserProfile(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUserProfileUpdateApi(t *testing.T) {
	token := fmt.Sprintf("Bearer %s", jwt)
	reqBody := fmt.Sprintf("{\"name\":\"%s\",\"pen_name\":\"%s\",\"password\": \"%s\",\"description\":\"%s\"}", testUser.Name, testUser.PenName, testUser.Password, testUser.Description)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")
	ctx.Set("userDetails", userDetails)

	if assert.NoError(t, controllers.ViewUserProfile(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
