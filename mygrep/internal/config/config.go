package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type NodeType string

const (
	NodeTypeCoordinator NodeType = "coordinator"
	NodeTypeWorker      NodeType = "worker"
)

type Config struct {
	NodeType NodeType
	Workers  []string
	Quorum   int
	Port     string

	Pattern string
	Flags   *Flags

	LineNum       int
	MatchesCount  int
	PostContext   int
	LastMatchLine int
	Buffer        []Line
}

type Line struct {
	Num   int
	Text  string
	Match bool
}

type Flags struct {
	IgnoreFlag       bool
	InvertFlag       bool
	StrictStringFlag bool
	LineNumberFlag   bool
	CountOfLineFlag  bool

	AdditionalLineFlag int // -A N
	PreviousLineFlag   int // -B N
	AroundLineFlag     int // -C N
}

func ParseConfig() (*Config, error) {
	args := os.Args[1:]
	var pattern string
	var remainingArgs []string

	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			pattern = arg
			remainingArgs = append(remainingArgs, args[:i]...)
			remainingArgs = append(remainingArgs, args[i+1:]...)
			break
		}
	}

	os.Args = []string{"mygrep"}
	os.Args = append(os.Args, remainingArgs...)

	var nodeTypeStr string = "coordinator"
	var workersStr string
	var port string = "8080"
	var quorum int = 1

	flags := &Flags{}

	flag.StringVar(&nodeTypeStr, "mode", "coordinator", "Node type: coordinator or worker")
	flag.StringVar(&workersStr, "workers", "", "List of workers: host1:port1,host2:port2")
	flag.IntVar(&quorum, "quorum", 1, "Quorum: min number of successful responses")
	flag.StringVar(&port, "port", "8080", "Port for worker mode")
	flag.BoolVar(&flags.IgnoreFlag, "i", false, "ignore register")
	flag.BoolVar(&flags.InvertFlag, "v", false, "invert order")
	flag.BoolVar(&flags.StrictStringFlag, "F", false, "setting fix string")
	flag.BoolVar(&flags.LineNumberFlag, "n", false, "show numbers of string per each")
	flag.BoolVar(&flags.CountOfLineFlag, "c", false, "show count of matching strings with pattern")
	flag.IntVar(&flags.AdditionalLineFlag, "A", 0, "after each find string additional show n string after")
	flag.IntVar(&flags.PreviousLineFlag, "B", 0, "before each find string show previous n strings")
	flag.IntVar(&flags.AroundLineFlag, "C", 0, "show before and after n strings")

	flag.Parse()
	nodeType := NodeType(nodeTypeStr)
	if nodeType != NodeTypeCoordinator && nodeType != NodeTypeWorker {
		return nil, fmt.Errorf("invalid mode: %s", nodeTypeStr)
	}
	var workers []string
	if workersStr != "" {
		workers = strings.Split(workersStr, ",")
	}

	if flags.AroundLineFlag > 0 {
		flags.AdditionalLineFlag = flags.AroundLineFlag
		flags.PreviousLineFlag = flags.AroundLineFlag
	}

	return &Config{
		NodeType: nodeType,
		Workers:  workers,
		Quorum:   quorum,
		Port:     port,
		Pattern:  pattern,
		Flags:    flags,
		Buffer:   make([]Line, 0),
	}, nil
}

func (c *Config) Matcher() (func(string) bool, error) {
	if !c.Flags.StrictStringFlag {
		pattern := c.Pattern
		if c.Flags.IgnoreFlag {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		return func(s string) bool {
			match := re.MatchString(s)
			if c.Flags.InvertFlag {
				return !match
			}
			return match
		}, nil
	} else {
		str := c.Pattern
		if c.Flags.IgnoreFlag {
			str = strings.ToLower(str)
		}
		return func(s string) bool {
			tmp := s
			if c.Flags.IgnoreFlag {
				tmp = strings.ToLower(tmp)
			}
			match := strings.Contains(tmp, str)
			if c.Flags.InvertFlag {
				return !match
			}
			return match
		}, nil
	}
}
