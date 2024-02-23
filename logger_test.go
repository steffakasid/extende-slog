package extendedslog_test

import (
	"errors"
	"io"
	"os"
	"testing"

	extendedslog "github.com/steffakasid/extended-slog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorf(t *testing.T) {

	tblTest := map[string]struct {
		format   string
		err      error
		args     []any
		expected string
	}{
		"success": {
			format:   "error %s",
			err:      errors.New("something"),
			args:     []any{},
			expected: "error something",
		},
		"nil error": {
			format:   "error %s",
			err:      nil,
			args:     []any{},
			expected: "",
		},
		"nil error, with args": {
			format:   "error %s",
			err:      nil,
			args:     []any{},
			expected: "",
		},
		"success with args": {
			format:   "error %s %s %s",
			err:      errors.New("test"),
			args:     []any{"test1", "test2"},
			expected: "error test test1 test2",
		},
		"success with args positions": {
			format:   "error %[2]s %[1]s %[3]s",
			err:      errors.New("test"),
			args:     []any{"test1", "test2"},
			expected: "error test1 test test2",
		},
	}

	for name, tst := range tblTest {
		t.Run(name, func(t *testing.T) {
			r, w, err := os.Pipe()
			require.NoError(t, err)

			t.Cleanup(func(t *testing.T) func() {
				extendedslog.Logger.SetOutput(w)

				return func() {
					extendedslog.Logger.SetOutput(os.Stdout)
				}
			}(t))

			extendedslog.Logger.Errorf(tst.format, tst.err, tst.args...)

			w.Close()

			out, err := io.ReadAll(r)
			require.NoError(t, err)

			assert.Contains(t, string(out), tst.expected)
		})
	}
}
