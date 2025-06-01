package user

import (
	"time"
)

// ID repräsentiert die ID des Users.
type ID = string

// User repäsentiert einen registrierten Benutzer im Haushaltsbuchsystem.
type User struct {
	iD              ID
	vorname         string
	nachname        string
	Aktiv           bool
	email           string
	passwort        []byte
	erstelltAm      time.Time
	aktuallisiertAm time.Time
}

// NewUser erzeugt einen neuen User mit expliziten Parametern.
func NewUser(id ID, vorname, nachname, email string, passwort []byte, erstelltAm, aktualisiertAm time.Time) *User {
	return &User{
		iD:              id,
		vorname:         vorname,
		nachname:        nachname,
		email:           email,
		passwort:        passwort,
		erstelltAm:      erstelltAm,
		aktuallisiertAm: aktualisiertAm,
	}
}

// ID gibt die ID des Users zurück.
func (u *User) ID() ID {
	return u.iD
}

// Vorname gibt den Vornamen des Users zurück.
func (u *User) Vorname() string {
	return u.vorname
}

// NeuerVorname aktualisiert den Vornamen des Users und validiert ihn.
func (u *User) NeuerVorname(vorname string) {
	u.vorname = vorname
}

// Nachname gibt den Nachnamen des Users zurück.
func (u *User) Nachname() string {
	return u.nachname
}

// NeuerNachname aktualisiert den Vornamen des Users und validiert ihn.
func (u *User) NeuerNachname(nachname string) {
	u.nachname = nachname
}

// IstAktiv gibt zurück, ob der User aktiv ist.
func (u *User) IstAktiv() bool {
	return u.Aktiv
}

// ErstelltAm gibt den Erstellungszeitpunkt des Users zurück.
func (u *User) ErstelltAm() time.Time {
	return u.erstelltAm
}

// AktualisiertAm gibt den Aktualisierungszeitpunkt des Users zurück.
func (u *User) AktualisiertAm() time.Time {
	return u.aktuallisiertAm
}

// Email gibt die Email des Users zurück.
func (u *User) Email() string {
	return u.email
}

// NeueEmail aktualisiert die Email des Users und validiert sie.
func (u *User) NeueEmail(email string) {
	u.email = email
}

// Passwort gibt das Passwort des Users zurück.
func (u *User) Passwort() []byte {
	return u.passwort
}

// NeuesPasswort aktualisiert das Passwort des Users.
func (u *User) NeuesPasswort(passwort []byte) {
	u.passwort = passwort
}

// Aktiviert aktiviert den User.
func (u *User) Aktiviert() {
	u.Aktiv = true
}

// Deaktiviert deaktiviert den User.
func (u *User) Deaktiviert() {
	u.Aktiv = false
}

// Aktualisert aktualisert den Aktualisierungszeitpunkt des Users.
func (u *User) Aktualisert() {
	u.aktuallisiertAm = time.Now().UTC()
}
