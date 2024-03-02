package eslog_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/steffakasid/eslog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorf(t *testing.T) {

	tblTest := map[string]struct {
		format   string
		args     []any
		expected string
	}{
		"success": {
			format:   "error %s",
			args:     []any{errors.New("something")},
			expected: "error something",
		},
		"nil error": {
			format:   "error %s",
			args:     []any{},
			expected: "",
		},
		"nil error with args": {
			format:   "error %s",
			args:     []any{},
			expected: "",
		},
		"success with args": {
			format:   "error %s %s %s",
			args:     []any{errors.New("test"), "test1", "test2"},
			expected: "error test test1 test2",
		},
		"success with args positions": {
			format:   "error %[2]s %[1]s %[3]s",
			args:     []any{errors.New("test"), "test1", "test2"},
			expected: "error test1 test test2",
		},
	}

	for name, tst := range tblTest {
		t.Run(name, func(t *testing.T) {
			r, w, err := os.Pipe()
			require.NoError(t, err)

			t.Cleanup(func(t *testing.T) func() {
				eslog.Logger.SetOutput(w)

				return func() {
					eslog.Logger.SetOutput(os.Stdout)
				}
			}(t))

			eslog.Logger.Errorf(tst.format, tst.args...)

			w.Close()

			out, err := io.ReadAll(r)
			require.NoError(t, err)

			assert.Contains(t, string(out), tst.expected)
		})
	}
}
