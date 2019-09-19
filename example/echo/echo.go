package main

import (
	"log"
	"net/http"

	"github.com/xlionet/monkey"
)

type Echo struct {
	monkey.WSProtocol
}

// OnTransportData ...
func (p *Echo) OnTransportData(transport monkey.Transport, env *monkey.Envelope) {
	if env == nil {
		log.Fatal("nil msg")
	}

	transport.SendData(env)

}
func main() {

	protocol := &Echo{}
	cfg := monkey.NewConfig()
	cfg.WSListernPort = 8080
	mk := monkey.New(protocol, cfg)

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		err := mk.HandleConnection(w, r)
		if err != nil {
			log.Fatal(err)
		}
	})
	http.HandleFunc("/", index)

	mk.Serve(http.DefaultServeMux)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "echo.html")
}
