package src

import (
	"strconv"
	"strings"
)

const (
	Duration90   Duration = "90"
	Duration120  Duration = "120"
	Duration120p Duration = "120p"
)

func parseDuration(duration string) Duration {
	switch duration {
	case "90":
		return Duration90
	case "120":
		return Duration120
	case "120p", "120+":
		return Duration120p
	default:
		return Duration120p
	}
}

func stringToMinutes(s string) int {
	s = strings.ToLower(strings.TrimSpace(s))

	if strings.Contains(s, "h") {
		parts := strings.Split(s, "h")
		hours, _ := strconv.Atoi(parts[0])
		minutes := 0
		if len(parts) > 1 && parts[1] != "" {
			m, _ := strconv.Atoi(parts[1])
			minutes = m
		}
		return hours*60 + minutes
	}

	if strings.Contains(s, "min") {
		m, _ := strconv.Atoi(strings.ReplaceAll(s, "min", ""))
		return m
	}

	m, _ := strconv.Atoi(s)
	return m
}

func parseFilmDuration(film *Film) Duration {
	minutes := stringToMinutes(film.Duration)

	if minutes <= 100 {
		return Duration90
	} else if minutes <= 140 {
		return Duration120
	} else {
		return Duration120p
	}
}

func filterFilmsByDuration(films []*Film, target Duration) []*Film {
	var result []*Film
	for _, f := range films {
		filmDur := parseFilmDuration(f)

		if shouldIncludeDuration(filmDur, target) {
			result = append(result, f)
		}
	}
	return result
}

func shouldIncludeDuration(filmDur, target Duration) bool {
	switch target {
	case Duration90:
		return filmDur == Duration90
	case Duration120:
		return filmDur == Duration90 || filmDur == Duration120
	case Duration120p:
		return true
	default:
		return false
	}
}
