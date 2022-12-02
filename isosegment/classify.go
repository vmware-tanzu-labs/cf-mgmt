package isosegment

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// classification is a grouping of isolation segments that should meet some property.
// - missing: segments that need to be updated to meet the property
// - extra: segments that meet the property but don't need to (can be cleaned up)
type classification struct {
	missing []*resource.IsolationSegment
	extra   []*resource.IsolationSegment
}

// classify accepts a list of segments that should meet some property and a list
// of segments that currently do meet that property, and returns a classification.
func classify(desired []*resource.IsolationSegment, current []*resource.IsolationSegment) classification {
	currentSegments := make(map[string]**resource.IsolationSegment)
	for i := range current {
		currentSegments[current[i].Name] = &current[i]
	}
	desiredSegments := make(map[string]**resource.IsolationSegment)
	for i := range desired {
		desiredSegments[desired[i].Name] = &desired[i]
	}

	var missing []*resource.IsolationSegment
	for name, seg := range desiredSegments {
		if _, ok := currentSegments[name]; !ok {
			missing = append(missing, *seg)
		}
	}

	var extra []*resource.IsolationSegment
	for name, seg := range currentSegments {
		if _, ok := desiredSegments[name]; !ok {
			extra = append(extra, *seg)
		}
	}

	return classification{missing, extra}
}
