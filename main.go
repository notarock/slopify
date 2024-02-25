package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/notarock/slopify/cmd"
)

func main() {
	cmd.Execute()
}
