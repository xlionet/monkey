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


