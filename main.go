package main

import(
    "net/http"
    "log"
    "os"

    "github.com/gorilla/mux"
    "github.com/unrolled/secure"
    "github.com/phyber/negroni-gzip/gzip"
    "github.com/rs/cors"
    "github.com/codegangsta/negroni"
    "git.txthinking.com/txthinking/signal"
)

func main(){
    r := mux.NewRouter()
    s := signal.New(func(r *http.Request) bool {
        allows := []string{
            "https://law.txthinking.com:444",
            "http://127.0.0.1:1634",
        }
        origin := r.Header.Get("Origin")
        for _, v := range allows {
            if v == origin {
                return true
            }
        }
        return false
    }, &Dog{})
    r.Handle("/signal/{id}", s)
    r.Methods("GET").Path("/hello").HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("[\"World\"]"))
    })

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())
    n.Use(cors.New(cors.Options{
        AllowedOrigins: []string{"https://law.txthinking.com:444", "http://127.0.0.1:1634"},
        AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
        AllowCredentials: true,
    }))
    n.Use(negroni.HandlerFunc(secure.New(secure.Options{
        AllowedHosts: []string{"law.txthinking.com:444"},
        SSLRedirect: true,
        SSLHost: "law.txthinking.com",
        SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
        STSSeconds: 315360000,
        STSIncludeSubdomains: true,
        STSPreload: true,
        FrameDeny: true,
        CustomFrameOptionsValue: "SAMEORIGIN",
        ContentTypeNosniff: true,
        BrowserXssFilter: true,
        ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval' blob: https://law.txthinking.com:444 wss://law.txthinking.com:444 https://fonts.googleapis.com https://fonts.gstatic.com https://webrtc.github.io https://ajax.googleapis.com",
    }).HandlerFuncWithNext))
    n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.Use(negroni.NewStatic(http.Dir("public")))
    n.UseHandler(r)

    go func() {
        //if err := http.ListenAndServe(":80", n); err != nil {
            //log.Fatal("http", err)
        //}
    }()
    cert := "/etc/letsencrypt/live/law.txthinking.com/cert.pem"
    privkey := "/etc/letsencrypt/live/law.txthinking.com/privkey.pem"
    if _, err := os.Open(cert); err != nil {
        cert = "./cert.pem"
    }
    if _, err := os.Open(privkey); err != nil {
        privkey = "./privkey.pem"
    }
    if err := http.ListenAndServeTLS(":444", cert, privkey, n); err != nil {
        log.Fatal("https", err)
    }
}

