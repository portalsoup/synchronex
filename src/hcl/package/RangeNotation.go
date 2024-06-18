package _package

import (
	"fmt"
	"github.com/blang/semver/v4"
	"strings"
)

// parseRange parses the Maven range notation and returns the corresponding semver.Range function.
func parseRange(rangeStr string) (semver.Range, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	var constraints []semver.Range

	// Split the range string by comma to handle multiple conditions
	parts := strings.Split(rangeStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		var constraint semver.Range
		var err error

		switch {
		case strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]"):
			versions := strings.Split(part[1:len(part)-1], ",")
			if len(versions) == 2 {
				constraint, err = semver.ParseRange(fmt.Sprintf(">=%s <=%s", versions[0], versions[1]))
			} else {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
		case strings.HasPrefix(part, "[") && strings.HasSuffix(part, ")"):
			versions := strings.Split(part[1:len(part)-1], ",")
			if len(versions) == 2 {
				constraint, err = semver.ParseRange(fmt.Sprintf(">=%s <%s", versions[0], versions[1]))
			} else {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
		case strings.HasPrefix(part, "(") && strings.HasSuffix(part, "]"):
			versions := strings.Split(part[1:len(part)-1], ",")
			if len(versions) == 2 {
				constraint, err = semver.ParseRange(fmt.Sprintf(">%s <=%s", versions[0], versions[1]))
			} else {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
		case strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")"):
			versions := strings.Split(part[1:len(part)-1], ",")
			if len(versions) == 2 {
				constraint, err = semver.ParseRange(fmt.Sprintf(">%s <%s", versions[0], versions[1]))
			} else {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
		default:
			return nil, fmt.Errorf("invalid range format: %s", part)
		}

		if err != nil {
			return nil, err
		}

		constraints = append(constraints, constraint)
	}

	// Combine all constraints using semver.RangeAnd
	return semver.RangeAnd(constraints...), nil
}

// isVersionInRange checks if a given version is within the specified Maven range.
func isVersionInRange(version, rangeStr string) (bool, error) {
	v, err := semver.Parse(version)
	if err != nil {
		return false, err
	}

	r, err := parseRange(rangeStr)
	if err != nil {
		return false, err
	}

	return r(v), nil
}

func main() {
	version := "1.2.3"
	rangeStr := "[1.0.0,2.0.0)"

	inRange, err := isVersionInRange(version, rangeStr)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Version %s in range %s: %v\n", version, rangeStr, inRange)
	}
}
