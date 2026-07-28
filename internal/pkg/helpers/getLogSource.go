package helpers

import (
	"fmt"
	"runtime"
)

func GetLogSource() string {
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		return fmt.Sprintf("%s (%s:%d)", funcName, file, line)
	}
	return "undefined"
}
