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
    r.Host("lawsroom.com").Methods("GET").Path("/signal/_/{id}").Handler(getSignalHandle())
    r.Host("lawsroom.com").Methods("GET").Path("/random").HandlerFunc(redirect)
    r.Host("lawsroom.com").Methods("GET").Path("/room/{id}").HandlerFunc(redirect)
    r.Host("lawsroom.com").Methods("GET").Path("/support").HandlerFunc(redirect)
    r.Host("api.lawsroom.com").Methods("GET").Path("/v1/room/prepare").HandlerFunc(api.RoomPrepare)
    r.Host("api.lawsroom.com").Methods("POST").Path("/v1/room/status").HandlerFunc(api.RoomStatus)

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())
    n.Use(cors.New(cors.Options{
        AllowedOrigins: []string{"https://lawsroom.com", "https://127.0.0.1"},
        AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
        AllowCredentials: true,
    }))
    n.Use(negroni.HandlerFunc(secure.New(secure.Options{
        AllowedHosts: []string{"lawsroom.com", "api.lawsroom.com"},
        SSLRedirect: false,
        SSLHost: "lawsroom.com",
        SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
        STSSeconds: 315360000,
        STSIncludeSubdomains: true,
        STSPreload: true,
        FrameDeny: true,
        CustomFrameOptionsValue: "SAMEORIGIN",
        ContentTypeNosniff: true,
        BrowserXssFilter: true,
        ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval' blob: data: https://lawsroom.com wss://lawsroom.com https://fonts.googleapis.com https://fonts.gstatic.com https://www.google-analytics.com https://127.0.0.1",
    }).HandlerFuncWithNext))
    n.UseHandler(r)

    if err := http.ListenAndServe(":1006", n); err != nil {
        log.Fatal("http", err)
    }
}

