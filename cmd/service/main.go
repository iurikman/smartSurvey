package main

import "github.com/iurikman/smartSurvey/cmd/internal"

func main() {
	serverOne := server.NewServer("8080")
	serverOne.Start()

}
