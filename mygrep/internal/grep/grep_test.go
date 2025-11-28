package grep

import (
	"testing"

	"mygrep/internal/config"
)

func TestGrepFromData(t *testing.T) {
	conf := &config.Config{
		Pattern: "error",
		Flags: &config.Flags{
			LineNumberFlag: true,
		},
		Buffer: make([]config.Line, 0),
	}

	data := "line1\nerror line2\nline3\nerror line4"
	lines, count, err := GrepFromData(data, conf)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Учитываем, что между группами может быть разделитель --
	expected := []string{"2:error line2", "--", "4:error line4"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, exp := range expected {
		if i >= len(lines) {
			t.Errorf("Missing line at index %d", i)
			continue
		}
		if lines[i] != exp {
			t.Errorf("Expected %s, got %s", exp, lines[i])
		}
	}
}
