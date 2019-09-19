# monkey

a high flexible websocket framwork base on gorilla

## what is about monkey

    monkey is a really simple websoket framwork base on gorilla, which makes you build your app in a simple way
    - high flexible API 
    - as little as possible dependency
    - easy to understand

### main API

    - OnTransportMade
    - OnTransportLost
    - OnTransportData
    - OnPing
    - OnPong
  
## how to use

generally, you just need to implement `Protocol` interface with your business, and then new and run a monkey instance, then it works

    1. monkey.New() got a monkey instance with your own protocol
    2. new a multiplexes with websocket path and handlers, 
    3. register protocol with `HandleConnection` in websocket handles
    4. start server with your multiplexes(usually is `http.DefaultServeMux`)

## get

   `go get -u github.com/xlionet/monkey`

## example

    ``` Golang

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
    ```
