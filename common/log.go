package common

import (
	logrus "github.com/sirupsen/logrus"
)

const componentField = "component"

// NewLogger returns a new logger instance with the configured component field.
func NewLogger(component string) *logrus.Entry {
	return logrus.WithField(componentField, component)
}
