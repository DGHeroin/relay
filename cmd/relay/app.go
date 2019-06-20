package main

import (
    "flag"
    "github.com/DGHeroin/relay"
    "log"
    "os"
    "os/signal"
    "syscall"
)

var (
    local  = flag.String("l", ":53", "Listen Local Address")
    remote = flag.String("r", "1.1.1.1:53", "Forward To Remote Address")
    u      = flag.Bool("u", true, "Forward UDP")
)

func main() {
    flag.Parse()
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.Printf("%v => %v \n", *local, *remote)
    tcp := relay.NewTCPRelay()
    go tcp.Serve(*local, *remote)

    if *u {
        udp := relay.NewUDPRelay()
        go udp.Serve(*local, *remote)
    }

    stopSig := make(chan os.Signal)
    signal.Notify(stopSig, os.Interrupt, os.Kill, syscall.SIGINT)
    select {
    case msg := <-stopSig:
        log.Println("exit with", msg)
        break
    }
}
