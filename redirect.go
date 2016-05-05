package main

import(
    "net/http"
    "os"
    "io/ioutil"
    "bytes"
)

func redirect(w http.ResponseWriter, r *http.Request){
    f, err := os.Open("./public/index.html")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    html, err := ioutil.ReadAll(f)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    js := "<script>page('" + r.URL.Path + "');</script>"
    html = bytes.Replace(html, []byte("<!-- redirect -->"), []byte(js), 1)
    w.Write(html)
}

