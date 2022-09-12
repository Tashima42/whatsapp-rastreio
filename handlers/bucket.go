package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tashima42/shared-expenses-manager-backend/data"
	"github.com/tashima42/shared-expenses-manager-backend/helpers"
)

type BucketHandler struct {
	DB *sql.DB
}

type CreateBucketResponseDTO struct {
	Success bool `json:"success"`
}

func (bh *BucketHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var bucket data.Bucket
	err := decoder.Decode(&bucket)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "CREATE-BUCKET-INVALID-BODY", "Unable to parse request body")
		return
	}

	// TODO: validate params?
	err = bucket.CreateBucket(bh.DB)
	if err != nil {
		fmt.Print(err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "CREATE-BUCKET-FAILED", "Unable to save bucket")
		return
	}

	// TODO: return created data?
	createBucketResponse := CreateBucketResponseDTO{
		Success: true,
	}
	helpers.RespondWithJSON(w, http.StatusOK, createBucketResponse)
}
