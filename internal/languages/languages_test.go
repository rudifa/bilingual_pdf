package languages

import "testing"

func TestValidate_ValidCodes(t *testing.T) {
	valid := []string{"en", "fr", "es", "de", "zh", "ja", "ko", "ar"}
	for _, code := range valid {
		if err := Validate(code); err != nil {
			t.Errorf("Validate(%q) should be valid, got error: %v", code, err)
		}
	}
}

func TestValidate_InvalidCodes(t *testing.T) {
	invalid := []string{"xx", "zz", "english", "", "123"}
	for _, code := range invalid {
		if err := Validate(code); err == nil {
			t.Errorf("Validate(%q) should return error", code)
		}
	}
}

func TestSupported_NonEmpty(t *testing.T) {
	langs := Supported()
	if len(langs) == 0 {
		t.Error("Supported() should return non-empty list")
	}
}

func TestSupported_Sorted(t *testing.T) {
	langs := Supported()
	for i := 1; i < len(langs); i++ {
		if langs[i].Code < langs[i-1].Code {
			t.Errorf("Supported() not sorted: %q comes after %q", langs[i].Code, langs[i-1].Code)
		}
	}
}

func TestName_Known(t *testing.T) {
	if got := Name("fr"); got != "French" {
		t.Errorf("Name(\"fr\") = %q, want \"French\"", got)
	}
}

func TestName_Unknown(t *testing.T) {
	if got := Name("xx"); got != "xx" {
		t.Errorf("Name(\"xx\") = %q, want \"xx\"", got)
	}
}
