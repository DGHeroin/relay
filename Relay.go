package relay

type Relay interface {
    Serve(src, dst string) error
}