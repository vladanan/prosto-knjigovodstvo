package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"golang.org/x/crypto/bcrypt"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/views/cp/emails"
)

var l = clr.GetErrorLogger()

/*


MAIL


*/

// sluzi za upisivanje izlaza iz templ komponenti u html tekst
// napravljen da bi se email tepml komponenta upisivala u html koji se salje iz smtp paketa
// Html.Write implementira interfejs io.Writer koji zahteva templ Render za upis html-a
type Html struct {
	Text string
}

func (h *Html) Write(p []byte) (i int, err error) {
	h.Text = string(p)
	// log.Println("html length:", len(h.text))
	return len(p), nil
}

type Mail interface {
	GetEmailData() (*http.Request, string, string, error)
}

type MailRegister struct {
	R        *http.Request
	Email    string
	UserName string
	Url      string
}

func (mr MailRegister) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.VerifyEmailForRegister(mr.R, mr.UserName, mr.Url).Render(mr.R.Context(), &Html); err != nil {
		return mr.R, "", "", l(mr.R, 7, err)
	}
	return mr.R, mr.Email, Html.Text, nil
}

type MailForgottenPasswordRequest struct {
	R        *http.Request
	Email    string
	UserName string
	Url      string
}

func (fpr MailForgottenPasswordRequest) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.VerifyForgottenPasswordRequest(fpr.R, fpr.UserName, fpr.Url).Render(fpr.R.Context(), &Html); err != nil {
		return fpr.R, "", "", l(fpr.R, 7, err)
	}
	return fpr.R, fpr.Email, Html.Text, nil
}

type MailFPSendNewPassword struct {
	R        *http.Request
	Email    string
	UserName string
	Password string
	Url      string
}

func (snp MailFPSendNewPassword) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.FPSendNewPassword(snp.R, snp.UserName, snp.Password, snp.Url).Render(snp.R.Context(), &Html); err != nil {
		return snp.R, "", "", l(snp.R, 7, err)
	}
	return snp.R, snp.Email, Html.Text, nil
}

type DeleteUserRequest struct {
	R        *http.Request
	Email    string
	UserName string
	Url      string
}

func (du DeleteUserRequest) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.DeleteUserRequest(du.R, du.UserName, du.Url).Render(du.R.Context(), &Html); err != nil {
		return du.R, "", "", l(du.R, 7, err)
	}
	return du.R, du.Email, Html.Text, nil
}

type DeletedUser struct {
	R        *http.Request
	Email    string
	UserName string
}

func (du DeletedUser) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.DeletedUser(du.R, du.UserName).Render(du.R.Context(), &Html); err != nil {
		return du.R, "", "", l(du.R, 7, err)
	}
	return du.R, du.Email, Html.Text, nil
}

type ChangeEmailRequest struct {
	R        *http.Request
	Email    string
	NewEmail string
	UserName string
	Url      string
}

func (cm ChangeEmailRequest) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.ChangeEmailRequest(cm.R, cm.UserName, cm.Email, cm.NewEmail, cm.Url).Render(cm.R.Context(), &Html); err != nil {
		return cm.R, "", "", l(cm.R, 7, err)
	}
	return cm.R, cm.Email, Html.Text, nil
}

type ChangedEmail struct {
	R        *http.Request
	Email    string
	UserName string
}

func (du ChangedEmail) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.ChangedEmail(du.R, du.UserName).Render(du.R.Context(), &Html); err != nil {
		return du.R, "", "", l(du.R, 7, err)
	}
	return du.R, du.Email, Html.Text, nil
}

type ChangeNameRequest struct {
	R           *http.Request
	Email       string
	UserName    string
	NewUserName string
	Url         string
}

func (cm ChangeNameRequest) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.ChangeNameRequest(cm.R, cm.UserName, cm.NewUserName, cm.Url).Render(cm.R.Context(), &Html); err != nil {
		return cm.R, "", "", l(cm.R, 7, err)
	}
	return cm.R, cm.Email, Html.Text, nil
}

type ChangePasswordRequest struct {
	R        *http.Request
	Email    string
	UserName string
	Url      string
}

func (cm ChangePasswordRequest) GetEmailData() (*http.Request, string, string, error) {
	var Html Html
	if err := emails.ChangePasswordRequest(cm.R, cm.UserName, cm.Url).Render(cm.R.Context(), &Html); err != nil {
		return cm.R, "", "", l(cm.R, 7, err)
	}
	return cm.R, cm.Email, Html.Text, nil
}

func SendEmail(m Mail) error {
	r, email, html, err := m.GetEmailData()
	if err != nil {
		return l(r, 4, err)
	}
	// Set up authentication information.
	auth := sasl.NewPlainClient("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_APP_PASSWORD_VEZBAMO"))

	// Connect to the server, authenticate, set the sender and recipient and send the email all in one step.
	to := []string{email}
	msg := strings.NewReader(
		`Content-Transfer-Encoding: quoted-printable` + "\r\n" +
			`Content-Type: text/html; charset="UTF-8"` + "\r\n" +
			`To: ` + email + "\r\n" +
			`Subject: Vezbamo portal` + "\r\n" +
			"\r\n" +
			html +
			"\r\n")
	if err = smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("SMTP_EMAIL"), to, msg); err != nil {
		l(r, 8, err)
		return l(r, 7, clr.NewAPIError(http.StatusNotAcceptable, "Verify_email_not_sent"))
	}

	return nil
}

/*


HASH & ENCRYPT


*/

// Pravi url key bilo vremenski sa opcionim dodatnim stringom ili neograniceni samo od stringa
func Get64UrlKey(check string, min int) string {
	var s string
	switch min {
	case 1:
		// seče datum-vreme string slice na 16. poziciji tako da ostaju jedinice minuta
		s = "auth-bot-stop-" + check + time.Now().Format(time.DateTime)[:16]
	case 10:
		// seče datum-vreme string slice na 15. poziciji tako da ostaju desetice minuta
		s = "auth-bot-stop-" + check + time.Now().Format(time.DateTime)[:15]
	case 0:
		s = "auth-bot-stop-" + check
	default:
		log.Println("pogresan parametar za funkciju za izradu key")
		return ""
	}
	ciphertext, err := bcrypt.GenerateFromPassword([]byte(s), 7)
	if err != nil {
		log.Println(err)
	}
	return base64.RawURLEncoding.EncodeToString(ciphertext)
}

// Proverava da li se base64 url crypto key slaze sa lozinkom tj. tekstom za proveru
// zbog granicnih situacija da se key napravi u sekundama pri kraju minuta ili minutima pri kraju desetice minuta
// za key od 1min proverava se i jedan minut manje a za kay od 10 i 10 minuta manje
// func dobija key za proveru, string u odnosu na koji se radi provera
// i broj da li da se proverava sa generisanim stringom za 1min ili 10min a ako je 0 onda nema vremenskog ogranicenja
func Check64UrlKey(key64, check string, min int) error {
	if data, err := base64.RawURLEncoding.DecodeString(key64); err != nil {
		return err
	} else {
		switch min {
		case 1:
			t := time.Now()
			t1 := t.Add(-(time.Minute * 1))
			s := "auth-bot-stop-" + check + t.Format(time.DateTime)[:16]
			s1 := "auth-bot-stop-" + check + t1.Format(time.DateTime)[:16]
			// dovoljno je da bar s ili s-1 bude okej
			if err := bcrypt.CompareHashAndPassword(data, []byte(s)); err != nil {
				if err := bcrypt.CompareHashAndPassword(data, []byte(s1)); err != nil {
					return err
				}
			}
		case 10:
			t := time.Now()
			t10 := t.Add(-(time.Minute * 10))
			s := "auth-bot-stop-" + check + t.Format(time.DateTime)[:15]
			s10 := "auth-bot-stop-" + check + t10.Format(time.DateTime)[:15]
			// dovoljno je da bar s ili s-10 bude okej
			if err := bcrypt.CompareHashAndPassword(data, []byte(s)); err != nil {
				if err := bcrypt.CompareHashAndPassword(data, []byte(s10)); err != nil {
					return err
				}
			}
		case 0:
			s := "auth-bot-stop-" + check
			if err := bcrypt.CompareHashAndPassword(data, []byte(s)); err != nil {
				return err
			}
		default:
			return errors.New("pogresni argumenti za proveru 64 url key")
		}
	}
	return nil
}

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		// panic(err)
		log.Println(err)
	}
	return &key
}

// Pravi key od datumskog stringa koji moze da se ponovi tako da dve razlicite func mogu da naprave isti key i dekodiraju tekst,
// uzima datum daleke godine ~200k godina u buducnosti sa nulovanim satom, minutima i sekundama
// kako bi se dobio cist datum za period od jednog dana i dodate su mikrosekunde radi kompleksnosti,
// uzima se deo od 16 bajta tj. deo bez zadnje tri nule (/1000) i spaja se sa specificim api delom kljuca takodje od 16 bajta
// koji se koristi i kao additionalData pa se i od njih pravi pointer na 32bitni niz
func DateEncryptionKey(apiAddZaKljuc string) *[32]byte {
	t1 := time.Date(294180, time.Now().Month(), time.Now().Day(), 0, 0, 0, 29418000000, time.UTC)
	s := fmt.Sprint(t1.UnixMicro()/1000) + apiAddZaKljuc
	s32 := [32]byte{}
	_, err := io.ReadFull(strings.NewReader(s), s32[:])
	if err != nil {
		log.Println("greska reader", err)
	}
	return &s32
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext, additionalData []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, additionalData), nil
}

// Wrap za utils.Encrypt tako da prihvata string sa porukama i vraca base64 url string
// koristi i specificni api deo kljuca koji se koristi i kao additionalData
func MsgEncrypt(msg, additionalData string, key *[32]byte) (string, error) {
	ciphertext, err := Encrypt([]byte(msg), []byte(additionalData), key)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext, additionalData []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		err := errors.New("malformed ciphertext")
		log.Println(err)
		return nil, err
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		additionalData,
	)
}

// Wrap za utils.Decrypt tako da prihvata base64 string i vraca slice dekriptovanih string poruka
// kojima prethodi specificni api deo kljuca koji se koristi i kao additionalData
func MsgDecrypt(msg, additionalData string, key *[32]byte) ([]string, error) {
	poruke := []string{additionalData}
	dec64, err := base64.RawURLEncoding.DecodeString(msg)
	if err != nil {
		return poruke, err
	}
	plaintext, err := Decrypt(dec64, []byte(additionalData), key)
	if err != nil {
		return poruke, err
	}
	poruke = append(poruke, strings.Split(string(plaintext), "||")...)
	if len(poruke) < 1 {
		return poruke, errors.New("greska u kreiranju poruka iz url keys")
	}
	return poruke, nil
}

// Pravi novu random lozinku i bcrypt key, a ako je poslata i lozinka onda daje nju i bcrypt key za nju
// poslata lozinka mora prethodno da se verifikuje sa vet.ValidatePassword()
// ali ne ovde jer se desava greska: import cycle not allowed
func GetNewPassword(password string) (string, string, error) {
	if password == "" {
		key := [32]byte{}
		if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
			log.Println(err)
			return "", "", err
		}
		keys := fmt.Sprint(key)
		password = base64.RawURLEncoding.EncodeToString([]byte(keys))[10:18]
	}
	if ciphertextFromPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12); err != nil {
		log.Println(err)
		return "", "", err
	} else {
		log.Println("get new pss i key", password, string(ciphertextFromPassword))
		return password, string(ciphertextFromPassword), nil
	}
}
