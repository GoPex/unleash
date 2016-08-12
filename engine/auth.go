package engine

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin/render"

	// plugin package
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	"gopkg.in/authboss.v0"
	// register authboss register module
	_ "gopkg.in/authboss.v0/register"
	// register authboss login module
	_ "gopkg.in/authboss.v0/auth"
	// to confirm authboss
	_ "gopkg.in/authboss.v0/confirm"
	// to lock user after N authentication failures
	_ "gopkg.in/authboss.v0/lock"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/remember"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

var ab *authboss.Authboss

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "username",
		authboss.FlashSuccessKey: ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

func initAuthBossPolicy(ab *authboss.Authboss) {
	ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       4,
			MaxLength:       8,
			AllowWhitespace: false,
		},
	}
}

func initAuthBossLayout(ab *authboss.Authboss, r *gin.Engine) {
	if os.Getenv(gin.ENV_GIN_MODE) == gin.ReleaseMode {
		ab.Layout = r.HTMLRender.(render.HTMLProduction).Template.Funcs(funcs).Lookup("authboss")
	} else {
		html := r.HTMLRender.(render.HTMLDebug).Instance("authboss.tmpl", nil).(render.HTML)
		ab.Layout = html.Template.Funcs(template.FuncMap(funcs)).Lookup("authboss.tmpl")
	}
}

var database = NewMemStorer()

func initAuthBossParam(r *gin.Engine) *authboss.Authboss {
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = securecookie.New(cookieStoreKey, nil)
	sessionStore = sessions.NewCookieStore(sessionStoreKey)

	ab = authboss.New()
	ab.Storer = database
	ab.CookieStoreMaker = NewCookieStorer
	ab.SessionStoreMaker = NewSessionStorer
	ab.ViewsPath = filepath.Join("ab_views")

	ab.LayoutDataMaker = layoutData

	ab.MountPath = "/auth"
	ab.LogWriter = os.Stdout

	ab.XSRFName = "csrf_token"
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	initAuthBossLayout(ab, r)
	ab.Mailer = authboss.LogMailer(os.Stdout)
	initAuthBossPolicy(ab)

	if err := ab.Init(); err != nil {
		// Handle error, don't let program continue to run
		log.Fatalln(err)
	}
	return ab
}
