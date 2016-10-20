package main

import(
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/unrolled/secure"
    "github.com/rs/cors"
    "github.com/codegangsta/negroni"
    "github.com/codegangsta/cli"
    "github.com/txthinking/law/api"
)

var domain string
var apiDomain string
var listen string
var origins []string

func main(){
    app := cli.NewApp()
    app.Name = "Law"
    app.Usage = "Law's room."
    app.Author = "Cloud"
    app.Email = "cloud@txthinking.com"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name: "domain",
            Value: "",
            Usage: "Your domain name.",
            Destination: &domain,
        },
        cli.StringFlag{
            Name: "apidomain",
            Value: "",
            Usage: "Your api domain name.",
            Destination: &apiDomain,
        },
        cli.StringFlag{
            Name: "listen",
            Value: "",
            Usage: "Listen address.",
            Destination: &listen,
        },
        cli.StringSliceFlag{
            Name: "origin",
            Usage: "Allow origins for CORS, can repeat more times.",
        },
    }
    app.Action = func(c *cli.Context) error {
        if c.String("domain") == "" {
            return cli.NewExitError("Domain is empty.", 86)
        }
        if c.String("apidomain") == "" {
            return cli.NewExitError("API domain is empty.", 86)
        }
        if c.String("listen") == "" {
            return cli.NewExitError("Listen address is empty.", 86)
        }
        origins = c.GlobalStringSlice("origin")
        return run()
    }
    app.Run(os.Args)
}

func run() error{
    r := mux.NewRouter()
    r.Host(domain).Methods("GET").Path("/signal/_/{id}").Handler(getSignalHandle(origins))
    r.Host(domain).Methods("GET").Path("/random").HandlerFunc(redirect)
    r.Host(domain).Methods("GET").Path("/room/{roomID}").HandlerFunc(redirect)
    r.Host(apiDomain).Methods("GET").Path("/v1/room/prepare").HandlerFunc(api.RoomPrepare)
    r.Host(apiDomain).Methods("POST").Path("/v1/room/status").HandlerFunc(api.RoomStatus)

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())
    n.Use(cors.New(cors.Options{
        AllowedOrigins: origins,
        AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
        AllowCredentials: true,
    }))
    n.Use(negroni.HandlerFunc(secure.New(secure.Options{
        AllowedHosts: []string{domain, apiDomain},
        SSLRedirect: false,
        SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
        STSSeconds: 315360000,
        STSIncludeSubdomains: true,
        STSPreload: true,
        FrameDeny: true,
        CustomFrameOptionsValue: "SAMEORIGIN",
        ContentTypeNosniff: true,
        BrowserXssFilter: true,
        ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval' blob: data: https://lawsroom.com wss://lawsroom.com https://dev-law.txthinking.com wss://dev-law.txthinking.com https://file.txthinking.com https://fonts.googleapis.com https://fonts.gstatic.com https://www.google-analytics.com",
    }).HandlerFuncWithNext))
    n.UseHandler(r)

    return http.ListenAndServe(listen, n)
}

