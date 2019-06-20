package main

import (
    "flag"
    "io"
    "log"
    "net"
)

func echoTCP(addr string) {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        panic(err)
    }
    defer ln.Close()
    for {
        conn, err := ln.Accept()
        if err != nil {
            panic(err)
        }
        
        go func(conn net.Conn) {
            defer conn.Close()
            io.Copy(conn, conn)
        }(conn)
    }
}

func echoUDP(addr string) {
    pc, err := net.ListenPacket("udp", addr)
    if err != nil {
        panic(err)
    }
    defer pc.Close()
    buf := make([]byte, 1500)
    for {
        n, conn, err := pc.ReadFrom(buf)
        if err != nil {
            panic(err)
        }
        data := buf[:n]
        log.Println(string(data))
        pc.WriteTo(data, conn)
    }
}

func main() {
    addr := flag.String("a", ":7771", "listen addr")
    flag.Parse()
    log.Println("listen on:", *addr)
    go echoTCP(*addr)
    go echoUDP(*addr)
    select {}
}
