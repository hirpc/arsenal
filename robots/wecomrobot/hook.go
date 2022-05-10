package wecomrobot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (w wecomrobot) Establish() error {
	return nil
}

func (w wecomrobot) Fire(entry *logrus.Entry) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	return w.Send(
		ctx,
		fmt.Sprintf(
			`{"msgtype":"markdown","markdown":{"content":"%v"}}`,
			toContent(entry),
		),
	)
}

func toContent(entry *logrus.Entry) string {
	var s strings.Builder
	s.WriteString("System Monitor\n")
	s.WriteString("**************\n")
	s.WriteString("level:  " + entry.Level.String() + "\n")
	for k, val := range entry.Data {
		if v, ok := val.(string); ok {
			s.WriteString(k + ": " + v + "\n")
		}
	}
	s.WriteString("**************\n")
	s.WriteString(entry.Message)
	return s.String()
}

func (w wecomrobot) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}
