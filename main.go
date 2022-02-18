package main

import (
	"os"

	"sourcability.com/asha-bot/ashabot"
)

func main() {
	os.Exit(ashabot.CLI(os.Args[1:]))
}
