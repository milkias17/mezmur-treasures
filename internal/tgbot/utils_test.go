package tgbot

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetMessageTexts(t *testing.T) {
	text := "Lorem"

	res := getMessageTexts(text, nil)
	if len(res) != 1 {
		t.Errorf("Expected 1, got %v", res)
	}
	if res[0] != text {
		t.Errorf("Expected %v, got %v", text, res)
	}

	res = getMessageTexts(text, &text)
	if res[0] != strings.Repeat(text, 2) {
		t.Errorf("Expected %v, got %v", strings.Repeat(text, 2), res)
	}

	text = "aaaaa"
	res = getMessageTexts(strings.Repeat(text, 5000/len(text)), nil)
	if len(res) != 2 {
		t.Errorf("Expected 2, got %v", res)
	}
	if !reflect.DeepEqual(res, []string{strings.Repeat("a", 4096), strings.Repeat("a", 904)}) {
		t.Errorf("Expected %v, got %v", []string{strings.Repeat("a", 4096), strings.Repeat("a", 904)}, res)
	}
}
