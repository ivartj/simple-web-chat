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
	"io"
)

const (
	mainProgramName = "simple-web-chat"
	mainProgramVersion = "0.1-SNAPSHOT"
)

func mainHelpMessage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [ -p <port-number> ] [ --assets <assets-directory> ]\n", mainProgramName)
}

func main() {
	tok := args.NewTokenizer(os.Args)
	var port uint64 = 80
	assets := "./assets"
	for tok.Next() {
		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				mainHelpMessage(os.Stdout)
				os.Exit(0)

			case "--version":
				fmt.Printf("%s version %s\n", mainProgramName, mainProgramVersion)
				os.Exit(0)

			case "-p", "--port":
				portstr, err := tok.TakeParameter()
				if err != nil {
					log.Fatalf("Error on parsing command-line arguments: %s\n", err.Error())
				}
				port, err = strconv.ParseUint(portstr, 10, 16)
				if err != nil {
					log.Fatalf("Failed to parse port number given on %s: %s\n", tok.Arg(), err.Error())
				}
			case "--assets":
				var err error
				assets, err = tok.TakeParameter()
				if err != nil {
					log.Fatalf("Error on parsing command-line arguments: %s\n", err.Error())
				}
			default:
				log.Fatalf("Unrecognized option, '%s'\n", tok.Arg())
			}

		} else {
			log.Fatalf("Unexpected non-option argument '%s'\n", tok.Arg())
		}
	}
	if tok.Err() != nil {
		log.Fatalf("Error on parsing command-line arguments: %s\n", tok.Err().Error())
	}
	rand.Seed(time.Now().Unix())
	ctx := newContext(assets)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), ctx))
}

