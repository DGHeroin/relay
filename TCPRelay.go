package relay

import (
    "io"
    "log"
    "net"
    "time"
)

type tcpRelay struct {
    source string
    target string

    listen net.Listener
}

func NewTCPRelay() Relay {
    r := tcpRelay{}
    return &r
}

func (t *tcpRelay) Serve(src, dst string) error {
    t.source = src
    t.target = dst
    log.Printf("Serve TCP: %v => %v", src, dst)
    if ln, err := net.Listen("tcp", src); err != nil {
        return err
    } else {
        t.listen = ln
    }

    for {
        if conn, err := t.listen.Accept(); err != nil {
            return err
        } else {
            go t.handleConnection(conn)
        }
    }
}

func (t *tcpRelay) handleConnection(conn net.Conn) error {
    var (
        rAddr  *net.TCPAddr
        remote net.Conn
        err    error
    )
    rAddr, err = net.ResolveTCPAddr("tcp", t.target)
    if err != nil {
        log.Println(err)
        return err
    }
    if rAddr.IP.To4() != nil {
        remote, err = net.DialTimeout("tcp", rAddr.String(), time.Second*5)
    } else {
        remote, err = net.DialTimeout("tcp6", rAddr.String(), time.Second*5)
    }
    if err != nil {
        log.Println(err)
        return err
    }
    defer conn.Close()
    defer remote.Close()

    go io.Copy(conn, remote)
    io.Copy(remote, conn)
    return nil
}
