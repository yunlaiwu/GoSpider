package logger

import (
	"fmt"
	"runtime"
	"strings"
	"strconv"
)

func callStack() (msg string) {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		msg += fmt.Sprintf("frame = %v, file = %v, line = %v, func = %v\n", skip, file, line, f.Name())
	}
	return msg
}

/*
 * 获取当前的goroutine的id，这个俺也是从网上抄的
 */
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}