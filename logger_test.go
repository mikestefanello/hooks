package hooks

import (
	"fmt"
	"testing"
)

func TestSetLogger(t *testing.T) {
	// No logger
	logf("test %d %d", 1, 2)

	// Add a logger that stores the output
	var output string
	SetLogger(func(format string, args ...any) {
		output = fmt.Sprintf(format, args...)
	})

	// Log
	logf("test %d %d", 1, 2)

	// Verify the output
	if output != "test 1 2" {
		t.Fail()
	}

	// Remove the logger
	SetLogger(nil)
}
