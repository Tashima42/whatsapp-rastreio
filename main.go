package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		os.Getenv("WHATSAPP_USER_ACCESS_TOKEN"),
		os.Getenv("WHATSAPP_BASE_URL"),
		os.Getenv("CORREIOS_BASE_URL"),
		os.Getenv("SECRET"),
		os.Getenv("GO_ENV"),
	)

	fmt.Println("Running on PORT", os.Getenv("APP_PORT"))
	a.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
