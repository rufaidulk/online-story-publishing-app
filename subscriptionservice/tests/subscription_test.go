package tests

import (
	"fmt"
	"os"
	"testing"
)

var jwt string

func TestMain(m *testing.M) {
	fmt.Println("Unit tests starting...")
	exitCode := m.Run()
	fmt.Println("Unit tests completed")
	os.Exit(exitCode)
}

//todo:: middleware is not calling, R&D required
// func TestBasicSubscriptionCreation(t *testing.T) {
// 	// var userUuid interface{}
// 	userUuid := uuid.New().String()
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set("useruuid", userUuid)
// 	rec := httptest.NewRecorder()
// 	ctx := e.NewContext(req, rec)
// 	ctx.Set("mode", "testing")

// 	if assert.NoError(t, controllers.CreateBaseSubscription(ctx)) {
// 		assert.Equal(t, http.StatusCreated, rec.Code)
// 	}
// 	fmt.Println("Response:")
// 	fmt.Println(rec.Body.String())
// }
