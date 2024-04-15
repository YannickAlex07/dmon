package cli

import (
	"fmt"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

// custom pterm hook for logrus
type PTermHook struct {
	Verbose bool
}

func (hook *PTermHook) entryToMessage(entry *log.Entry) string {
	msg := entry.Message

	if hook.Verbose {
		if len(entry.Data) > 0 {
			msg += " >>"
		}

		for k, v := range entry.Data {
			msg += fmt.Sprintf(" %s=\"%v\"", k, v)
		}
	}

	return msg
}

func (hook *PTermHook) Fire(entry *log.Entry) error {
	pterm.PrintDebugMessages = true
	msg := hook.entryToMessage(entry)

	switch entry.Level {
	case log.ErrorLevel:
		pterm.Error.Println(msg)
	case log.WarnLevel:
		pterm.Warning.Println(msg)
	case log.DebugLevel:
		pterm.Debug.Println(msg)
	default:
		pterm.Info.Println(msg)
	}

	return nil
}

func (hook *PTermHook) Levels() []log.Level {
	levels := []log.Level{
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
	}

	if hook.Verbose {
		levels = append(levels, log.DebugLevel)
	}

	return levels
}
