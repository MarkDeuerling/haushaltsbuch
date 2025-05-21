package currency

import (
	"errors"
	"math/big"
)

// Currency repräsentiert die Währung im System.
type Currency struct {
	amount *big.Int
	code   string
}

// NewCurrency erstellt eine neue Währung mit dem Betrag und dem Währungs-Code.
func NewCurrency(amount, code string) (*Currency, error) {
	bigInt := new(big.Int)
	bigIntamount, ok := bigInt.SetString(amount, 10)
	if !ok {
		return nil, errors.New("amount is not a valid number")
	}
	return &Currency{code: code, amount: bigIntamount}, nil
}

// Code gibt den Währungs-Code zurück.
func (c *Currency) Code() string {
	return c.code
}

// Amount gibt den Betrag der Währung zurück.
func (c *Currency) Amount() string {
	return c.amount.String()
}

// Add addiert zwei Currency und gibt eine neue Currency mit dem neuen Betrag zurück.
// Es wird der Währungs-Code des Objekts verwendet.
func (c *Currency) Add(currency *Currency) *Currency {
	result := new(big.Int)
	result.Add(c.amount, currency.amount)
	return &Currency{amount: result, code: c.code}
}

// Sub subtrahiert zwei Currency und gibt eine neue Currency mit dem neuen Betrag zurück.
// Es wird der Währungs-Code des Objekts verwendet.
func (c *Currency) Sub(currency *Currency) *Currency {
	result := new(big.Int)
	result.Sub(c.amount, currency.amount)
	return &Currency{amount: result, code: c.code}
}
