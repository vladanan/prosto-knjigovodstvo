// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/vladanan/prosto/src/controllers/i18n"
	"log"
	"net/http"
	"strings"
)

var key = i18n.SessionKey("session")

func getLang(r *http.Request) []string {
	// session, err := store.Get(r, "vezbamo.onrender.com-users")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
		log.Println("header: getLang: Error on get session:", err)
	}
	langMap := session.Values["language"]
	accept := r.Header.Get("Accept-Language")
	language := ""

	if langMap != nil {
		language = langMap.(string)
	}
	if language == "" {
		language = strings.Split(accept, ",")[0]
	}
	// var languageNames = map[string]string {
	// 	"ar": "Arabic    - العربية: ar",
	// 	"zh": "Chinese - 中文 (汉语): zh",
	// 	"en": "English  : en",
	// 	"sh": "Exyu      &nbsp;srpskohrvatski: sh",
	// 	"fr": "French   - français: fr",
	// 	"de": "German - Deutch: de",
	// 	"hi": "Hindi      - हिन्दी: hi",
	// 	"it": "Italian    ;- italiano: it",
	// 	"ru": "Russian  - русский: ru",
	// 	"sr": "Serbian  - српски: sr",
	// 	"es": "Spanish  - español: es",
	// }
	// fmt.Println("lang:", language, "accept:", accept, "0:", strings.Split(accept, ",")[0])
	// return []string{language, languageNames[language]}

	return []string{language, language}
}

var f i18n.F
var t = f.T // sluzi za lokalnu upotrebu za views paket
var T = f.T // sluzi za import za ostale pakete
func initt(r *http.Request) string {
	f.R = r
	return ""
}

func Header(r *http.Request) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(initt(r))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 59, Col: 13}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><div class=\"relative top-0 left-2 w-6\"><a href=\"/\" class=\"\"><img src=\"static/site/prosto.png\" height=\"25\" width=\"25\" alt=\"Vezbamo\"></a></div><a href=\"/\" class=\"absolute top-1 left-10 text-sm\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(t("Home"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 70, Col: 62}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a><div class=\"absolute top-0 right-2\"><select onchange=\"sendLang(event)\" name=\"lang\" id=\"lang\" class=\"dark:bg-black dark:border dark:border-slate-300 dark:px-1 w-12 text-sm \"><option class=\"font-bold \" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(getLang(r)[0])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 74, Col: 50}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(getLang(r)[1])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 74, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option> <option value=\"en\">en &nbsp;- English</option> <option value=\"sh\">sh &nbsp;&nbsp;- Srpskohrvatski - Ex-yu</option> <option disabled value=\"ar\">ar &nbsp;- Arabic - العربية</option> <option disabled value=\"zh\">zh &nbsp;- Chinese - 中文 (汉语)</option> <option disabled value=\"fr\">fr &nbsp;- French - français</option> <option disabled value=\"de\">de &nbsp;- German - Deutch</option> <option disabled value=\"hi\">hi &nbsp;- Hindi - हिन्दी</option> <option disabled value=\"it\">it &nbsp;- Italian - italiano</option> <option disabled value=\"ru\">ru &nbsp;- Russian - русский</option> <option disabled value=\"sr\">sr &nbsp;- Serbian - српски</option> <option disabled value=\"es\">es &nbsp;- Spanish - español</option></select> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if session, err := i18n.GetSessionFromContext(r.Context(), key); err == nil {
			if auth, _ := session.Values["authenticated"].(bool); !auth {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"text-sm ml-1 px-1\" type=\"button\"><a href=\"/sign_in\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(t("Sign_in"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 91, Col: 37}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></button> <button class=\"text-sm ml-1 px-1\" type=\"button\"><a href=\"/sign_up\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(t("Sign_up_b"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 94, Col: 39}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"text-sm ml-1 px-1\" type=\"button\"><a href=\"/sign_out\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(t("Sign_out"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 98, Col: 39}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></button> <button class=\"text-sm ml-1 px-1 dark:text-black bg-gradient-to-r from-blue-400 to-blue-200 rounded-sm\" type=\"button\"><a href=\"/dashboard\">&#64;</a></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
