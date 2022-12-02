package client

import (
	"time"
)

type FilterModifier int

const (
	FilterModifierNone FilterModifier = iota
	FilterModifierGreaterThan
	FilterModifierLessThan
	FilterModifierGreaterThanOrEqual
	FilterModifierLessThanOrEqual
)

func (r FilterModifier) String() string {
	switch r {
	case FilterModifierGreaterThan:
		return "gt"
	case FilterModifierGreaterThanOrEqual:
		return "gte"
	case FilterModifierLessThan:
		return "lt"
	case FilterModifierLessThanOrEqual:
		return "lte"
	}
	return ""
}

type TimestampFilter struct {
	Timestamp []time.Time
	Operator  FilterModifier
}

func (t *TimestampFilter) EqualTo(createdAt ...time.Time) {
	t.Timestamp = createdAt
}

func (t *TimestampFilter) Before(createdAt time.Time) {
	t.Timestamp = []time.Time{
		createdAt,
	}
	t.Operator = FilterModifierLessThan
}

func (t *TimestampFilter) BeforeOrEqualTo(createdAt time.Time) {
	t.Timestamp = []time.Time{
		createdAt,
	}
	t.Operator = FilterModifierLessThanOrEqual
}

func (t *TimestampFilter) After(createdAt time.Time) {
	t.Timestamp = []time.Time{
		createdAt,
	}
	t.Operator = FilterModifierGreaterThan
}

func (t *TimestampFilter) AfterOrEqualTo(createdAt time.Time) {
	t.Timestamp = []time.Time{
		createdAt,
	}
	t.Operator = FilterModifierGreaterThanOrEqual
}

type Filter struct {
	Values []string
	Not    bool
}

func (f *Filter) EqualTo(v ...string) {
	f.Values = v
}

func (f *Filter) NotEqualTo(v ...string) {
	f.Values = v
	f.Not = true
}
