package db

import "testing"

func TestIsEnglish(t *testing.T) {
	text := ""

	if !isEnglish(text) {
		t.Errorf("Expected true, got false for empty string")
	}

	text = "I am an english text"

	if !isEnglish(text) {
		t.Errorf("Expected true, got false for %s", text)
	}

	text = "ዘማሪ ተስፋዬ ጋቢሶ ርዕስ ሥራህ ያመሰግንሃል"

	if isEnglish(text) {
		t.Errorf("Expected false, got true for %s", text)
	}

}
