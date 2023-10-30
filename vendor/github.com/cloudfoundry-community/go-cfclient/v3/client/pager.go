package client

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

var ErrNoResultsReturned = errors.New("expected 1 or more results, but got 0")
var ErrExactlyOneResultNotReturned = errors.New("expected exactly 1 result, but got less or more than 1")

type Pager struct {
	NextPageReader     *path.QuerystringReader
	PreviousPageReader *path.QuerystringReader

	TotalResults int
	TotalPages   int
}

func NewPager(pagination resource.Pagination) *Pager {
	nextPageReader, _ := path.NewQuerystringReader(pagination.Next.Href)
	previousPageReader, _ := path.NewQuerystringReader(pagination.Previous.Href)

	return &Pager{
		NextPageReader:     nextPageReader,
		PreviousPageReader: previousPageReader,
		TotalResults:       pagination.TotalResults,
		TotalPages:         pagination.TotalPages,
	}
}

func (p *Pager) HasNextPage() bool {
	return p.NextPageReader != nil
}

func (p *Pager) NextPage(opts ListOptioner) {
	if p.HasNextPage() {
		page := p.NextPageReader.Int(PageField)
		perPage := p.NextPageReader.Int(PerPageField)
		opts.CurrentPage(page, perPage)
	}
}

func (p *Pager) HasPreviousPage() bool {
	return p.PreviousPageReader != nil
}

func (p *Pager) PreviousPage(opts ListOptioner) {
	if p.HasPreviousPage() {
		page := p.PreviousPageReader.Int(PageField)
		perPage := p.PreviousPageReader.Int(PerPageField)
		opts.CurrentPage(page, perPage)
	}
}

type ListFunc[T ListOptioner, R any] func(opts T) ([]R, *Pager, error)

func AutoPage[T ListOptioner, R any](opts T, list ListFunc[T, R]) ([]R, error) {
	var all []R
	for {
		page, pager, err := list(opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, nil
}

// Single returns a single object from the call to list or an error if matches > 1 or matches < 1
func Single[T ListOptioner, R any](opts T, list ListFunc[T, R]) (R, error) {
	matches, _, err := list(opts)
	if err != nil {
		return *new(R), err
	}
	if len(matches) != 1 {
		return *new(R), ErrExactlyOneResultNotReturned
	}
	return matches[0], nil
}

// First returns the first object from the call to list or an error if matches < 1
func First[T ListOptioner, R any](opts T, list ListFunc[T, R]) (R, error) {
	matches, _, err := list(opts)
	if err != nil {
		return *new(R), err
	}
	if len(matches) < 1 {
		return *new(R), ErrNoResultsReturned
	}
	return matches[0], nil
}
