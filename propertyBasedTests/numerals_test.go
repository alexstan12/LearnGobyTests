package propertyBasedTests

import (
	"fmt"
	"log"
	"testing"
	"testing/quick"
)

var cases = []struct {
	Description string
	Arabic      uint16
	Roman       string
}{
	{"1 gets converted to I", 1, "I"},
	{"2 gets converted to II", 2, "II"},
	{"3 gets converted to III", 3, "III"},
	{"4 gets converted to IV (can't repeat more than 3 times", 4, "IV"},
	{"5 gets converted to V", 5, "V"},
	{"9 gets converted to IX", 9, "IX"},
	{"10 gets converted to X", 10, "X"},
	{"10 gets converted to X", 10, "X"},
	{"14 gets converted to XIV", 14, "XIV"},
	{"18 gets converted to XVIII", 18, "XVIII"},
	{"20 gets converted to XX", 20, "XX"},
	{"39 gets converted to XXXIX", 39, "XXXIX"},
	{"40 gets converted to XL", 40, "XL"},
	{"47 gets converted to XLVII", 47, "XLVII"},
	{"49 gets converted to XLIX", 49, "XLIX"},
	{"50 gets converted to L", 50, "L"},
}

func TestRomanNumerals(t *testing.T) {

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)
			want := test.Roman
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestConvertingToArabic(t *testing.T) {
	for _, test := range cases {
		t.Run(fmt.Sprintf("%q gets converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
			got := ConvertToArabic(test.Roman)
			want := test.Arabic
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	}
}

func TestPropertiesOfConversion(t *testing.T) {
	assertion := func(arabic uint16) bool {
		if arabic < 0 || arabic > 3999 {
			log.Println(arabic)
			return true
		}
		roman := ConvertToRoman(arabic)
		fromRoman := ConvertToArabic(roman)
		return fromRoman == arabic
	}

	if err := quick.Check(assertion, &quick.Config{
		MaxCount: 1000,
	}); err != nil {
		t.Error("failed checks", err)
	}
}
