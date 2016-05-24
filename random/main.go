package main

import(
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/unrolled/secure"
    "github.com/rs/cors"
    "github.com/codegangsta/negroni"
    "github.com/codegangsta/cli"
)

var domain string
var listen string
var origins []string

func main(){
    app := cli.NewApp()
    app.Name = "Random of Law"
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
    r.Host(domain).Methods("GET").Path("/signal/r/{id}").Handler(getSignalHandle(origins))

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())
    n.Use(cors.New(cors.Options{
        AllowedOrigins: origins,
        AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
        AllowCredentials: true,
    }))
    n.Use(negroni.HandlerFunc(secure.New(secure.Options{
        AllowedHosts: []string{domain},
        SSLRedirect: false,
        SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
        STSSeconds: 315360000,
        STSIncludeSubdomains: true,
        STSPreload: true,
        FrameDeny: true,
        CustomFrameOptionsValue: "SAMEORIGIN",
        ContentTypeNosniff: true,
        BrowserXssFilter: true,
        ContentSecurityPolicy: "default-src 'self'",
    }).HandlerFuncWithNext))
    n.UseHandler(r)

    return http.ListenAndServe(listen, n)
}

