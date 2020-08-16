package web

// SearchResult abstracts the most important fields for generic Bleve search results.
// type searchResult struct {
// 	ID        string                       `json:"id"`
// 	Score     float64                      `json:"score"`
// 	Locations *search.FieldTermLocationMap `json:"locations,omitempty"`
// 	Fragments *search.FieldFragmentMap     `json:"fragments,omitempty"`
// }

// // toSearchResults converts Bleve search results into abstract search results.
// func toSearchResults(results *bleve.SearchResult) []searchResult {
// 	searchResults := make([]searchResult, len(results.Hits))

// 	for i := 0; i < len(searchResults); i++ {
// 		result := results.Hits[i]

// 		searchResult := searchResult{
// 			ID:    result.ID,
// 			Score: result.Score,
// 		}

// 		if len(result.Locations) > 0 {
// 			searchResult.Locations = &result.Locations
// 		}

// 		if len(result.Fragments) > 0 {
// 			searchResult.Fragments = &result.Fragments
// 		}

// 		searchResults[i] = searchResult
// 	}

// 	return searchResults
// }
