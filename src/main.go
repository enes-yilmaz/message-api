package main

import (
	"MESSAGEAPI/src/cmd/message"
	"MESSAGEAPI/src/pkg/utils"
)

func main() {
	env := utils.GetGoEnv()
	message.Execute(env)
}
