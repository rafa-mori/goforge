package main

import (
	gl "github.com/faelmori/goforge/logger"
	l "github.com/faelmori/logz"
)

// This file is the entry point for the GoForge CLI application.
// It initializes the logger and starts the application by executing the main command.
// It allows the application to be run as a standalone CLI tool.

var logger l.Logger

// main initializes the logger and creates a new GoBE instance.
func main() {
	if err := RegX().Command().Execute(); err != nil {
		gl.Log("fatal", err.Error())
	}
}
