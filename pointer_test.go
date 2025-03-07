package jsonpointergo

import (
	"reflect"
	"testing"
)

func TestNewJSONPointer(t *testing.T) {
	tests := []struct {
		jsonPointer string
		wantErr     bool
	}{
		{"/foo/bar", false},
		{"", false},
		{"foo/bar", true},
	}

	for _, test := range tests {
		_, err := NewJSONPointer(test.jsonPointer)
		if (err != nil) != test.wantErr {
			t.Errorf(
				"NewJSONPointer(%v) error = %v, wantErr %v",
				test.jsonPointer,
				err,
				test.wantErr,
			)
		}
	}
}

func TestGetValue(t *testing.T) {
	document := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
		"array": []any{1, 2, 3},
	}

	tests := []struct {
		jsonPointer string
		want        any
		wantErr     bool
	}{
		{"/foo/bar", "baz", false},
		{"/array/0", 1, false},
		{"/array/3", nil, true},
		{"/nonexistent", nil, true},
		{"", document, false},
	}

	for _, test := range tests {
		jp, err := NewJSONPointer(test.jsonPointer)
		if err != nil {
			t.Fatalf(
				"NewJSONPointer(%v) error = %v",
				test.jsonPointer,
				err,
			)
		}
		got, err := jp.GetValue(document)
		if (err != nil) != test.wantErr {
			t.Errorf(
				"JSONPointer.GetValue() error = %v, wantErr %v",
				err,
				test.wantErr,
			)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(
				"JSONPointer.GetValue() = %v, want %v",
				got,
				test.want,
			)
		}
	}
}

func TestDecodeJSONPointerReference(t *testing.T) {
	tests := []struct {
		ref  string
		want string
	}{
		{"~1", "/"},
		{"~0", "~"},
		{"~01", "~1"},
	}

	for _, test := range tests {
		if got := decodeJSONPointerReference(test.ref); got != test.want {
			t.Errorf(
				"decodeJSONPointerReference(%v) = %v, want %v",
				test.ref,
				got,
				test.want,
			)
		}
	}
}
