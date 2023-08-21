package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Jimbo8702/goreservation/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	// insertedUser := insertTestUser(t, tdb.User)
	insertedUser := fixtures.AddUser(tdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email: "james@foo.com",
		Password: "james_foo",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	
	if authResp.Token == "" {
		t.Fatalf("expected the JWT to be present in the auth response")
	} 
	// set the encrypted password to an empty string, because we do NOT return that in any 
	// json response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User){
		t.Fatalf("expected user to equal inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	// insertTestUser(t, tdb.User)
	fixtures.AddUser(tdb.Store, "james", "foo", false)


	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email: "james@foo.com",
		Password: "passwordnotcorrect",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}
	var errResp Error
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatal(err)
	}
	if errResp.Err != "invalid credentials" {
		t.Fatalf("expected error to be invalid credentials but got %s", errResp.Err)
	}
}