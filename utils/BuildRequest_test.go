package utils

import (
	"testing"
)

func TestBuildRequest_NoQueryParams(t *testing.T) {
	url := "https://example.com/path"
	result := BuildRequest(url, map[string]string{})
	if result != url {
		t.Errorf("Expected %s, got %s", url, result)
	}
}

func TestBuildRequest_SingleParam(t *testing.T) {
	url := "https://example.com/path"
	params := map[string]string{"foo": "bar"}
	result := BuildRequest(url, params)
	expected := "https://example.com/path?foo=bar"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestBuildRequest_MultipleParams(t *testing.T) {
	url := "https://example.com/path"
	params := map[string]string{"foo": "bar", "baz": "qux"}
	result := BuildRequest(url, params)
	if !(result == "https://example.com/path?baz=qux&foo=bar" || result == "https://example.com/path?foo=bar&baz=qux") {
		t.Errorf("Unexpected result: %s", result)
	}
}

func TestBuildRequest_SpecialChars(t *testing.T) {
	url := "https://example.com/path"
	params := map[string]string{"q": "hello world", "x": "a&b"}
	result := BuildRequest(url, params)
	if !(result == "https://example.com/path?q=hello+world&x=a%26b" || result == "https://example.com/path?x=a%26b&q=hello+world") {
		t.Errorf("Unexpected result: %s", result)
	}
}

func TestBuildRequest_OverrideExistingQuery(t *testing.T) {
	url := "https://example.com/path?foo=old"
	params := map[string]string{"foo": "new", "bar": "baz"}
	result := BuildRequest(url, params)
	if !(result == "https://example.com/path?bar=baz&foo=new" || result == "https://example.com/path?foo=new&bar=baz") {
		t.Errorf("Unexpected result: %s", result)
	}
}
