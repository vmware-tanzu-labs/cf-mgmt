package resource

import (
	"fmt"
	"strings"
)

// Metadata allows you to tag API resources with information that does not directly affect its functionality.
type Metadata struct {
	Labels      map[string]*string `json:"labels"`
	Annotations map[string]*string `json:"annotations"`
}

// NewMetadata creates a new metadata instance
func NewMetadata() *Metadata {
	return &Metadata{}
}

// WithAnnotation is a fluent method alias for SetAnnotation
func (m *Metadata) WithAnnotation(prefix, key string, v string) *Metadata {
	m.SetAnnotation(prefix, key, v)
	return m
}

// WithLabel is a fluent method alias for SetLabel
func (m *Metadata) WithLabel(prefix, key string, v string) *Metadata {
	m.SetLabel(prefix, key, v)
	return m
}

// SetAnnotation to the metadata instance
//
// The prefix and value are optional and may be an empty string. The key must be at least 1 character in length.
func (m *Metadata) SetAnnotation(prefix, key string, v string) {
	if m.Annotations == nil {
		m.Annotations = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Annotations[fmt.Sprintf("%s/%s", prefix, key)] = &v
	} else {
		m.Annotations[key] = &v
	}
}

// RemoveAnnotation removes an annotation by setting the specified key's value to nil which can then be passed to the API
func (m *Metadata) RemoveAnnotation(prefix, key string) {
	if m.Annotations == nil {
		m.Annotations = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Annotations[fmt.Sprintf("%s/%s", prefix, key)] = nil
	} else {
		m.Annotations[key] = nil
	}
}

// SetLabel to the metadata instance
//
// The prefix and value are optional and may be an empty string. The key must be at least 1 character in length.
func (m *Metadata) SetLabel(prefix, key string, v string) {
	if m.Labels == nil {
		m.Labels = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Labels[fmt.Sprintf("%s/%s", prefix, key)] = &v
	} else {
		m.Labels[key] = &v
	}
}

// RemoveLabel removes a label by setting the specified key's value to nil which can then be passed to the API
func (m *Metadata) RemoveLabel(prefix, key string) {
	if m.Labels == nil {
		m.Labels = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Labels[fmt.Sprintf("%s/%s", prefix, key)] = nil
	} else {
		m.Labels[key] = nil
	}
}

// Clear automatically calls Remove on all annotations and labels present in the metadata instance
func (m *Metadata) Clear() {
	splitKey := func(k string) (string, string) {
		p := strings.Split(k, "/")
		if len(p) == 1 {
			return "", p[0]
		}
		return p[0], p[1]
	}
	for k := range m.Annotations {
		prefix, key := splitKey(k)
		m.RemoveAnnotation(prefix, key)
	}
	for k := range m.Labels {
		prefix, key := splitKey(k)
		m.RemoveLabel(prefix, key)
	}
}
