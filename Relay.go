package relay

type Relay interface {
	Serve(src, dst string, stopCh chan struct{}) error
}
