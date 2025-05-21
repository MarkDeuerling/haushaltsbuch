package user

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrInvalideEmail wird zurückgegeben, wenn die Email ungültig ist.
	ErrInvalideEmail = errors.New("Ungültige Email")
	// ErrLeererVorname wird zurückgegeben, wenn der Vorname ungültig ist.
	ErrLeererVorname = errors.New("Vorname darf nicht leer sein")
	// ErrLeererNachname wird zurückgegeben, wenn der Nachname ungültig ist.
	ErrLeererNachname = errors.New("Nachname darf nicht leer sein")
	// ErrVornameZuKurz wird zurückgegeben, wenn der Vorname zu kurz ist.
	ErrVornameZuKurz = errors.New("Vorname muss mindestens 1 Zeichen lang sein")
	// ErrVornameZuLang wird zurückgegeben, wenn der Vorname zu lang ist.
	ErrVornameZuLang = errors.New("Vorname darf maximal 100 Zeichen lang sein")
	// ErrNachnameZuKurz wird zurückgegeben, wenn der Vorname zu kurz ist.
	ErrNachnameZuKurz = errors.New("Nachname muss mindestens 1 Zeichen lang sein")
	// ErrNachnameZuLang wird zurückgegeben, wenn der Vorname zu lang ist.
	ErrNachnameZuLang = errors.New("Nachname darf maximal 100 Zeichen lang sein")
)

// ID repräsentiert die ID des Users.
type ID = string

// User repäsentiert einen registrierten Benutzer im Haushaltsbuchsystem.
type User struct {
	iD              ID
	vorname         string
	nachname        string
	logo            string
	Aktiv           bool
	email           string
	passwort        []byte
	erstelltAm      time.Time
	aktuallisiertAm time.Time
}

// NewUser erzeugt einen neuen User.
func NewUser(id ID, vorname, nachname, logo string, email string, password []byte) (*User, error) {
	jetzt := time.Now().UTC()

	user := &User{
		iD:              id,
		vorname:         vorname,
		nachname:        nachname,
		logo:            logo,
		email:           email,
		passwort:        password,
		erstelltAm:      jetzt,
		aktuallisiertAm: jetzt,
	}
	if err := user.validiereVorname(vorname); err != nil {
		return nil, err
	}
	if err := user.validiereNachname(nachname); err != nil {
		return nil, err
	}
	if ok := user.validiereEmail(); !ok {
		return nil, ErrInvalideEmail
	}

	return user, nil
}

// ID gibt die ID des Users zurück.
func (u *User) ID() ID {
	return u.iD
}

// Vorname gibt den Vornamen des Users zurück.
func (u *User) Vorname() string {
	return u.vorname
}

func (u *User) validiereVorname(vorname string) error {
	vorname = strings.TrimSpace(vorname)
	if vorname == "" {
		return ErrLeererVorname
	}

	if len(vorname) < 1 {
		return ErrVornameZuKurz
	}

	if len(vorname) > 100 {
		return ErrVornameZuLang
	}
	return nil
}

func (u *User) validiereNachname(nachname string) error {
	nachname = strings.TrimSpace(nachname)
	if nachname == "" {
		return ErrLeererNachname
	}

	if len(nachname) < 1 {
		return ErrNachnameZuKurz
	}

	if len(nachname) > 100 {
		return ErrNachnameZuLang
	}
	return nil
}

// AktualisiereVorname aktualisiert den Vornamen des Users und validiert ihn.
func (u *User) AktualisiereVorname(vorname string) error {
	err := u.validiereVorname(vorname)
	if err != nil {
		return err
	}
	u.vorname = vorname
	return nil
}

// Nachname gibt den Nachnamen des Users zurück.
func (u *User) Nachname() string {
	return u.nachname
}

// AktualisiereNachname aktualisiert den Vornamen des Users und validiert ihn.
func (u *User) AktualisiereNachname(nachname string) error {
	err := u.validiereNachname(nachname)
	if err != nil {
		return err
	}
	u.nachname = nachname
	return nil
}

// Logo gibt das Logo des Users zurück.
func (u *User) Logo() string {
	return u.logo
}

// IstAktiv gibt zurück, ob der User aktiv ist.
func (u *User) IstAktiv() bool {
	return u.Aktiv
}

// ErstelltAm gibt den Erstellungszeitpunkt des Users zurück.
func (u *User) ErstelltAm() time.Time {
	return u.erstelltAm
}

// SetzeErstelltAm setzt den Erstellungszeitpunkt des Users.
func (u *User) SetzeErstelltAm(erstellt time.Time) {
	u.erstelltAm = erstellt
}

// AktualisiertAm gibt den Aktualisierungszeitpunkt des Users zurück.
func (u *User) AktualisiertAm() time.Time {
	return u.aktuallisiertAm
}

// Email gibt die Email des Users zurück.
func (u *User) Email() string {
	return u.email
}

// Passwort gibt das Passwort des Users zurück.
func (u *User) Passwort() []byte {
	return u.passwort
}

// AktualisiereEmail aktualisiert die Email des Users und validiert sie.
func (u *User) AktualisiereEmail(email string) error {
	if !u.validiereEmail() {
		return ErrInvalideEmail
	}
	return nil
}

func (u *User) validiereEmail() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(u.email)
}

// AktualisierePasswort aktualisiert das Passwort des Users.
func (u *User) AktualisierePasswort(passwort []byte) {
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

// Aktualisere aktualisert den Aktualisierungszeitpunkt des Users.
func (u *User) Aktualisere() {
	u.aktuallisiertAm = time.Now().UTC()
}
