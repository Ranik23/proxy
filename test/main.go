package main

import (
	"flag"
	"log"

	"github.com/Ranik23/proxy"
)



func main() {

	var pemPath string
    flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
    var keyPath string
    flag.StringVar(&keyPath, "key", "server.key", "path to key file")
    var proto string
    flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
    var address string
    flag.StringVar(&address, "address", "localhost:9091", "address on which the serevr will listen")
    flag.Parse()
    if proto != "http" && proto != "https" {
        log.Fatal("Protocol must be either http or https")
    }

	p := proxy.NewProxy(address, pemPath, keyPath, proto)

    p.OnRequest(proxy.Is("http://github.com/Ranik23")).Do(proxy.AlwaysReject)
    
	if err := p.Run(); err != nil {
        log.Fatal(err)
    }
}