package main

import(
    "net/http"
    "github.com/txthinking/signal"
)

func getSignalHandle() *signal.Signal{
    signal.ROOM_CAPACITY = 5
    return signal.New(func(r *http.Request) bool {
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
}

