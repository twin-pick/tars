package src

import (
	"regexp"
	"strconv"
	"strings"
)

const MAX_SHORT_DURATION = 100
const MAX_MEDIUM_DURATION = 135

func filterFilmsByDuration(films []*Film, duration string) []*Film {
	var result []*Film

	for _, f := range films {
		minutes := stringToMinutes(f.Duration)

		switch duration {
		case "short":
			if minutes > 0 && minutes <= MAX_SHORT_DURATION {
				result = append(result, f)
			}
		case "medium":
			if minutes > MAX_SHORT_DURATION && minutes <= MAX_MEDIUM_DURATION {
				result = append(result, f)
			}
		case "long":
			result = append(result, f)
		}

		if f.Duration == "" && duration == "long" {
			result = append(result, f)
		}
	}

	return result
}

func stringToMinutes(minutes string) int {
	minutes = strings.TrimSpace(minutes)
	if minutes == "" {
		return 0
	}

	reMin := regexp.MustCompile(`(\d+)\s*min`)
	if match := reMin.FindStringSubmatch(minutes); len(match) == 2 {
		if val, err := strconv.Atoi(match[1]); err == nil {
			return val
		}
	}

	if val, err := strconv.Atoi(minutes); err == nil {
		return val
	}

	return 0
}