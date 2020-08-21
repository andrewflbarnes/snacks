package main

import (
    "fmt"
    "net"
    "bufio"
)

func main() {
    ln, err := net.Listen("tcp", ":8111")

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err.Error())
        } else {
            go handleConnection(conn)
        }
    }
}

func handleConnection(conn net.Conn) {
    for {
        msg, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println(err.Error())
            fmt.Printf("Closing connection: %q\n", conn)
            conn.Close()
            return
        } else {
            fmt.Println("PING")
            fmt.Println(msg)
            fmt.Fprintf(conn, msg)
            fmt.Println("PONG")
        }
    }
}
