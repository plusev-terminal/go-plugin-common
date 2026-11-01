package account

import "testing"

func TestOptionalPrice(t *testing.T) {
	for _, zero := range []string{"", " ", "0", "0.0", "0.00000000", "-0", ".0"} {
		if got := OptionalPrice(zero); got != "" {
			t.Errorf("OptionalPrice(%q) = %q, want empty", zero, got)
		}
	}
	for _, price := range []string{"9300", "0.5", "100.000", "0.00000001"} {
		if got := OptionalPrice(price); got != price {
			t.Errorf("OptionalPrice(%q) = %q, want unchanged", price, got)
		}
	}
}

func TestTriggerSourceValid(t *testing.T) {
	for _, s := range []TriggerSource{TriggerSourceLast, TriggerSourceMark, TriggerSourceIndex} {
		if !s.Valid() {
			t.Errorf("%s must be valid", s)
		}
	}
	for _, s := range []TriggerSource{"", "MARK_PRICE", "contract"} {
		if s.Valid() {
			t.Errorf("%q must not be valid", s)
		}
	}
}
