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
    signal.ROOM_CAPACITY = 5
    s := signal.New(func(r *http.Request) bool {
        allows := []string{
            "https://lawsroom.com",
            "https://127.0.0.1",
        }
        origin := r.Header.Get("Origin")
        for _, v := range allows {
            if v == origin {
                return true
            }
        }
        return false
    }, nil)
    r.Handle("/signal/{id}", s)
    r.Methods("GET").Path("/hello").HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("[\"World\"]"))
    })

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())
    n.Use(cors.New(cors.Options{
        AllowedOrigins: []string{"https://lawsroom.com", "https://127.0.0.1"},
        AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
        AllowCredentials: true,
    }))
    n.Use(negroni.HandlerFunc(secure.New(secure.Options{
        AllowedHosts: []string{"lawsroom.com"},
        SSLRedirect: true,
        SSLHost: "lawsroom.com",
        SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
        STSSeconds: 315360000,
        STSIncludeSubdomains: true,
        STSPreload: true,
        FrameDeny: true,
        CustomFrameOptionsValue: "SAMEORIGIN",
        ContentTypeNosniff: true,
        BrowserXssFilter: true,
        ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval' blob: data: https://lawsroom.com https://lawsroom.com:1443 wss://lawsroom.com wss://lawsroom.com:1443 https://fonts.googleapis.com https://fonts.gstatic.com https://www.google-analytics.com https://127.0.0.1",
    }).HandlerFuncWithNext))
    n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.Use(negroni.NewStatic(http.Dir("public")))
    n.UseHandler(r)

    //go func() {
        //if err := http.ListenAndServe(":80", n); err != nil {
            //log.Fatal("http", err)
        //}
    //}()
    cert := "/etc/letsencrypt/live/lawsroom.com/cert.pem"
    privkey := "/etc/letsencrypt/live/lawsroom.com/privkey.pem"
    if _, err := os.Open(cert); err != nil {
        cert = "./cert.pem"
    }
    if _, err := os.Open(privkey); err != nil {
        privkey = "./privkey.pem"
    }
    if err := http.ListenAndServeTLS(":443", cert, privkey, n); err != nil {
        log.Fatal("https", err)
    }
}

