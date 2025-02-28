package test

import (
	"strings"
)

func filterRoutineWASM(stack string) bool {
	// Nested t.Run on Go 1.14-1.21 and go1.22 WASM have these routines
	return strings.Contains(stack, "runtime.goexit()") ||
		strings.Contains(stack, "runtime.goexit({})")
}
