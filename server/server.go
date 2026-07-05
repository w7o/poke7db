// Heavy assistance from ChatGPT 5.3 mini

// Not portable
// Used only for production purposes

package main

import (
    "fmt"
    "net/http"
)

func Start(addr string) error {
    // Register HTTP route; directs everything to handler
    http.HandleFunc("/", handler)

    fmt.Println("Server running on localhost:", addr)

    // Listen on this address and serve HTTP requests using default handler
    return http.ListenAndServe(":" + addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // w is what is sent back to the client
    // r contains all information about the incoming request
    // assumes all items in given folder is a json file
    http.ServeFile(w, r, "./localdata"+r.URL.Path+".json")
} 

func main() {
    Start("8080")
}