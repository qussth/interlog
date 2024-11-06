package interlog

import (
	"testing"
)

// Should return
// 2024-11-06 13:49:03 | INFO  | info   1:test_value_1=1   2:test_value_2=2
// 2024-11-06 13:49:03 | ERROR |   1:test_value_1=1   2:test_value_2=2   3:test_value_3=3   4:test_value_4=4
func TestLogger_With(t *testing.T) {
	l := New()
	c := l.With([]Value{
		{"test_value_1", "1"},
		{"test_value_2", "2"},
	})

	c.Info("info", nil)
	c.Error(nil, Values{
		{"test_value_3", "3"},
		{"test_value_4", "4"},
	})
}
