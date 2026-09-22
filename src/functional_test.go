package main

import (
	"testing"

	"github.com/dacanizares/rapidou/tst"
)

const testAdminEmail = "admin@rapidou.test"
const testAdminPassword = "correct-horse"

func testApp(t *testing.T) (*App, func()) {
	t.Helper()
	app, err := newApp(Config{
		DatabasePath:  t.TempDir() + "/rapidou.db",
		JWTSecret:     "test-secret-that-never-leaves-the-test",
		AdminEmail:    testAdminEmail,
		AdminPassword: testAdminPassword,
	})
	if err != nil {
		t.Fatal(err)
	}
	return app, func() { app.Close() }
}

func TestAPI(t *testing.T) {
	tst.RunAPI(t, func(t *testing.T) (tst.Application, func()) {
		app, close := testApp(t)
		return app.Handler(), close
	}, tst.Credentials{Email: testAdminEmail, Password: testAdminPassword})
}

func TestUI(t *testing.T) {
	tst.RunUI(t, func(t *testing.T) (tst.Application, func()) {
		app, close := testApp(t)
		return app.Handler(), close
	}, tst.Credentials{Email: testAdminEmail, Password: testAdminPassword})
}
