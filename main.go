package main

import (
	"net/http"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	ctx := newContext()
	log.Fatal(http.ListenAndServe(":9999", ctx))
}

