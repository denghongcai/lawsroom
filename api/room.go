package api

import(
    "net/http"
    "strings"
    "github.com/satori/go.uuid"
)

func RoomPrepare(w http.ResponseWriter, r *http.Request){
    id := strings.Replace(uuid.NewV4().String(), "-", "", -1)
    w.Header().Set("Content-Type", "application/json")
    w.Write(ok(map[string]string{"name": id}))
}


