package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	treePrefix      = "├─ "
	treeLastPrefix  = "└─ "
	treePadding     = "│  "
	treeEmptyPading = "   "
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorBlue       = "\033[34m"
	colorReset      = "\033[0m"
)

type TreeFormatter struct {
	TimestampFormat string
	ShowCaller      bool
	UseColors       bool
}

func (f *TreeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// Write timestamp and level if present
	if entry.Time.IsZero() == false {
		b.WriteString(entry.Time.Format(f.TimestampFormat))
		b.WriteString(" ")
	}

	// Add color to level
	levelColor := ""
	if f.UseColors {
		switch entry.Level {
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = colorRed
		case logrus.WarnLevel:
			levelColor = colorYellow
		case logrus.InfoLevel:
			levelColor = colorGreen
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = colorBlue
		}
	}

	if levelColor != "" {
		b.WriteString(levelColor)
		b.WriteString(strings.ToUpper(entry.Level.String()))
		b.WriteString(colorReset)
	} else {
		b.WriteString(strings.ToUpper(entry.Level.String()))
	}

	b.WriteString(": ")

	// Write main message
	b.WriteString(entry.Message)
	b.WriteString("\n")

	// Sort fields for consistent output
	var fields []string
	for field := range entry.Data {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	// Write fields in tree format
	for i, field := range fields {
		isLast := i == len(fields)-1
		prefix := treePrefix
		if isLast {
			prefix = treeLastPrefix
		}

		value := entry.Data[field]
		b.WriteString(prefix)
		b.WriteString(fmt.Sprintf("%s: %v\n", field, value))
	}

	return b.Bytes(), nil
}
