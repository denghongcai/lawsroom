package api

import(
    "encoding/json"
)

func ok(v interface{}) []byte {
    data := map[string]interface{}{
        "error": nil,
        "result": v,
    }
    r, _ := json.Marshal(data)
    return r
}

func err(code int, message string) []byte {
    data := map[string]interface{}{
        "error": map[string]interface{}{
            "code": code,
            "message": message,
        },
        "result": nil,
    }
    r, _ := json.Marshal(data)
    return r
}


