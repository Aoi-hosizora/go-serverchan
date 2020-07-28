package serverchan

import (
	"log"
	"strings"
)

// Serverchan logger, please override `Log`.
type Logger interface {
	// sckey and title may need to be masked.
	// if code == -1, it is a http error or unmarshal error.
	// else if code != -1, it is a response error (maybe 1024).
	Log(sckey string, title string, code int32, err error)
}

// Serverchan DefaultLogger's log-mode, use this to control DefaultLogger behavior.
type DefaultLogMode uint8

const (
	// Not to log.
	LMNone DefaultLogMode = iota

	// Log error only.
	LMErr

	// Log all, but with masked sckey and title.
	LMMask

	// Log all, include full sckey.
	LMAll
)

type defaultLogger struct {
	mode DefaultLogMode
}

// A serverchan default logger with log.Logger.
func DefaultLogger(mode DefaultLogMode) *defaultLogger {
	return &defaultLogger{mode: mode}
}

func (d *defaultLogger) Log(sckey string, title string, code int32, err error) {
	mode := d.mode
	if mode <= LMNone {
		return
	}

	if mode >= LMErr {
		// Err, Mask, All
		if err != nil {
			log.Printf("[Serverchan] failed to send message to %s: %v", Mask(sckey), err)
			return
		}
	}

	if mode >= LMMask {
		// Mask, All
		if mode == LMMask {
			sckey = Mask(sckey)
			title = Mask(title)
		}

		if code == ErrnoSuccess {
			log.Printf("[Serverchan] <- %3d | %s | %s", 0, sckey, title)
		}
	}
}

type noLogger struct{}

// A serverchan logger that not to do anything.
func NoLogger() *noLogger {
	return &noLogger{}
}

// noinspection GoUnusedParameter
func (n *noLogger) Log(sckey string, title string, code int32, err error) {}

// Mask sckey as `*`.
func Mask(tok string) string {
	switch len(tok) {
	case 0:
		return ""
	case 1:
		return "*"
	case 2:
		return "*" + tok[1:]
	case 3:
		return "**" + tok[2:3]
	case 4:
		return tok[0:1] + "**" + tok[3:4]
	case 5:
		return tok[0:1] + "***" + tok[4:5]
	default:
		return tok[0:2] + strings.Repeat("*", len(tok)-4) + tok[len(tok)-2:] // <<< Default
	}
}
