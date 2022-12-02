package client

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const filterTagName = "qs"

const (
	DefaultPage     = 1
	DefaultPageSize = 50
)

const (
	PageField    = "page"
	PerPageField = "per_page"
)

type ListOptioner interface {
	CurrentPage(page, perPage int)
	ToQueryString() url.Values
}

// ListOptions is the shared common type for all other list option types
type ListOptions struct {
	Page    int    `qs:"page"`
	PerPage int    `qs:"per_page"`
	OrderBy string `qs:"order_by"`

	LabelSelector Filter          `qs:"label_selector"`
	CreateAts     TimestampFilter `qs:"created_ats"`
	UpdatedAts    TimestampFilter `qs:"updated_ats"`
}

// NewListOptions creates a default list options with page and page size set
func NewListOptions() *ListOptions {
	return &ListOptions{
		Page:    DefaultPage,
		PerPage: DefaultPageSize,
	}
}

func (lo *ListOptions) CurrentPage(page, perPage int) {
	lo.PerPage = perPage
	lo.Page = page
}

func (lo *ListOptions) ToQueryString(subOptionsPtr any) url.Values {
	s := ListOptionsSerializer{}
	s.Add(&lo)
	s.Add(subOptionsPtr)
	val, err := s.Serialize()
	if err != nil {
		panic("go-cfclient filter bug: " + err.Error())
	}
	return val
}

func appendQueryStrings(a, b url.Values) url.Values {
	for k, v := range a {
		b.Set(k, v[0]) // url.Values get only returns 1st item anyway
	}
	return b
}

var filterType = reflect.TypeOf(Filter{})
var timeFilterType = reflect.TypeOf(TimestampFilter{})
var timeType = reflect.TypeOf(time.Time{})

type ListOptionsSerializer struct {
	optStructs []any
}

func (los *ListOptionsSerializer) Add(optStruct any) {
	los.optStructs = append(los.optStructs, optStruct)
}

func (los *ListOptionsSerializer) Serialize() (url.Values, error) {
	var values url.Values
	for _, opt := range los.optStructs {
		val, err := los.serializeOptionStruct(opt)
		if err != nil {
			return url.Values{}, err
		}
		values = appendQueryStrings(values, val)
	}
	return values, nil
}

func (los *ListOptionsSerializer) serializeOptionStruct(o any) (url.Values, error) {
	if o == nil {
		return url.Values{}, nil
	}

	val := reflect.ValueOf(o)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return url.Values{}, nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return url.Values{}, nil
	}

	return los.reflectValues(val)
}

func (los *ListOptionsSerializer) reflectValues(val reflect.Value) (url.Values, error) {
	values := url.Values{}
	for i := 0; i < val.Type().NumField(); i++ {
		sv := val.Field(i)
		rawTag := val.Type().Field(i).Tag.Get(filterTagName)
		if rawTag == "" || rawTag == "-" {
			continue
		}
		if los.isEmptyValue(sv) {
			continue
		}

		sv = los.getNonPointerValue(sv)
		switch sv.Type() {
		case filterType:
			err := los.reflectFilter(sv, rawTag, values)
			if err != nil {
				return url.Values{}, err
			}
		case timeFilterType:
			err := los.reflectTimestampFilter(sv, rawTag, values)
			if err != nil {
				return url.Values{}, err
			}
		default:
			values.Add(rawTag, fmt.Sprint(sv.Interface()))
		}
	}

	return values, nil
}

func (los *ListOptionsSerializer) reflectFilter(val reflect.Value, tag string, values url.Values) error {
	var filterStrings []string
	var not bool

	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Field(i)
		if f.Kind() == reflect.Slice || f.Kind() == reflect.Array {
			for ti := 0; ti < f.Len(); ti++ {
				tv := f.Index(ti)
				if tv.Kind() == reflect.String {
					s := tv.Interface().(string)
					filterStrings = append(filterStrings, s)
				}
			}
		} else if f.Kind() == reflect.Bool {
			not = f.Interface().(bool)
		}
	}

	if len(filterStrings) > 0 {
		key := tag
		if not {
			key = key + "[not]"
		}
		values.Add(key, strings.Join(filterStrings, ","))
	}

	return nil
}

func (los *ListOptionsSerializer) reflectTimestampFilter(val reflect.Value, tag string, values url.Values) error {
	var timestamps []string
	var relationalOperator FilterModifier

	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Field(i)
		if f.Kind() == reflect.Slice || f.Kind() == reflect.Array {
			for ti := 0; ti < f.Len(); ti++ {
				tv := f.Index(ti)
				if tv.Type() == timeType {
					timestamp := tv.Interface().(time.Time)
					timestamps = append(timestamps, timestamp.Format(time.RFC3339))
				}
			}
		} else if f.Kind() == reflect.Int {
			relationalOperator = f.Interface().(FilterModifier)
		}
	}

	if len(timestamps) > 0 {
		key := tag
		if relationalOperator != FilterModifierNone {
			key = key + "[" + relationalOperator.String() + "]"
		}
		values.Add(key, strings.Join(timestamps, ","))
	}

	return nil
}

func (los *ListOptionsSerializer) isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (los *ListOptionsSerializer) getNonPointerValue(sv reflect.Value) reflect.Value {
	for sv.Kind() == reflect.Ptr {
		if sv.IsNil() {
			break
		}
		sv = sv.Elem()
	}
	return sv
}
