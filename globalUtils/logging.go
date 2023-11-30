package globalUtils

import (
	"errors"
	"log"
)

func CreateError(message string, logger *log.Logger) error {
	const colorRed = "\033[0;31m"
	const colorNone = "\033[0m"
	logger.Printf("%sError: %s%s", colorRed, message, colorNone)
	return errors.New(message)
}
