package isosegment

// classification is a grouping of isolation segments that should meet some property.
// - missing: segments that need to be updated to meet the property
// - extra: segments that meet the property but don't need to (can be cleaned up)
type classification struct {
	missing []Segment
	extra   []Segment
}

// classify accepts a list of segments that should meet some property and a list
// of segments that currently do meet that property, and returns a classification.
func classify(desired []Segment, current []Segment) classification {
	currentSegments := make(map[string]*Segment)
	for i := range current {
		currentSegments[current[i].Name] = &current[i]
	}
	desiredSegments := make(map[string]*Segment)
	for i := range desired {
		desiredSegments[desired[i].Name] = &desired[i]
	}

	var missing []Segment
	for name, seg := range desiredSegments {
		if _, ok := currentSegments[name]; !ok {
			missing = append(missing, *seg)
		}
	}

	var extra []Segment
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
func (c classification) update(arg string, apply, cleanup func(s *Segment, arg string) error) error {
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
