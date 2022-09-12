package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tashima42/shared-expenses-manager-backend/data"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	clearTable()
	ensureTableExists()
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected code %d. Got %d\n", expected, actual)
	}
}

func ensureTableExists() {
	b, _ := ioutil.ReadFile("./schema.sql")
	tableCreationQuery := string(b)
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM transactions;")
	a.DB.Exec("DELETE FROM bucket_user_account;")
	a.DB.Exec("DELETE FROM buckets;")
	a.DB.Exec("DELETE FROM user_accounts;")
	a.DB.Exec("ALTER SEQUENCE transactions_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE bucket_user_account_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE buckets_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE user_accounts_id_seq RESTART WITH 1")
}

func populateDatabaseWithUserAccount() data.UserAccount {
	uc := data.UserAccount{
		Username: "user1",
		Email:    "user1@example.com",
		Name:     "User One",
		City:     "Anta Gorda",
		PixKey:   "1939485-1kjdjfu1-lsidiri1-kyv7829",
		Role:     "admin",
	}
	uc.CreateUserAccount(a.DB)
	return uc
}
