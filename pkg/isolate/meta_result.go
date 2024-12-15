package isolate

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MetaResult struct {
	Time         float64
	TimeWall     float64
	MaxRSS       int
	CSWVoluntary int
	CSWForced    int
	CGMem        int
	CGOMMKilled  int
	ExitCode     int
	Status       string
	Message      string
}

func NewMetaResultFromFile(file string) (*MetaResult, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Initialize the result
	result := &MetaResult{}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(f)

	// Process each line
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Split the line into key and value
		parts := strings.Split(strings.TrimPrefix(line, "// "), ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Parse the value based on the key
		switch key {
		case "time":
			result.Time, _ = strconv.ParseFloat(value, 64)
		case "time-wall":
			result.TimeWall, _ = strconv.ParseFloat(value, 64)
		case "max-rss":
			result.MaxRSS, _ = strconv.Atoi(value)
		case "csw-voluntary":
			result.CSWVoluntary, _ = strconv.Atoi(value)
		case "csw-forced":
			result.CSWForced, _ = strconv.Atoi(value)
		case "cg-mem":
			result.CGMem, _ = strconv.Atoi(value)
		case "cg-oom-killed":
			result.CGOMMKilled, _ = strconv.Atoi(value)
		case "exitcode":
			result.ExitCode, _ = strconv.Atoi(value)
		case "status":
			result.Status = value
		case "message":
			result.Message = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
