package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tashima42/shared-expenses-manager-backend/data"
	"github.com/tashima42/shared-expenses-manager-backend/helpers"
)

type GroupHandler struct {
	DB *sql.DB
}

func (gh *GroupHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var group data.Group
	err := decoder.Decode(&group)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "CREATE-GROUP-INVALID-BODY", "Unable to parse request body")
		return
	}

	// TODO: validate params
	err = group.CreateGroup(gh.DB)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "CREATE-GROUP-FAILED", "Unable to save group")
		return
	}
	// TODO: add user to group as admin
	helpers.RespondWithJSON(w, http.StatusOK, group)
}

type addUserToGroupDTO struct {
	GroupId int    `json:"groupId"`
	UserId  int    `json:"userId"`
	Role    string `json:"role"`
}

func (gh *GroupHandler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userBucket addUserToGroupDTO
	err := decoder.Decode(&userBucket)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "ADD-USER-GROUP-INVALID-BODY", err.Error())
		return
	}
	if userBucket.Role != "admin" && userBucket.Role != "member" {
		helpers.RespondWithError(w, http.StatusBadRequest, "ADD-USER-GROUP-INVALID-ROLE", "Role must be admin or member")
		return
	}

	group := data.Group{ID: userBucket.GroupId}
	err = group.GetById(gh.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "ADD-USER-GROUP-FAILED-GET-GROUP", err.Error())
		return
	}
	user := data.UserAccount{ID: userBucket.UserId}
	err = user.GetById(gh.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "ADD-USER-GROUP-FAILED-GET-USER", err.Error())
		return
	}

	err = group.AddUserToGroup(gh.DB, userBucket.UserId, userBucket.Role)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "ADD-USER-GROUP-FAILED", err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, userBucket)
}
