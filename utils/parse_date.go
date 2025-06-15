package utils

import (
    "time"
    "errors"
)

// ParseDate parses a date string in "YYYY-MM-DD" format and returns a time.Time object.
func ParseDate(dateStr string) (time.Time, error) {
    const layout = "2006-01-02"
    t, err := time.Parse(layout, dateStr)
    if err != nil {
        return time.Time{}, errors.New("date must be in YYYY-MM-DD format")
    }
    return t, nil
}