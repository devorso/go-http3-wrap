package main

import (
	"flag"
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
	"os"
	"strings"
)

func checkFile(name string, list []os.DirEntry) bool {
	for _, v := range list {
		if strings.EqualFold(v.Name(), name) {
			return true
		}
	}
	return false
}

/**
Used for SPA apps.
*/

func main() {

	port := flag.String("port", "80", "The port")
	certkey := flag.String("certkey", "", "Path absolute to certificat key (tls)")
	privkey := flag.String("privatekey", "", "Path absolute to private key")
	flag.Parse()

	list := flag.NFlag()

	if list < 3 {
		log.Println("Please specify the arguments")
		return
	}
	log.Println("Go HTTP3 Wrapper by Devorso")
	log.Println("Check directory www...")
	mux := http.NewServeMux()
	listData, _ := os.ReadDir("./www")
	if len(listData) == 0 {
		log.Println("Wrapper closed. Nothing in www directoy.")
		return
	}
	log.Println("www directory contains files.. ok")
	fsStatic := http.FileServer(http.Dir("www/static"))
	fsRoot := http.FileServer(http.Dir("www"))

	// mux.Handle("/", http.FileServer(http.Dir("www")))
	mux.Handle("/static", fsStatic)

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		pathD := strings.Split(request.URL.Path, "/")

		if checkFile(pathD[1], listData) {
			fsRoot.ServeHTTP(writer, request)

		} else {
			if len(pathD) == 2 && pathD[1] == "" {
				fsRoot.ServeHTTP(writer, request)
			} else {
				// serve file index.html
				http.ServeFile(writer, request, "./www/index.html")
			}

		}
	})

	log.Println("HTTP3 Server started on ", *port)
	err := http3.ListenAndServe(`0.0.0.0:`+*port, *certkey, *privkey, mux)
	if err != nil {
		log.Println(err)
		return
	}

	select {}

}
