package jsonpointergo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	JSONPointerEmptyPointer   = ""
	JSONPointerSeparatorToken = "/"
	JSONPointerEscapeToken    = "~"
	JSONPointerSlashEncoded   = "~1"
	JSONPointerTildaEncoded   = "~0"
)

// JSONObject is a type alias for a map with string keys and values of
// any type
type JSONObject = map[string]any

// JSONPointer struct holds the parsed reference tokens of a JSON
// Pointer
type JSONPointer struct {
	// Slice of reference tokens derived from the JSON Pointer
	referenceTokens []string
}

// NewJSONPointer creates a new JSONPointer instance from a JSON
// Pointer string
func NewJSONPointer(jsonPointer string) (*JSONPointer, error) {
	tokens, err := parseJSONPointerString(jsonPointer)
	if err != nil {
		return nil, err
	}
	return &JSONPointer{
		referenceTokens: tokens,
	}, nil
}

// parseJSONPointerString parses a JSON Pointer string into its
// reference tokens
func parseJSONPointerString(jsonPointer string) ([]string, error) {
	if jsonPointer == JSONPointerEmptyPointer {
		return []string{}, nil
	}
	if !strings.HasPrefix(jsonPointer, JSONPointerSeparatorToken) {
		return nil, fmt.Errorf(
			"jsonpointer: a jsonpointer should start with a reference to the root value: %v",
			JSONPointerSeparatorToken,
		)
	}
	// Split the JSON Pointer into tokens
	tokens := strings.Split(jsonPointer, JSONPointerSeparatorToken)
	return tokens[1:], nil
}

// GetValue retrieves the value from the JSON document based on the
// JSON Pointer
func (jp *JSONPointer) GetValue(document JSONObject) (any, error) {
	if document == nil {
		return nil, fmt.Errorf(
			"jsonpointer: the JSON document provided is nil",
		)
	}
	var subDocument any
	// Start with the root of the JSON document
	subDocument = document
	for i, tokenRefEncoded := range jp.referenceTokens {
		tokenRef := decodeJSONPointerReference(tokenRefEncoded)
		// Handle the case where the value is a JSON object
		jsonDoc, ok := subDocument.(JSONObject)
		if ok {
			value, ok := jsonDoc[tokenRef]
			if !ok {
				return nil, fmt.Errorf(
					"jsonpointer: the document provided does not have the following reference: %v, %v",
					tokenRef,
					i,
				)
			}
			subDocument = value
			continue
		}
		// Handle the case where the value is an Array
		jsonArray, ok := subDocument.([]any)
		if ok {
			index, err := strconv.Atoi(tokenRef)
			if err != nil {
				return nil, fmt.Errorf(
					"jsonpointer: the reference is trying to access a field on an array: %v",
					tokenRef,
				)
			}
			if index < 0 || index >= len(jsonArray) {
				return nil, fmt.Errorf(
					"jsonpointer: the index provided [%v] is trying to access an out of bond item on an array of length %v",
					index,
					len(jsonArray),
				)
			}
			subDocument = jsonArray[index]
			continue
		}
		// Handle the case where the value is a single value
		return nil, fmt.Errorf(
			"jsonpointer: the reference is trying to access a single value: %v. Type of subdocument: %T",
			tokenRef,
			subDocument,
		)
	}
	return subDocument, nil
}

// decodeJSONPointerReference decodes a reference token by replacing
// escape sequences
func decodeJSONPointerReference(ref string) string {
	// Replace "~1" with "/"
	ref = strings.ReplaceAll(
		ref,
		JSONPointerSlashEncoded,
		JSONPointerSeparatorToken,
	)
	// Replace "~0" with "~"
	return strings.ReplaceAll(
		ref,
		JSONPointerTildaEncoded,
		JSONPointerEscapeToken,
	)
}
