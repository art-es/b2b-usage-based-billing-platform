package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type logfield struct {
	key string
	val string
}

type logger struct {
	level        Level
	format       Format
	output       io.Writer
	getCreatedAt func() string
	fields       []logfield
}

type log struct {
	logger *logger
	level  Level
	fields []logfield
}

func (l logger) Set(key, val string) Logger {
	if l.level == Disabled {
		return l
	}

	newFields := make([]logfield, len(l.fields)+1)
	copy(newFields, l.fields)
	newFields[len(newFields)-1] = logfield{key: key, val: val}
	l.fields = newFields
	return l
}

func (l logger) Log(lvl Level) Log {
	return log{
		logger: &l,
		level:  lvl,
	}
}

func (l log) Set(key string, val string) Log {
	if l.logger.level == Disabled || l.level > l.logger.level {
		return l
	}

	newFields := make([]logfield, len(l.fields)+1)
	copy(newFields, l.fields)
	newFields[len(newFields)-1] = logfield{key: key, val: val}
	l.fields = newFields
	return l
}

func (l log) Write() {
	if l.logger.level == Disabled || l.level > l.logger.level {
		return
	}

	fieldsMap := make(map[string]string, len(l.logger.fields)+len(l.fields)+2)

	for _, f := range l.logger.fields {
		fieldsMap[f.key] = f.val
	}

	for _, f := range l.fields {
		fieldsMap[f.key] = f.val
	}

	fieldsMap["level"] = l.level.String()
	fieldsMap["created_at"] = l.logger.getCreatedAt()

	switch l.logger.format {
	case FormatText:
		writeAsText(l.logger.output, fieldsMap)
	default:
		writeAsJSON(l.logger.output, fieldsMap)
	}
}

func writeAsJSON(w io.Writer, data map[string]string) {
	json.NewEncoder(w).Encode(data)
}

func writeAsText(w io.Writer, data map[string]string) {
	line := fmt.Sprintf("%s %s:", strings.ToUpper(data["level"]), data["created_at"])
	for key, val := range data {
		line += " " + key + "=" + val
	}

	w.Write([]byte(line + "\n"))
}
