package i18n

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/vladanan/prosto/src/controllers/clr"
	"golang.org/x/text/language"
)

var l = clr.GetErrorLogger()

// posredni tip za key, value pair u context Value
type SessionKey string

func GetSessionFromContext(ctx context.Context, s SessionKey) (*sessions.Session, error) {
	if v := ctx.Value(s); v != nil {
		// log.Println("found value:", v)
		if p, ok := v.(*sessions.Session); ok {
			// log.Println("jeste pointer za store")
			return p, nil
		} else {
			return nil, fmt.Errorf("poslat je pointer koji nije tip sessions.Store")
		}
	} else {
		return nil, fmt.Errorf("request context key nije validan: %v", s)
	}
}

type F struct {
	R *http.Request
}

func (f *F) T(tkey string) string {

	key := SessionKey("session")
	session, err := GetSessionFromContext(f.R.Context(), key)
	if err != nil {
		l(f.R, 7, err)
		return ""
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err = bundle.LoadMessageFile("assets/site/active.en.toml")
	if err != nil {
		l(f.R, 8, err)
		return ""
	}
	_, err = bundle.LoadMessageFile("assets/site/active.sh.toml")
	if err != nil {
		l(f.R, 8, err)
		return ""
	}
	// bundle.MustLoadMessageFile("assets/i18n/active.en.toml")
	// bundle.MustLoadMessageFile("assets/i18n/active.sh.toml")
	// bundle.MustLoadMessageFile("assets/i18n/active.es.toml")

	langMap := session.Values["language"]
	sessionLanguage := ""

	if langMap != nil {
		sessionLanguage = langMap.(string)
	}

	lang := f.R.FormValue("lang")
	accept := f.R.Header.Get("Accept-Language")

	if sessionLanguage != "" {
		accept = sessionLanguage
	}

	//fmt.Println("language: ", lang, "header: ", accept) ,jk

	localizer := i18n.NewLocalizer(bundle, lang, accept)

	// prevod := localizer.MustLocalize(&i18n.LocalizeConfig{
	prevod, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: tkey, //id je u stvari key u key value pair u toml fajlu
		},
	})
	if err != nil {
		l(f.R, 8, err)
		fe := strings.Split(strings.ReplaceAll(err.Error(), "\"", ""), " ")
		return fe[6] + ": TRANSLATION ERROR FOR: " + fe[1]
	}

	return prevod

}

// get sa receiverom tako da se u templ stranicama lakse napravi kratka t funkcija
// mada je u templ to odradjeno direktno preko instance F structa tako da je ova func nepotrebna osim ako zatreba negde van views sto je malo verovatno
func GetT(r *http.Request) func(tkey string) string {
	f := F{R: r}
	return f.T
}
