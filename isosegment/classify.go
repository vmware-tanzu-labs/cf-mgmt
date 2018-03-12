package isosegment

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// classification is a grouping of isolation segments that should meet some property.
// - missing: segments that need to be updated to meet the property
// - extra: segments that meet the property but don't need to (can be cleaned up)
type classification struct {
	missing []cfclient.IsolationSegment
	extra   []cfclient.IsolationSegment
}

// classify accepts a list of segments that should meet some property and a list
// of segments that currently do meet that property, and returns a classification.
func classify(desired []cfclient.IsolationSegment, current []cfclient.IsolationSegment) classification {
	currentSegments := make(map[string]*cfclient.IsolationSegment)
	for i := range current {
		currentSegments[current[i].Name] = &current[i]
	}
	desiredSegments := make(map[string]*cfclient.IsolationSegment)
	for i := range desired {
		desiredSegments[desired[i].Name] = &desired[i]
	}

	var missing []cfclient.IsolationSegment
	for name, seg := range desiredSegments {
		if _, ok := currentSegments[name]; !ok {
			missing = append(missing, *seg)
		}
	}

	var extra []cfclient.IsolationSegment
	for name, seg := range currentSegments {
		if _, ok := desiredSegments[name]; !ok {
			extra = append(extra, *seg)
		}
	}

	return classification{missing, extra}
}

// update accepts two actions:
// 1. apply - run on every 'missing' segment in the classification
// 2. cleanup - run on every 'extra' segment in the classification
func (c classification) update(arg string, apply, cleanup func(s *cfclient.IsolationSegment, arg string) error) error {
	for i := range c.missing {
		err := apply(&c.missing[i], arg)
		if err != nil {
			return err
		}
	}
	for i := range c.extra {
		err := cleanup(&c.extra[i], arg)
		if err != nil {
			return err
		}
	}
	return nil
}
