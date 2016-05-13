package main

import(
    "net/http"
    "log"

    "github.com/gorilla/mux"
    "github.com/unrolled/secure"
    //"github.com/phyber/negroni-gzip/gzip"
    "github.com/rs/cors"
    "github.com/codegangsta/negroni"
    "github.com/txthinking/law/api"
)

func main(){
    r := mux.NewRouter()
    r.Handle("/signal/_/{id}", getSignalHandle())
    r.Methods("GET").Path("/random").HandlerFunc(redirect)
    r.Methods("GET").Path("/room/{id}").HandlerFunc(redirect)
    r.Methods("GET").Path("/unsupport").HandlerFunc(redirect)
    r.Methods("GET").Path("/v1/room/prepare").HandlerFunc(api.RoomPrepare)
    r.Methods("POST").Path("/v1/room/capacity").HandlerFunc(api.RoomCapacity)

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
    //n.Use(gzip.Gzip(gzip.DefaultCompression))
    //n.Use(negroni.NewStatic(http.Dir("public")))
    n.UseHandler(r)

    //go func() {
        if err := http.ListenAndServe(":1006", n); err != nil {
            log.Fatal("http", err)
        }
    //}()
    //if err := http.ListenAndServeTLS(":443", "./cert.pem", "./privkey.pem", n); err != nil {
        //log.Fatal("https", err)
    //}
}

