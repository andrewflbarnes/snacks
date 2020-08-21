package main

import (
    "fmt"
    "net"
    "bufio"
    "andrewflbarnes/snacks/loris"
    "strconv"
)

var (
    port = 8989
)

func main() {
    body := `{"ab":"cd"}`

    l := loris.New();
    lParams := loris.LorisVals{
        Endpoint: "/",
        ContentType: "application/json",
        Length: len(body),
        Body: body,
    }

    payload := l.Build(lParams)

    serverReady := make(chan bool)

    go server(port, serverReady)

    <-serverReady

    sendPayload("localhost", port, payload)
}

func sendPayload(host string, port int, payload string) {
    dest := host + ":" + strconv.Itoa(port)
    fmt.Println("CLIENT: Opening connection to " + dest)

    conn, err := net.Dial("tcp", dest)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    fmt.Println("CLIENT: Sending payload to " + dest)
    fmt.Fprintf(conn, payload)
    status, err := bufio.NewReader(conn).ReadString('\n');
    if  err != nil {
        panic(err)
    }

    fmt.Println("CLIENT: Received response:\n" + status)
}

func server(port int, ready chan bool) {
    strPort := strconv.Itoa(port)
    fmt.Println("SERVER: Opening server on: " + strPort)
    ln, err := net.Listen("tcp", "127.0.0.1:" + strPort)

    if err != nil {
        panic(err)
    }

    ready<-true

    for {
        conn, err := ln.Accept()
        fmt.Println("SERVER: Accepted connection on " + strPort)
        if err != nil {
            panic(err)
        } else {
            go handleConnection(conn)
        }
    }
}

func handleConnection(conn net.Conn) {
    lastBlank := false
    fullMsg := ""

    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            panic(err)
        }
        
        fullMsg += msg
        if msg == "\n" {
            if lastBlank {
                fmt.Println("SERVER: Received payload:\n" + fullMsg)
                fullMsg = ""
                lastBlank = false

                response := "Received\n"
                fmt.Print("SERVER: Sending response: " + response)
                fmt.Fprintf(conn, response)
            } else {
                lastBlank = true
            }
        }
    }
}
