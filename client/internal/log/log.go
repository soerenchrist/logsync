package log

import (
	"fmt"
)

func Info(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Error(format string, a ...any) {
	fmt.Printf("ERROR: "+format+"%v\n", a...)
}
