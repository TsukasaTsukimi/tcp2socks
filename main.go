package main

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/net/proxy"
)

func forward(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}

func process(client net.Conn, dialer proxy.Dialer) {
	target, err := dialer.Dial("tcp", "198.251.81.97:80")
	if err != nil {
		fmt.Printf("Connect failed: %v\n", err)
		return
	}
	go forward(client, target)
	go forward(target, client)
}

func main() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:2801", nil, proxy.Direct)
	if err != nil {
		fmt.Printf("Socks5 failed: %v\n", err)
	}
	server, err := net.Listen("tcp", ":2805")
	if err != nil {
		fmt.Printf("Listen failed: %v\n", err)
		return
	}

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Printf("Accept failed: %v\n", err)
			continue
		}
		go process(client, dialer)
	}
}
