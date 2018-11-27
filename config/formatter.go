package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	MEGABYTE     = 1.0
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

var (
	bytesPattern = regexp.MustCompile(`(?i)^(-?\d+)([KMGT])B?$`)
)

func invalidByteQuantityError() error {
	return errors.New("Byte quantity must be an integer with a unit of measurement like M, MB, G, or GB")
}
