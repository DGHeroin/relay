package relay

import (
    "log"
    "net"
)

var (
    DefaultMTU = 1500
)

type udpRelay struct {
    source string
    target string
    remote map[string] *packet
}

func NewUDPRelay () Relay {
    r := udpRelay{}
    return &r
}

type packet struct {
    conn   net.Addr
    data   []byte
    pc     net.PacketConn
    remote *net.UDPConn
}

func (t *udpRelay) Serve(src, dst string) error {
    t.source = src
    t.target = dst
    t.remote = make(map[string]*packet)
    log.Printf("Serve UDP: %v => %v", src, dst)
    ln, err := net.ListenPacket("udp", src)
    if err != nil {
        log.Println(err)
        return err
    }

    buf := make([]byte, DefaultMTU)
    for {
        n, conn, err := ln.ReadFrom(buf)
        if err != nil {
            log.Println(err)
        }
        if pkt, ok := t.remote[conn.String()]; ok {
            pkt.data = buf[:n]
            go t.handlePacket(pkt)
        } else {
            pkt = &packet{conn: conn, data:buf[:n], pc:ln}
            t.remote[conn.String()] = pkt
            go t.handlePacket(pkt)
        }
    }
    return nil
}

func (t *udpRelay) handlePacket(p *packet){
    if p.remote == nil {
        rAddr, err := net.ResolveUDPAddr("udp", t.target)
        if err != nil {
            log.Println(err)
            t.die(p)
            return
        }
        remote, err := net.DialUDP("udp", nil, rAddr)
        if err != nil {
            log.Println(err)
            t.die(p)
            return
        }
        p.remote = remote
        go func() {
            buf := make([]byte, DefaultMTU)
            for {
                n, err := remote.Read(buf)
                if err != nil {
                    log.Println(err)
                    t.die(p)
                    return
                }
                data := buf[:n]
                p.pc.WriteTo(data, p.conn)
            }
        }()
    }
    _, err := p.remote.Write(p.data)
    if err != nil {
        t.die(p)
    }
}

func (t *udpRelay) die(p *packet) {
    log.Println("断开", p)
    go func() {
        delete(t.remote, p.conn.String())
    }()
}