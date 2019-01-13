package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MEGABYTE = 1.0
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

func ByteSize(bytes int) string {
	if bytes == -1 {
		return "unlimited"
	}
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes == 0:
		return "0"
	case bytes < GIGABYTE:
		unit = "M"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

func StringToMegabytes(s string) (string, error) {
	i, err := ToMegabytes(s)
	if err != nil {
		return "", fmt.Errorf("must be an integer instead of [%s]", s)
	}
	if i == -1 {
		return "unlimited", nil
	}
	return ByteSize(i), nil
}

func AsString(i int) string {
	if i == -1 {
		return "unlimited"
	}
	return strconv.Itoa(i)
}

func ToInteger(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	if strings.EqualFold(s, "unlimited") {
		return -1, nil
	}
	return strconv.Atoi(strings.TrimSpace(s))
}

func ToMegabytes(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	if strings.EqualFold(s, "unlimited") {
		return -1, nil
	}
	value, err := strconv.Atoi(strings.TrimSpace(s))
	if err == nil {
		return value, nil
	}
	parts := bytesPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(parts) < 3 {
		return 0, invalidByteQuantityError()
	}

	value, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, invalidByteQuantityError()
	}

	var bytes int
	unit := strings.ToUpper(parts[2])
	switch unit {
	case "T", "TB":
		bytes = value * TERABYTE
	case "G", "GB":
		bytes = value * GIGABYTE
	case "M", "MB":
		bytes = value * MEGABYTE
	}

	return bytes / MEGABYTE, nil
}

func FutureTime(t time.Time, timeToAdd string) (string, error) {
	if timeToAdd == "" {
		return t.Format(time.RFC3339), nil
	}
	parts := timePattern.FindStringSubmatch(strings.TrimSpace(timeToAdd))
	if len(parts) < 3 {
		return "", invalidTimePatternError()
	}
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", invalidTimePatternError()
	}
	unit := strings.ToUpper(parts[2])
	switch unit {
	case "D":
		t = t.Add(time.Hour * 24 * time.Duration(value))
	case "H":
		t = t.Add(time.Hour * time.Duration(value))
	case "M":
		t = t.Add(time.Minute * time.Duration(value))
	default:
		return "", invalidTimePatternError()
	}

	return t.Format(time.RFC3339), nil
}

var (
	bytesPattern = regexp.MustCompile(`(?i)^(-?\d+)([KMGT])B?$`)
	timePattern  = regexp.MustCompile(`(?i)^(-?\d+)([DHM])$`)
)

func invalidByteQuantityError() error {
	return errors.New("Byte quantity must be an integer with a unit of measurement like M, MB, G, or GB")
}

func invalidTimePatternError() error {
	return errors.New("Time to add must have format like D, M or H")
}
