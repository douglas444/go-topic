package main

import (
    "fmt"
    "net/http"
)


func serverStart(exit chan bool) {
    
    http.HandleFunc("/exit", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("exiting");
        w.WriteHeader(202);
        exit <- true;
    });

    http.ListenAndServe(":8080", nil);

}

func main() {
    exit := make(chan bool);
    go serverStart(exit);
    fmt.Println("running on 8080");
    <- exit;
}
