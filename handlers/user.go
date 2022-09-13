package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tashima42/shared-expenses-manager-backend/data"
	"github.com/tashima42/shared-expenses-manager-backend/helpers"
)

type UserHandler struct {
	DB *sql.DB
}

type CreateUserResponseDTO struct {
	Success bool `json:"success"`
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user data.UserAccount
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "CREATE-USER-INVALID-BODY", "Unable to parse request body")
		return
	}

	// TODO: validate params?
	err = user.CreateUserAccount(uh.DB)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "CREATE-USER-FAILED", "Unable to save user")
		return
	}

	// TODO: return created data?
	createUserResponse := CreateUserResponseDTO{
		Success: true,
	}
	helpers.RespondWithJSON(w, http.StatusOK, createUserResponse)
}
