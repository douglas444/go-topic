package main

import (
    "fmt"
    "strings"
    "encoding/json"
    "net/http"
    "sync"
    "bytes"
    "log"
    "io/ioutil"
)

type Subscriber struct {
    Endpoint string `json:"endpoint"`
}

func notify(subscribers []Subscriber, body []byte) {
    for _, subscriber := range subscribers {
        res, err := http.Post(subscriber.Endpoint, "application/json", bytes.NewReader(body));
        if err != nil {
            log.Println(err);
        }
        res.Body.Close();
    }
}

func find(slice []Subscriber, value Subscriber) (int, bool) {
    for i, item := range slice {
        if item.Endpoint == value.Endpoint {
            return i, true
        }
    }
    return -1, false
}

func serverStart(subscribersByTopic map[string][]Subscriber) {

    var mutex sync.RWMutex;

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "no one's around to help");
    });

    http.HandleFunc("/write/", func(w http.ResponseWriter, r *http.Request) {

        if r.Body == nil {
            http.Error(w, "request body is missing", 400);
            return;
        }

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "invalid request body", 400);
            return;
        }

        topic := strings.TrimPrefix(r.URL.Path, "/write/")

        mutex.RLock();
        subscribers := subscribersByTopic[topic];
        mutex.RUnlock();

        go notify(subscribers, body);

    });

    http.HandleFunc("/subscribe/", func(w http.ResponseWriter, r *http.Request) {

        if r.Body == nil {
            http.Error(w, "request body is missing", 400);
            return;
        }

        var subscriber Subscriber;
        err := json.NewDecoder(r.Body).Decode(&subscriber)
        if err != nil {
            http.Error(w, "invalid request body", 400);
            return;
        }

        if (subscriber.Endpoint == "") {
            http.Error(w, "endpoint field cannot be empty", 400);
            return;
        }

        topic := strings.TrimPrefix(r.URL.Path, "/subscribe/");

        mutex.Lock();

        if _, contains := subscribersByTopic[topic]; !contains {
            subscribersByTopic[topic] = make([]Subscriber, 0);
        }

        if _, contains := find(subscribersByTopic[topic], subscriber); !contains {
            subscribersByTopic[topic] = append(subscribersByTopic[topic], subscriber);
        } else {
            http.Error(w, "already subscribed to this topic", 400);
        }
        mutex.Unlock();

    });

    http.ListenAndServe(":8080", nil);

}

func main() {

    topics := make(map[string][]Subscriber);
    go serverStart(topics);
    fmt.Println("running on 8080");

    c := make(chan int)
    <- c
}
