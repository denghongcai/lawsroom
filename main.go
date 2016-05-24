package main

import(
    "net/http"
    "log"

    "github.com/gorilla/mux"
    "github.com/unrolled/secure"
    "github.com/rs/cors"
    "github.com/codegangsta/negroni"
    "github.com/txthinking/law/api"
)

func main(){
    r := mux.NewRouter()
    r.Host("law.txthinking.com").Methods("GET").Path("/signal/_/{id}").Handler(getSignalHandle())
    r.Host("law.txthinking.com").Methods("GET").Path("/random").HandlerFunc(redirect)
    r.Host("law.txthinking.com").Methods("GET").Path("/room/{id}").HandlerFunc(redirect)
    r.Host("api.law.txthinking.com").Methods("GET").Path("/v1/room/prepare").HandlerFunc(api.RoomPrepare)
    r.Host("api.law.txthinking.com").Methods("POST").Path("/v1/room/status").HandlerFunc(api.RoomStatus)

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())
    n.Use(cors.New(cors.Options{
        AllowedOrigins: []string{"https://law.txthinking.com", "https://127.0.0.1"},
        AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
        AllowCredentials: true,
    }))
    n.Use(negroni.HandlerFunc(secure.New(secure.Options{
        AllowedHosts: []string{"law.txthinking.com", "api.law.txthinking.com"},
        SSLRedirect: false,
        SSLHost: "law.txthinking.com",
        SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
        STSSeconds: 315360000,
        STSIncludeSubdomains: true,
        STSPreload: true,
        FrameDeny: true,
        CustomFrameOptionsValue: "SAMEORIGIN",
        ContentTypeNosniff: true,
        BrowserXssFilter: true,
        ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval' blob: data: https://law.txthinking.com wss://law.txthinking.com https://file.txthinking.com https://fonts.googleapis.com https://fonts.gstatic.com https://www.google-analytics.com https://127.0.0.1",
    }).HandlerFuncWithNext))
    n.UseHandler(r)

    if err := http.ListenAndServe(":1008", n); err != nil {
        log.Fatal("http", err)
    }
}

