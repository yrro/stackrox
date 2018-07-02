package blevesearch

import (
	"fmt"
	"strings"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/logging"
	searchPkg "bitbucket.org/stack-rox/apollo/pkg/search"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

const maxSearchResponses = 2000

var logger = logging.LoggerForModule()

// NewPrefixQuery generates new query that matches prefixes
func NewPrefixQuery(field, prefix string) query.Query {
	prefixQuery := bleve.NewPrefixQuery(strings.ToLower(prefix))
	prefixQuery.SetField(field)
	return prefixQuery
}

// GetScopesQuery generates a disjunct query based on the scope values
func GetScopesQuery(scopes []*v1.Scope, scopeToQuery func(scope *v1.Scope) query.Query) query.Query {
	if len(scopes) != 0 {
		disjunctionQuery := bleve.NewDisjunctionQuery()
		for _, scope := range scopes {
			// Check if nil as some resources may not be applicable to scopes
			disjunctionQuery.AddQuery(scopeToQuery(scope))
		}
		return disjunctionQuery
	}
	return bleve.NewMatchAllQuery()
}

// FieldsToQuery converts a request and the options for the data type to Bleve types
func FieldsToQuery(request *v1.ParsedSearchRequest, optionsMap map[string]*v1.SearchField) (*query.ConjunctionQuery, error) {
	conjunctionQuery := bleve.NewConjunctionQuery()
	for fieldName, field := range request.GetFields() {
		searchField, ok := optionsMap[fieldName]
		if !ok {
			continue
		}
		queryFunc, ok := datatypeToQueryFunc[searchField.GetType()]
		if !ok {
			return nil, fmt.Errorf("Query for type %s is not implemented", searchField.GetType())
		}
		conjunct, err := queryFunc(searchField.GetFieldPath(), field.GetValues())
		if err != nil {
			return nil, err
		}
		conjunctionQuery.AddQuery(conjunct)
	}
	return conjunctionQuery, nil
}

// RunSearchRequest does something
func RunSearchRequest(objType string, request *v1.ParsedSearchRequest, index bleve.Index, scopeToQuery func(scope *v1.Scope) query.Query, optionsMap map[string]*v1.SearchField) ([]searchPkg.Result, error) {
	conjunctionQuery := bleve.NewConjunctionQuery(typeQuery(objType))
	queries, err := buildQuery(request, scopeToQuery, optionsMap)
	if err != nil {
		return nil, err
	}
	conjunctionQuery.AddQuery(queries...)
	return RunQuery(conjunctionQuery, index)
}

// RunQuery runs the actual query and then collapses the results into a simpler format
func RunQuery(query query.Query, index bleve.Index) ([]searchPkg.Result, error) {
	searchRequest := bleve.NewSearchRequest(query)
	// Initial size is 10 which seems small
	searchRequest.Size = maxSearchResponses
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	return collapseResults(searchResult), nil
}

func buildQuery(request *v1.ParsedSearchRequest, scopeToQuery func(scope *v1.Scope) query.Query, optionsMap map[string]*v1.SearchField) ([]query.Query, error) {
	var queries []query.Query
	queries = append(queries, GetScopesQuery(request.GetScopes(), scopeToQuery))
	if request.GetFields() != nil && len(request.GetFields()) != 0 {
		q, err := FieldsToQuery(request, optionsMap)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}
	if request.GetStringQuery() != "" {
		queries = append(queries, NewPrefixQuery("", request.GetStringQuery()))
	}
	return queries, nil
}

func typeQuery(objType string) query.Query {
	q := bleve.NewMatchQuery(objType)
	q.SetField("type")
	return q
}

func collapseResults(searchResult *bleve.SearchResult) (results []searchPkg.Result) {
	results = make([]searchPkg.Result, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		results = append(results, searchPkg.Result{
			ID:      hit.ID,
			Matches: hit.Fragments,
			Score:   hit.Score,
		})
	}
	return
}
