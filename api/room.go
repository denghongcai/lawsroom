package api

import(
    "net/http"
    "strings"
    "encoding/json"
    "github.com/satori/go.uuid"
    "io/ioutil"
    "git.txthinking.com/txthinking/signal"
)

type PrepareOutput struct {
    Name string
}
func RoomPrepare(w http.ResponseWriter, r *http.Request){
    id := strings.Replace(uuid.NewV4().String(), "-", "", -1)
    w.Header().Set("Content-Type", "application/json")
    output := PrepareOutput{
        Name: id,
    }
    w.Write(_ok(output))
}

type CapacityInput struct {
    Name string
}
type CapacityOutput struct {
     Capacity int
     Used int
     Idle int
}
func RoomCapacity(w http.ResponseWriter, r *http.Request){
    input := CapacityInput{}
    defer r.Body.Close()
    var err error
    var body []byte
    if body, err = ioutil.ReadAll(r.Body); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    if err = json.Unmarshal(body, &input); err != nil {
        w.Write(_error(100, err.Error()))
        return
    }
    c := signal.ROOM_CAPACITY
    u := 0
    if _, ok := signal.Rooms[input.Name]; ok {
        u = len(signal.Rooms[input.Name].Peers)
    }
    i := c - u
    output := CapacityOutput{
        Capacity: c,
        Used: u,
        Idle: i,
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(_ok(output))
}


