package web

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

// SearchResult abstracts the most important fields for generic Bleve search results.
type SearchResult struct {
	ID        string                       `json:"id"`
	Score     float64                      `json:"score"`
	Locations *search.FieldTermLocationMap `json:"locations,omitempty"`
	Fragments *search.FieldFragmentMap     `json:"fragments,omitempty"`
}

// ToSearchResults converts Bleve search results into abstract search results.
func ToSearchResults(results *bleve.SearchResult) []SearchResult {
	searchResults := make([]SearchResult, len(results.Hits))

	for i := 0; i < len(searchResults); i++ {
		result := results.Hits[i]

		searchResult := SearchResult{
			ID:    result.ID,
			Score: result.Score,
		}

		if len(result.Locations) > 0 {
			searchResult.Locations = &result.Locations
		}

		if len(result.Fragments) > 0 {
			searchResult.Fragments = &result.Fragments
		}

		searchResults[i] = searchResult
	}

	return searchResults
}
