package main

import (
	"fmt"
	"log"
	"net" // standard network package
	"strings"
)

func main() {
	// config
	port := 8000
	protocol := "tcp"

	// resolve TCP address
	addr, err := net.ResolveTCPAddr(protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln(err)
	}

	// get TCP socket
	socket, err := net.ListenTCP(protocol, addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listen: ", socket.Addr().String())

	// keep listening
	for {
		// wait for connection
		conn, err := socket.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Connected by ", conn.RemoteAddr().String())

		// yield connection to concurrent process
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// close connection when this function ends
	defer conn.Close()

	// read request
	buf := make([]byte, 1024)
	conn.Read(buf)
	met, path := getRequestMethodAndPath(string(buf))
	fmt.Printf("Method:%s\nPath:%s\n", met, path)

	// write response
	conn.Write([]byte(makeResponse(string(buf))))
}

func getRequestMethodAndPath(s string) (string, string){
	var met, path string

	line := strings.Split(s, "\n")[0]
	met = strings.Split(line, " ")[0]
	path = strings.Split(line, " ")[1]

	return met, path
}

func makeResponse(s string) (string) {
	lines := strings.Split(s, "\n")
	content := "Hello World!"

	// line 1
	//res := strings.Split(lines[0], " ")[2]
	res := "HTTP/1.1"
	res = fmt.Sprintf("%s 200 OK\n", res)

	// line 2-8
	res += "Server: Apache\nSet-Cookie: csrftoken=2d1u2#sk2oi*\n"
	res += lines[6] + "\n" + lines[7] + "\n" + lines[8] + "\n"
	res += "Content-Type: text/html; charset=utf-8\n"
	res += fmt.Sprintf("Content-Length: %d\n", len(content))

	// line 9
	res += "\n"

	// line 10
	res += content

	fmt.Println(res)

	return res
}
