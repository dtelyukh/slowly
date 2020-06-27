package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
}
