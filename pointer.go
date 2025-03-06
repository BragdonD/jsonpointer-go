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

type JSONObject = map[string]any

type JSONPointer struct {
	referenceTokens []string
}

func NewJSONPointer(jsonPointer string) (*JSONPointer, error) {
	tokens, err := parseJSONPointerString(jsonPointer)
	if err != nil {
		return nil, err
	}
	return &JSONPointer{
		referenceTokens: tokens,
	}, nil
}

func parseJSONPointerString(jsonPointer string) ([]string, error) {
	if jsonPointer == JSONPointerEmptyPointer {
		return nil, fmt.Errorf(
			"jsonpointer: the jsonpointer is empty",
		)
	}
	if !strings.HasPrefix(jsonPointer, JSONPointerSeparatorToken) {
		return nil, fmt.Errorf(
			"jsonpointer: a jsonpointer should start with a reference to the root value: %v",
			JSONPointerSeparatorToken,
		)
	}
	tokens := strings.Split(jsonPointer, JSONPointerSeparatorToken)
	return tokens[1:], nil
}

func (jp *JSONPointer) GetValue(document JSONObject) (any, error) {
	if document == nil {
		return nil, fmt.Errorf(
			"jsonpointer: the JSON document provided is nil",
		)
	}
	var subDocument any
	subDocument = document
	for i, tokenRefEncoded := range jp.referenceTokens {
		tokenRef := decodeJSONPointerReference(tokenRefEncoded)
		jsonDoc, ok := subDocument.(JSONObject)
		if ok {
			value, ok := jsonDoc[tokenRef]
			if !ok {
				return nil, fmt.Errorf(
					"jsonpointer: the document provided does not have the following reference: %v, %v",
					tokenRef, i,
				)
			}
			subDocument = value
			continue
		}
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
		return nil, fmt.Errorf("jsonpointer: the reference is trying to access a single value: %v. Type of subdocument: %T", tokenRef, subDocument)
	}
	return subDocument, nil
}

func decodeJSONPointerReference(ref string) string {
	refWithSlash := strings.ReplaceAll(
		ref,
		JSONPointerSlashEncoded,
		JSONPointerSeparatorToken,
	)
	return strings.ReplaceAll(
		refWithSlash,
		JSONPointerTildaEncoded,
		JSONPointerEscapeToken,
	)
}
