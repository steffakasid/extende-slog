package eslog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAnytoString(t *testing.T) {
	tblTest := map[string]struct {
		anything []any
		expected []string
	}{
		"success int": {
			anything: []any{int(1)},
			expected: []string{"1"},
		},
		"success float": {
			anything: []any{float32(1.4)},
			expected: []string{"1.4"},
		},
		"success bool": {
			anything: []any{bool(true)},
			expected: []string{"true"},
		},
		"success string": {
			anything: []any{string("test")},
			expected: []string{"test"},
		},
		"success struct": {
			anything: []any{struct {
				key string
			}{
				key: "test",
			}},
			expected: []string{"{test}"},
		},
		"success error": {
			anything: []any{errors.New("Error")},
			expected: []string{"Error"},
		},
	}
	for name, test := range tblTest {
		t.Run(name, func(t *testing.T) {
			testResult := convertAnyToString(test.anything...)
			assert.Equal(t, test.expected, testResult)
		})
	}
}
