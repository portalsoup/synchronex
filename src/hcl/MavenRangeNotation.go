package hcl

import (
	"fmt"
	"strings"
)

type Boundary struct {
	Version   string
	Inclusive bool
}

type Range struct {
	Start Boundary
	End   Boundary
}

func (r Range) IsInRange(version string) bool {
	if r.Start.Version != "" {
		if r.Start.Inclusive {
			if version < r.Start.Version {
				return false
			}
		} else {
			if version <= r.Start.Version {
				return false
			}
		}
	}

	if r.End.Version != "" {
		if r.End.Inclusive {
			if r.End.Version < version {
				return false
			}
		} else {
			if r.End.Version <= version {
				return false
			}
		}
	}
	return true
}

func tokenizeBoundary(boundaryStr string) (Boundary, error) {
	boundaryStr = strings.TrimSpace(boundaryStr)

	var b Boundary

	// Determine if the boundary is inclusive or exclusive
	if boundaryStr[0] == '[' || boundaryStr[0] == '(' {
		b.Inclusive = boundaryStr[0] == '['
		b.Version = boundaryStr[1:] // Remove the first character
	} else if boundaryStr[len(boundaryStr)-1] == ']' || boundaryStr[len(boundaryStr)-1] == ')' {
		b.Inclusive = boundaryStr[len(boundaryStr)-1] == ']'
		b.Version = boundaryStr[:len(boundaryStr)-1] // Remove the last character
	} else {
		return Boundary{}, fmt.Errorf("invalid boundary string: %s", boundaryStr)
	}

	return b, nil
}

// TokenizeRange parses a Maven range notation string into an array of Range structs.
func TokenizeRange(rangeStr string) ([]Range, error) {

	var ranges []Range

	// Split the range string by comma to handle multiple range specifications
	rawBoundaries := strings.Split(rangeStr, ",")

	if len(rawBoundaries)%2 != 0 {
		return []Range{}, fmt.Errorf("Invalid range string: %s", Boundary{})
	}

	for i := 0; i < len(rawBoundaries); i += 2 {
		start, err := tokenizeBoundary(rawBoundaries[i])
		if err != nil {
			return ranges, err
		}
		end, err := tokenizeBoundary(rawBoundaries[i+1])
		if err != nil {
			return ranges, err
		}
		r := Range{
			Start: start,
			End:   end,
		}
		ranges = append(ranges, r)
	}

	return ranges, nil
}
