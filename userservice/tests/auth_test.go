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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var jwt string

func TestMain(m *testing.M) {
	fmt.Println("Unit tests starting...")
	exitCode := m.Run()
	fmt.Println("Unit tests completed")
	os.Exit(exitCode)
}

func TestRegistration(t *testing.T) {
	reqBody := `{"name":"alex","email":"alex@test.com","password":"123456"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")

	if assert.NoError(t, controllers.Registration(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
	fmt.Println("Response:")
	fmt.Println(rec.Body.String())
}

func TestLogin(t *testing.T) {
	reqBody := `{"email":"alex@test.com","password":"123456"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("mode", "testing")

	if assert.NoError(t, controllers.Login(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	fmt.Println("Response:")
	fmt.Println(rec.Body.String())
	x := helper.SuccessResponse{}
	json.Unmarshal(rec.Body.Bytes(), &x)
	data := x.Data.(map[string]interface{})
	jwt = data["token"].(string)
}

//todo:: to be debugged
// func TestAuthorization(t *testing.T) {
// 	// token := "Bearer "
// 	// token += jwt
// 	token := fmt.Sprintf("Bearer %s", jwt)
// 	fmt.Println(token)
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, token)
// 	rec := httptest.NewRecorder()
// 	ctx := e.NewContext(req, rec)
// 	ctx.Set("mode", "testing")

// 	if assert.NoError(t, controllers.AuthorizeUser(ctx)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// 	fmt.Println("Response:")
// 	fmt.Println(rec.Body.String())
// }
