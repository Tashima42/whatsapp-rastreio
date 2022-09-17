package helpers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func NowPlusSeconds(seconds int) time.Time {
	return time.Now().Local().Add(time.Second * time.Duration(seconds))
}
func NowMinusMinutes(minutes int) time.Time {
	return time.Now().Local().Add(time.Minute * time.Duration(minutes))
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func ParseDateIso(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", date)
}

func FormatDateIso(date time.Time) string {
	return date.Format("2006-01-02T15:04:05-0700")
}

func RespondWithError(w http.ResponseWriter, code int, errorCode string, message string) {
	fmt.Printf("ErrorCode: %v, Message: %v", errorCode, message)
	RespondWithJSON(w, code, map[string]interface{}{"success": false, "errorCode": errorCode, "message": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type Logger struct {
	DB  *sql.DB
	Env string
}

func (lg *Logger) Log(log string, label string) {
	err := lg.DB.QueryRow(
		"INSERT INTO logs(log, env, date, label) VALUES($1, $2, $3, $4)",
		log,
		lg.Env,
		time.Now(),
		label,
	)
	if err != nil {
		fmt.Println(err.Err())
	}
}
