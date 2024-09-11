package utils

import (
	"os"
)

func GetGoEnv() string {
	e, ok := os.LookupEnv("ENV")
	if !ok {
		e = defaultGoEnv
	}
	return e
}

func GetSwagHostEnv() string {
	e, ok := os.LookupEnv("SWAG_URL")
	if !ok {
		e = defaultSwagUrl
	}
	return e
}

const defaultSwagUrl = "localhost:4001/message-api"
const defaultGoEnv = "qa"
