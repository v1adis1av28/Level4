package grep

import (
	"bufio"
	"fmt"
	"mygrep/internal/config"
	"strings"
)

func GrepFromData(data string, c *config.Config) ([]string, int, error) {
	lines := []string{}
	count := 0

	matcher, err := c.Matcher()
	if err != nil {
		return nil, -1, err
	}

	scanner := bufio.NewScanner(strings.NewReader(data))
	firstMatch := true

	for scanner.Scan() {
		c.LineNum++
		line := scanner.Text()
		match := matcher(line)

		c.Buffer = append(c.Buffer, config.Line{Num: c.LineNum, Text: line, Match: match})

		if len(c.Buffer) > c.Flags.PreviousLineFlag+1 {
			c.Buffer = c.Buffer[1:]
		}

		if match {
			count++
			if c.Flags.CountOfLineFlag {
				continue
			}

			if !firstMatch && c.LastMatchLine+c.Flags.AdditionalLineFlag+1 < c.LineNum {
				lines = append(lines, "--")
			}
			firstMatch = false
			c.LastMatchLine = c.LineNum

			for i := 0; i < len(c.Buffer)-1; i++ {
				output := formatLine(c, c.Buffer[i], true)
				lines = append(lines, output)
			}

			output := formatLine(c, c.Buffer[len(c.Buffer)-1], false)
			lines = append(lines, output)

			c.PostContext = c.Flags.AdditionalLineFlag

		} else if c.PostContext > 0 {
			output := formatLine(c, config.Line{Num: c.LineNum, Text: line, Match: false}, true)
			lines = append(lines, output)
			c.PostContext--
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, -1, err
	}

	if c.Flags.CountOfLineFlag {
		return []string{fmt.Sprintf("%d", c.MatchesCount)}, c.MatchesCount, nil
	}

	return lines, count, nil
}

func formatLine(c *config.Config, line config.Line, isContext bool) string {
	if c.Flags.CountOfLineFlag {
		return ""
	}

	prefix := ":"
	if isContext {
		prefix = "-"
	}

	if c.Flags.LineNumberFlag {
		return fmt.Sprintf("%d%s%s", line.Num, prefix, line.Text)
	}
	return line.Text
}
