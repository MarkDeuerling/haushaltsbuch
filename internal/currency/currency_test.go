package currency

import (
	"testing"
)

func TestAddAmount(t *testing.T) {
	eur1, _ := NewCurrency("100", "EUR")
	eur2, _ := NewCurrency("10", "EUR")
	sum := eur1.Add(eur2)

	want, _ := NewCurrency("110", "EUR")

	if sum.Amount() != want.Amount() {
		t.Errorf("Amount is incorrect, got: %s, want: %s", sum.Amount(), want.Amount())
	}
}

func TestSubAmount(t *testing.T) {
	eur1, _ := NewCurrency("100", "EUR")
	eur2, _ := NewCurrency("10", "EUR")
	sum := eur1.Sub(eur2)

	want, _ := NewCurrency("90", "EUR")

	if sum.Amount() != want.Amount() {
		t.Errorf("Amount is incorrect, got: %s, want: %s", sum.Amount(), want.Amount())
	}
}

func TestCurrencyCode(t *testing.T) {
	want, _ := NewCurrency("9", "EUR")
	got, _ := NewCurrency("100", "EUR")

	if got.Code() != want.Code() {
		t.Errorf("code is incorrect, got: %s, want: %s", got.Code(), want.Code())
	}
}

func TestCompareCurrencyCode(t *testing.T) {
	want, _ := NewCurrency("9", "EUR")
	got, _ := NewCurrency("100", "EUR")

	if want.Code() != got.Code() {
		t.Errorf("currency code is not the same, got: %s, want: %s", got.Code(), want.Code())
	}
}

func TestWrongAmount(t *testing.T) {
	_, err := NewCurrency("10b0", "EUR")

	if err == nil {
		t.Errorf("want error for invalid input")
	}
}

func TestOverflowAdd(t *testing.T) {
	eur1, _ := NewCurrency("9223372036854775807", "EUR")
	eur2, _ := NewCurrency("1", "EUR")
	sum := eur1.Add(eur2)

	want, _ := NewCurrency("9223372036854775808", "EUR")

	if sum.Amount() != want.Amount() {
		t.Errorf("Amount is incorrect, got: %s, want: %s", sum.Amount(), want.Amount())
	}
}

func TestOverflowSub(t *testing.T) {
	eur1, _ := NewCurrency("-9223372036854775807", "EUR")
	eur2, _ := NewCurrency("1", "EUR")
	sum := eur1.Sub(eur2)

	want, _ := NewCurrency("-9223372036854775808", "EUR")

	if sum.Amount() != want.Amount() {
		t.Errorf("Amount is incorrect, got: %s, want: %s", sum.Amount(), want.Amount())
	}
}
