package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/controllers/i18n"
	"github.com/vladanan/prosto/src/controllers/routes"
	"github.com/vladanan/prosto/src/models"
)

type IP string

type tracking struct {
	time  time.Time
	calls int
}

var cls map[IP]tracking
var bls map[IP]tracking
var claa map[IP]tracking
var blaa map[IP]tracking

func flodMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 192.168.0.1
		// ruter server notbuk wifi 25, notbuk lan 226
		// telefon 21, stari notbuk 26
		// ruter ip int 100.116.172.232
		// ruter ip out 178.220.151.98
		// 18 req u 1s učitavanje index
		// 1000 poziva u 0.118s za stranu custom_apis

		var lock = sync.Mutex{}

		ips, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("userip: %q is not IP:port", r.RemoteAddr)
		}
		ip := IP(ips)

		if os.Getenv("PRODUCTION") == "FALSE" {
			// fmt.Fprintf(os.Stderr, "...%v\n", r.URL)
			fmt.Println(
				time.Now().Format(time.RFC3339),
				"ip:", ip, "port:", port,
				"url:", r.URL,
				// "path:", r.URL.Path,
				"xff:", r.Header.Get("X-Forwarded-For"), "rf:", r.Header.Get("Referer"),
				"claa:", len(claa), "blaa:", len(blaa), fmt.Sprint(blaa),
				"cls:", len(cls), "bls:", len(bls))
		}

		if strings.Contains(r.URL.Path, "/api") ||
			strings.Contains(r.URL.Path, "/auth") ||
			strings.Contains(r.URL.Path, "/vmk") {

			// ako je poziv api/auth ne sme da se dozvoli više od dva poziva jer recimo za sign_up
			// jedan je poziv za url auth/sign_up_post a drugi za post /api/v/user što čini najmanje dva 0 i 1
			// kod dobijanje svih testova, prvo ide sajt poziv /htmx_get_tests i onda api /api/v/test
			// ali prilikom ucitavanja svega sa dashboarda to je sve pod auth rutom i aktivira zastitu a to je oko 5-6 fajlova
			calls := 6
			interval := int64(500)

			// for testing only
			// testId := "none"
			// if len(r.URL.Query()["call"]) != 0 {
			// 	testId = r.URL.Query()["call"][0]
			// }

			lock.Lock()
			defer lock.Unlock()

			// služi da se sve pobriše iz claa ako duže od n min nema upisa u claa
			enough := 0
			for _, MapRequest := range claa {
				// log.Println("aa:", time.Since(MapRequest.time).Seconds(), MapString, r.URL)
				if time.Since(MapRequest.time).Seconds() < 7 {
					// log.Println("IMA u claa skorijih", MapString, testId)
					enough++
				}
			}
			// fmt.Println("enough", enough) 12
			if enough == 0 {
				// log.Println("aa: nema u claa skorijih i sve se briše", testId, r.URL)
				claa = make(map[IP]tracking)
				// blaa = make(map[string]request) //ne treba svi da se brišu iz blaa nego samo pojedinačni u donjoj proveri
			}

			if b, found := blaa[ip]; found {
				// log.Println("aa: ima ga u bls: provera intervala", time.Since(b.time).Seconds(), ip, testId, b.calls)
				if time.Since(b.time).Seconds() > 3 {
					// log.Println("aa: prošlo više od 3sec za blaa del iz blaa i update claa", ip, testId, b.calls)
					delete(blaa, ip)
					claa[ip] = tracking{time: time.Now(), calls: b.calls}
					w.WriteHeader(http.StatusFound)
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte("0x3_aa found 0x01_2E"))
					// next.ServeHTTP(w, r)
				} else {
					// log.Println("BLaa  ***  MALICIUOS:", ip, testId, b.calls)
					// updatuje se i claa i blaa za novi malicious poziv inače bi im isteklo vreme i bili bi izbrisani
					blaa[ip] = tracking{time: time.Now(), calls: b.calls + 1}
					claa[ip] = tracking{time: time.Now(), calls: b.calls + 1}
					http.Error(
						w,
						"0x2_aa internal user error 0x02_9a",
						http.StatusForbidden)
				}
			} else if c, found := claa[ip]; found {
				// log.Println("aa: ima ga u claa", time.Since(c.time).Milliseconds(), ip, r.URL, testId, c.calls)
				if c.calls > calls && time.Since(c.time).Milliseconds() < interval {
					// log.Println("aa: CLaa pao na kriterijumu: IDE U BLaa", ip, testId, c.calls)
					blaa[ip] = tracking{time: time.Now(), calls: c.calls + 1}
					http.Error(
						w,
						"0x1_aa internal user error 0x02_9a",
						http.StatusForbidden)
				} else {
					// log.Println("aa: manje od kriterijuma: proverava se interval za claa", ip, testId)
					if time.Since(c.time).Seconds() > 3 {
						// log.Println("aa: interval za claa veći od 3 sve okej i briše se iz claa", ip)
						delete(claa, ip)
						next.ServeHTTP(w, r)
					} else {
						// log.Println("aa: interval za cl manji od 3 sve okej i updatuje se claa", ip, testId, c.calls)
						claa[ip] = tracking{time: time.Now(), calls: c.calls + 1}
						next.ServeHTTP(w, r)
					}
				}
			} else {
				// log.Println("aa: ip nije u claa dodaje se NOVI:", ip, testId)
				claa[ip] = tracking{time: time.Now(), calls: 1}
				next.ServeHTTP(w, r)
			}

		} else {

			// ako je običan poziv njih bude 18 u sec za index stranu i svi su po 0 milisekindi
			// dakle može puno njih u kratkom vremenu da bi se obične strane normalno učitavale 33
			calls := 50
			interval := int64(-1)

			// for testing only
			// testId := "none"
			// if len(r.URL.Query()["call"]) != 0 {
			// 	testId = r.URL.Query()["call"][0]
			// }

			lock.Lock()
			defer lock.Unlock()

			// služi da se sve pobriše iz cls ako duže od n sec nema upisa u cls
			enough := 0
			for _, MapRequest := range cls {
				// log.Println("s:", time.Since(MapRequest.time).Seconds(), MapString, r.URL)
				if time.Since(MapRequest.time).Seconds() < 7 {
					// log.Println("IMA u cls skorijih", MapString, testId)
					enough++
				}
			}
			// fmt.Println("enough", enough)
			if enough == 0 {
				// log.Println("s: nema u cls skorijih i sve se briše", testId, r.URL)
				cls = make(map[IP]tracking)
				// bls = make(map[string]request) //ne treba svi da se brišu iz bls nego samo pojedinačni u donjoj proveri
			}

			if b, found := bls[ip]; found {
				// log.Println("s: ima ga u bls: provera intervala", time.Since(b.time).Seconds(), ip, testId, b.calls)
				if time.Since(b.time).Seconds() > 3 {
					// log.Println("s: prošlo više od 3sec za bls del iz bls i update cls", ip, testId, b.calls)
					delete(bls, ip)
					cls[ip] = tracking{time: time.Now(), calls: b.calls}
					w.WriteHeader(http.StatusFound)
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte("0x3_s found 0x01_2E"))
					// next.ServeHTTP(w, r)
				} else {
					// log.Println("BLs  ***  MALICIUOS:", ip, testId, b.calls)
					// updatuje se i cls i bls za novi malicious poziv inače bi im isteklo vreme i bili bi izbrisani
					bls[ip] = tracking{time: time.Now(), calls: b.calls + 1}
					cls[ip] = tracking{time: time.Now(), calls: b.calls + 1}
					http.Error(
						w,
						"0x2_s internal user error 0x02_9a",
						http.StatusServiceUnavailable)
				}
			} else if c, found := cls[ip]; found {
				// log.Println("s: ima ga u cls", time.Since(c.time).Milliseconds(), ip, r.URL, testId, c.calls)
				if c.calls >= calls && time.Since(c.time).Milliseconds() > interval {
					// log.Println("s: CLs pao na kriterijumu: IDE U BLs", ip, testId, c.calls)
					bls[ip] = tracking{time: time.Now(), calls: c.calls + 1}
					http.Error(
						w,
						"0x1_s internal user error 0x02_9a",
						http.StatusServiceUnavailable)
				} else {
					// log.Println("manje od kriterijuma: proverava se interval za cls", ip, testId)
					if time.Since(c.time).Seconds() > 3 {
						// log.Println("s: interval za cls veći od 3 sve okej i briše se iz cls", ip)
						delete(cls, ip)
						next.ServeHTTP(w, r)
					} else {
						// log.Println("s: interval za cls manji od 3 sve okej i updatuje se cls", ip, testId, c.calls)
						cls[ip] = tracking{time: time.Now(), calls: c.calls + 1}
						next.ServeHTTP(w, r)
					}
				}
			} else {
				// log.Println("s: ip nije u cls dodaje se NOVI:", ip, testId)
				cls[ip] = tracking{time: time.Now(), calls: 1}
				next.ServeHTTP(w, r)
			}
		}

	})
}

var pool *pgxpool.Pool

func createLocalPool() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(os.Getenv("FEDORA_CONNECTION_STRING"))
	if err != nil {
		log.Println(err)
	}
	// config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	// 	// do something with every new connection
	// }
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Println(err)
	}
	// s := pool.Stat()
	// log.Println("konekcije:", s.MaxConns(), s.AcquiredConns())
	return pool
}
func insertLocalPoolConnMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := pool.Acquire(r.Context())
		if err != nil {
			log.Println(err)
		}
		s := pool.Stat()
		log.Println("local pool konekcije:", s.MaxConns(), s.AcquiredConns())
		connKey := models.ConnKey("pool_conn")
		ctx := context.WithValue(r.Context(), connKey, conn)
		rp := r.Clone(ctx)
		next.ServeHTTP(w, rp)
	})
}

var store *sessions.CookieStore

func insertSessionMiddleware(next http.Handler) http.Handler {
	// https://pkg.go.dev/github.com/gorilla/sessions@v1.2.2#section-documentation
	// https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-cookie-same-site-00
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
	// https://stackoverflow.com/questions/67821709/this-set-cookie-didnt-specify-a-samesite-attribute-and-was-default-to-samesi

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "vezbamo.onrender.com-users")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionKey := i18n.SessionKey("session")
		ctx := context.WithValue(r.Context(), sessionKey, session)
		// log.Println(session)
		rs := r.Clone(ctx)
		next.ServeHTTP(w, rs)
	})
}

func main() {

	var err_godotenv = godotenv.Load(".env")
	if err_godotenv != nil {
		log.Println(err_godotenv)
	}
	// key   = []byte("12345678901234567890123456789012")
	var key = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
		// SameSite: http.SameSiteDefaultMode,
		// SameSite: http.SameSiteLaxMode,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSite(0),
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r.NotFoundHandler = http.HandlerFunc(routes.GoTo404)

	routes.RouterSite(r)
	routes.RouterUsers(r)
	routes.RouterPausalDashboard(r)
	routes.RouterPausalForms(r)

	apirouter := r.PathPrefix("/api").Subrouter()
	authrouter := r.PathPrefix("/auth").Subrouter()

	// komentovati donje tri linije kada se ne koristi lokalni pool i obratno
	// pool = createLocalPool()
	// apirouter.Use(insertLocalPoolConnMiddleware) //123
	// authrouter.Use(insertLocalPoolConnMiddleware)

	routes.RouterCRUD(apirouter)
	routes.RouterAuth(authrouter)

	routes.RouterI18n(r)

	routes.ServeStatic(r, "/static/")

	cls = make(map[IP]tracking)
	bls = make(map[IP]tracking)
	claa = make(map[IP]tracking)
	blaa = make(map[IP]tracking)

	r.Use(flodMiddleware)

	r.Use(insertSessionMiddleware)

	log.Print(clr.Green + "main go" + clr.Reset)

	// go gamesForLearningChannelsAndLogs()

	server := &http.Server{
		Addr: "0.0.0.0:10000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      r,
	}

	// Run our server in a go routine so it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("server closed")
			} else {
				log.Println("error starting server, error:", err.Error())
				// os.Exit(1)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	// var err = srv.ListenAndServe("0.0.0.0:10000", r)
	// var err = server.ListenAndServe()
	// if errors.Is(err, http.ErrServerClosed) {
	// 	slog.Info("server closed")
	// } else if err != nil {
	// 	slog.Error("error starting server", "error", err.Error())
	// 	os.Exit(1)
	// }

}

// isključio u staticcheck "-U1000" da ne javlja za postojeće a neiskorišćene funkcije
func gamesForLearningChannelsAndLogs() {

	s := "jedan"
	cc := make(chan string)
	go func() {
		log.Println(s)
		dva := " dva"
		time.Sleep(time.Second * 0)
		// slog.Info(dva)
		cc <- s + dva
	}()
	slog.Info(<-cc)

}

/*

[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/ /*;q=0.8,application/signed-exchange;v=b3;q=0.7] Accept-Encoding:[gzip, deflate] Accept-Language:[en-US,en;q=0.9] Cache-Control:[max-age=0] Connection:[keep-alive] Referer:[http://192.168.0.25:10000/] Upgrade-Insecure-Requests:[1] User-Agent:[Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36 OPR/83.0.0.0]] {} <nil> 0 [] false 192.168.0.25:10000 map[] map[] <nil> map[] 192.168.0.21:41700 / <nil> <nil> <nil> 0xc0001b70e0 <nil> [] map[]}


 */
