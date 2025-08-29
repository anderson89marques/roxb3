package utils

import (
	"fmt"
	"time"
)

func ParseStringHourMilisToTime(timeString string) (*time.Time, error) {
	// Converts strings in this layout HHMMSSmmm to Time
	if len(timeString) != 9 {
		return nil, fmt.Errorf("invalid time string format: Expected HHMMSSmmm receive: %s", timeString)
	}
	secondsPart := timeString[:6]
	millisecondsPart := timeString[6:]

	// Parse the HHMMSS part
	parsedTime, err := time.Parse("150405", secondsPart)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %w", err)
	}

	// Convert milliseconds string to int
	var milliseconds int
	_, err = fmt.Sscanf(millisecondsPart, "%3d", &milliseconds)
	if err != nil {
		return nil, fmt.Errorf("error converting milliseconds: %w", err)
	}

	// Add milliseconds to the parsed time
	finalTime := parsedTime.Add(time.Duration(milliseconds) * time.Millisecond)
	return &finalTime, nil
}
