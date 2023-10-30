package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
)

const (
	MEGABYTE = 1.0
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1000000 * MEGABYTE
)

var (
	bytesPattern = regexp.MustCompile(`(?i)^(-?\d+\.?\d*)([KMGT])B?$`)
	timePattern  = regexp.MustCompile(`(?i)^(-?\d+)([DHM])$`)
)

func ByteSize(bytes *int) string {
	if bytes == nil {
		return "unlimited"
	}
	unit := ""
	value := float32(*bytes)
	switch {
	case *bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case *bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case *bytes == 0:
		return "0"
	case *bytes < GIGABYTE:
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
	if i == nil {
		return "unlimited", nil
	}
	return ByteSize(i), nil
}

func AsString(i *int) string {
	if i == nil {
		return "unlimited"
	}
	return strconv.Itoa(*i)
}

func ToInteger(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	if strings.EqualFold(s, "unlimited") {
		return nil, nil
	}
	if strings.EqualFold(s, "-1") {
		return nil, nil
	}
	val, err := strconv.Atoi(strings.TrimSpace(s))
	if err == nil {
		return util.GetIntPointer(val), nil
	}
	return nil, err
}

func ToMegabytes(s string) (*int, error) {
	if s == "" {
		return util.GetIntPointer(0), nil
	}
	if strings.EqualFold(s, "unlimited") {
		return nil, nil
	}
	if strings.EqualFold(s, "-1") {
		return nil, nil
	}
	value, err := strconv.Atoi(strings.TrimSpace(s))
	if err == nil {
		return util.GetIntPointer(value), nil
	}
	parts := bytesPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(parts) < 3 {
		return nil, errors.Wrap(invalidByteQuantityError(), "Unable to find match by pattern")
	}

	floatValue, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, errors.Wrap(invalidByteQuantityError(), "Unable to convert to integer")
	}

	var bytes float64
	unit := strings.ToUpper(parts[2])
	switch unit {
	case "T", "TB":
		bytes = floatValue * TERABYTE
	case "G", "GB":
		bytes = floatValue * GIGABYTE
	case "M", "MB":
		bytes = floatValue * MEGABYTE
	}

	return util.GetIntPointer(int(bytes / float64(MEGABYTE))), nil
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

func invalidByteQuantityError() error {
	return errors.New("Byte quantity must be an integer with a unit of measurement like M, MB, G, or GB")
}

func invalidTimePatternError() error {
	return errors.New("Time to add must have format like D, M or H")
}
