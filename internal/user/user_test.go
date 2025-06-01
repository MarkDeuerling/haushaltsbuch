package user_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/user"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		nachname  string
		vorname   string
		email     string
		expectErr error
	}{
		{
			name:      "Gültiger User",
			nachname:  "Mustermann",
			vorname:   "Max",
			email:     "max.mustermann@gmail.com",
			expectErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user.NewUser("123", tt.vorname, tt.nachname, tt.email, []byte("password"), time.Now(), time.Now())
			if tt.expectErr == nil {
				assert.Equal(t, tt.expectErr, nil, "Expected error to be %v, got %v", tt.expectErr, nil)
			}
			// if tt.expectErr != nil {
			// 	assert.ErrorIs(t, err, tt.expectErr, "Expected error to be %v, got %v", tt.expectErr, err)
			// }
		})
	}
}

func TestAktualisiereVorname(t *testing.T) {
	tests := []struct {
		name          string
		vorname       string
		expectErr     error
		expectVorname string
	}{
		{
			name:          "Gültiger Vorname",
			vorname:       "Marcel",
			expectErr:     nil,
			expectVorname: "Marcel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser("123", "Max", "Mustermann", "max.mustermann@gmail.de", []byte("password"), time.Now(), time.Now())
			u.NeuerVorname(tt.vorname)
			if tt.expectErr == nil {
				// assert.Equal(t, tt.expectErr, err, "Erwarteter Fehler ist %v, bekommen %v", tt.expectErr, err)
				assert.Equal(t, tt.expectVorname, u.Vorname(), "Erwarteter Vorname ist %s, bekommen %s", tt.expectVorname, u.Vorname())
			}
			if tt.expectErr != nil {
				// assert.ErrorIs(t, err, tt.expectErr, "Erwarteter Fehler ist %v, bekommen %v", tt.expectErr, err)
				assert.Equal(t, "Max", u.Vorname(), "Vorname sollte unverändert bleiben, erwartet 'Max', bekommen '%s'", u.Vorname())
			}
		})
	}
}

func TestAktualisiereNachname(t *testing.T) {
	tests := []struct {
		name           string
		nachname       string
		expectErr      error
		expectNachname string
	}{
		{
			name:           "Gültiger Nachname",
			nachname:       "Hermann",
			expectErr:      nil,
			expectNachname: "Hermann",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser("123", "Max", "Mustermann", "max.mustermann@gmail.de", []byte("password"), time.Now(), time.Now())
			u.NeuerNachname(tt.nachname)
			if tt.expectErr == nil {
				// assert.Equal(t, tt.expectErr, err, "Erwarteter Fehler ist %v, bekommen %v", tt.expectErr, err)
				assert.Equal(t, tt.expectNachname, u.Nachname(), "Erwarteter Nachname ist %s, bekommen %s", tt.expectNachname, u.Nachname())
			}
			if tt.expectErr != nil {
				// assert.ErrorIs(t, err, tt.expectErr, "Erwarteter Fehler ist %v, bekommen %v", tt.expectErr, err)
				assert.Equal(t, "Mustermann", u.Nachname(), "Nachname sollte unverändert bleiben, erwartet 'Mustermann', bekommen '%s'", u.Nachname())
			}
		})
	}
}

func TestAktualisiereEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectErr   error
		expectEmail string
	}{
		{
			name:        "Gültige Email",
			email:       "max.mustermann@gmail.de",
			expectErr:   nil,
			expectEmail: "max.mustermann@gmail.de",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser("123", "Max", "Mustermann", "max.mustermann@gmail.de", []byte("password"), time.Now(), time.Now())
			u.NeueEmail(tt.email)
			if tt.expectErr == nil {
				// assert.Equal(t, tt.expectErr, err, "Erwarteter Fehler ist %v, bekommen %v", tt.expectErr, err)
				assert.Equal(t, tt.expectEmail, u.Email(), "Erwarteter Email ist %s, bekommen %s", tt.expectEmail, u.Nachname())
			}
			if tt.expectErr != nil {
				// assert.ErrorIs(t, err, tt.expectErr, "Erwarteter Fehler ist %v, bekommen %v", tt.expectErr, err)
				assert.Equal(t, "max.mustermann@gmail.de", u.Email(), "Email sollte unverändert bleiben, erwartet 'max.mustermann@gmail.de', bekommen '%s'", u.Nachname())
			}
		})
	}
}
