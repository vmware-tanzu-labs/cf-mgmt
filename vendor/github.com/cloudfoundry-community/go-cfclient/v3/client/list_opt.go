package client

import (
	"fmt"
	"net/url"
	"reflect"
)

const (
	filterTagName = "qs"

	DefaultPage     = 1
	DefaultPageSize = 50

	PageField    = "page"
	PerPageField = "per_page"
)

type ListOptionsSerializer interface {
	Serialize(values url.Values, tag string) error
}

var listOptionsSerializerType = reflect.TypeOf((*ListOptionsSerializer)(nil)).Elem()

type ListOptioner interface {
	CurrentPage(page, perPage int)
	ToQueryString() (url.Values, error)
}

// ListOptions is the shared common type for all other list option types
type ListOptions struct {
	Page       int             `qs:"page"`
	PerPage    int             `qs:"per_page"`
	OrderBy    string          `qs:"order_by"`
	LabelSel   LabelSelector   `qs:"label_selector"`
	CreateAts  TimestampFilter `qs:"created_ats"`
	UpdatedAts TimestampFilter `qs:"updated_ats"`
}

// NewListOptions creates a default list options with page and page size set
func NewListOptions() *ListOptions {
	return &ListOptions{
		Page:    DefaultPage,
		PerPage: DefaultPageSize,
	}
}

func (lo *ListOptions) CurrentPage(page, perPage int) {
	lo.Page = page
	lo.PerPage = perPage
}

func (lo ListOptions) Serialize(values url.Values, _ string) error {
	return serializeField(values, reflect.ValueOf(lo))
}

func (lo *ListOptions) ToQueryString(subOptionsPtr any) (url.Values, error) {
	if subOptionsPtr != nil {
		values := url.Values{}
		err := serializeField(values, reflect.ValueOf(subOptionsPtr))
		return values, err
	}
	return nil, nil
}

func serializeField(values url.Values, val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	valTypes := val.Type()
	for i := 0; i < valTypes.NumField(); i++ {
		fieldType := valTypes.Field(i)
		rawTag := fieldType.Tag.Get(filterTagName)
		if (rawTag != "" && rawTag != "-") || fieldType.Type.Implements(listOptionsSerializerType) {
			sv := val.Field(i)
			if sv.Kind() == reflect.Ptr {
				if sv.IsNil() {
					continue
				}
				sv = sv.Elem()
			}
			if sv.IsZero() {
				continue
			}
			svi := sv.Interface()
			if filter, ok := svi.(ListOptionsSerializer); ok {
				if err := filter.Serialize(values, rawTag); err != nil {
					return err
				}
			} else {
				values.Add(rawTag, fmt.Sprintf("%v", svi))
			}
		}
	}
	return nil
}
