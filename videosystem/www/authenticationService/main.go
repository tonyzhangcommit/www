package main

import (
	"auth/bootstrap"
)

func main() {
	bootstrap.InitializeConfig()
	bootstrap.RunServer()
}
