package main

import (
	"go-std/internal/utils"
)

func main() {
	logger := utils.NewLogger(utils.DEBUG, true)
	logger.Info("Hello, World!")

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	logger.SetLevel(utils.INFO)
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	logger.SetLevel(utils.WARN)
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

}
