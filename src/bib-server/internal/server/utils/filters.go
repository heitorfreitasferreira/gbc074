package utils

import (
	"errors"
	"library-manager/bib-server/internal/database"
	"regexp"
	"strings"
)

// FilterBooks filters books based on the query string.
func FilterBooks(books []database.Book, query string) ([]database.Book, error) {
	// Define the regular expression for the query.
	queryRegex := `^(titulo:[^&|]*|autor:[^&|]*|isbn:[^&|]*)([&|\|](titulo:[^&|]*|autor:[^&|]*|isbn:[^&|]*))?$`
	matched, err := regexp.MatchString(queryRegex, query)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("invalid query format")
	}

	// Split the query by the operators.
	operator := "\000"
	if strings.Contains(query, "&") {
		operator = "&"
	} else if strings.Contains(query, "|") {
		operator = "|"
	}

	// Split the query into filters.
	filters := strings.Split(query, operator)
	if len(filters) > 2 {
		return nil, errors.New("query contains more than two filters")
	}

	// Check for duplicate filters.
	filterMap := make(map[string]string)
	for _, filter := range filters {
		parts := strings.Split(filter, ":")
		key := parts[0]
		value := parts[1]
		if _, exists := filterMap[key]; exists {
			return nil, errors.New("duplicate filter detected in the query")
		}
		filterMap[key] = value
	}

	// Filter the books based on the filters and operator.
	var result []database.Book
	for _, book := range books {
		match := evaluateBook(book, filterMap, operator)
		if match {
			result = append(result, book)
		}
	}

	return result, nil
}

// evaluateBook checks if a book matches the given filters based on the operator.
func evaluateBook(book database.Book, filters map[string]string, operator string) bool {
	matchCount := 0

	for key, value := range filters {
		switch key {
		case "titulo":
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(value)) {
				matchCount++
			}
		case "autor":
			if strings.Contains(strings.ToLower(book.Author), strings.ToLower(value)) {
				matchCount++
			}
		case "isbn":
			if strings.Contains(strings.ToLower(book.ISBN), strings.ToLower(value)) {
				matchCount++
			}
		}
	}

	if operator == "&" {
		return matchCount == len(filters)
	} else if operator == "|" {
		return matchCount > 0
	} else {
		return matchCount == 1
	}
}
