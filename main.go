package main

import (
	"net/http"
	"log"
	"math/rand"
	"time"
	"github.com/ivartj/minn/args"
	"strconv"
	"os"
	"fmt"
)

func main() {
	tok := args.NewTokenizer(os.Args)
	var port uint64 = 80
	assets := "./assets"
	for tok.Next() {
		if tok.IsOption() {

			switch tok.Arg() {
			case "-p", "--port":
				portstr, err := tok.TakeParameter()
				if err != nil {
					log.Fatalf("Error on parsing command-line arguments: %s\n", err.Error())
				}
				port, err = strconv.ParseUint(portstr, 10, 16)
				if err != nil {
					log.Fatalf("Failed to parse port number given on %s: %s", tok.Arg(), err.Error())
				}
			case "--assets":
				var err error
				assets, err = tok.TakeParameter()
				if err != nil {
					log.Fatalf("Error on parsing command-line arguments: %s\n", err.Error())
				}
			default:
				log.Fatalf("Unrecognized option, '%s'", tok.Arg())
			}

		} else {
			log.Fatalf("Unexpected non-option argument '%s'", tok.Arg())
		}
	}
	if tok.Err() != nil {
		log.Fatalf("Error on parsing command-line arguments: %s\n", tok.Err().Error())
	}
	rand.Seed(time.Now().Unix())
	ctx := newContext(assets)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), ctx))
}

