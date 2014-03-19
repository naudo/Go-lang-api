// https://github.com/fiorix/go-redis
package main

import(
    "fmt"
    "net/http"
    "github.com/fiorix/go-redis/redis"
    "encoding/json"
    "log"
)

func main() {
    fmt.Println("binding to / on port 4000")
    http.HandleFunc("/", hello)

    log.Fatal(http.ListenAndServe(":4000", nil))
} 


type user struct {
    Username string
    Password string
}

func hello(res http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()
    decoder := json.NewDecoder(req.Body)
    var user_struct user 

    err := decoder.Decode(&user_struct)
    if err != nil {
        panic("I can't parse this stuff")
    }

    js, _ := json.Marshal(user_struct)
    go save_to_redis("new_key", string(js))

    wrap_response(res, js)
}

func save_to_redis(key string, value string){
    redis_client := redis.New("127.0.0.1:6379")
    redis_client.Set(key, value)
}

func wrap_response(res http.ResponseWriter, body []byte) {
    // res.Header().Set("Content-Type", "application/json")
    res.Header().Set("Connection", "close")
    res.WriteHeader(http.StatusAccepted)
    res.Write(body)
}