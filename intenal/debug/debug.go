package debug

import (
	"github.com/davecgh/go-spew/spew"
	"io"
)

// =============================================================================
// Singleton logger
// =============================================================================

var logWrite io.Writer

func SetupLogger(w io.Writer) {
	if logWrite != nil {
		panic("Logger already initialized")
	}
	logWrite = w
}

type Level int

const(
	None Level = iota
	V
	VV
	VVV
)
var DEBUG Level = V


func DebugDump(level Level, text string, objs ...any) {
	if level == None {
		return
	}
	if level <= DEBUG {
		logWrite.Write([]byte(text))
		logWrite.Write([]byte("\n"))
		spew.Fdump(logWrite, objs...)
		logWrite.Write([]byte("\n"))
	}
}