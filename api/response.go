package api

import(
    "encoding/json"
)

func _ok(v interface{}) []byte {
    data := map[string]interface{}{
        "Error": nil,
        "Result": v,
    }
    r, _ := json.Marshal(data)
    return r
}

func _error(code int, message string) []byte {
    data := map[string]interface{}{
        "Error": map[string]interface{}{
            "Code": code,
            "Message": message,
        },
        "Result": nil,
    }
    r, _ := json.Marshal(data)
    return r
}


