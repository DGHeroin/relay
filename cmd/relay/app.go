package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DGHeroin/relay"
)

var (
	local  = flag.String("l", ":53", "Listen Local Address")
	remote = flag.String("r", "1.1.1.1:53", "Forward To Remote Address")
	u      = flag.Bool("u", true, "Forward UDP")
	config = flag.String("c", "", "config file")
)

type Config struct {
	Configs []LinkConfig
}

type LinkConfig struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
	U   bool   `json:"udp"`
}

func OpenRelay(src, dst string, enableUDP bool) func() {
	tcp := relay.NewTCPRelay()
	stopCh := make(chan struct{})
	go tcp.Serve(src, dst, stopCh)

	if *u {
		udp := relay.NewUDPRelay()
		go udp.Serve(src, dst, stopCh)
	}

	return func() {
		close(stopCh)
	}
}

func ServeConfig(cfg *Config) {
	for _, c := range cfg.Configs {
		OpenRelay(c.Src, c.Dst, c.U)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *config != "" {
		if data, err := os.ReadFile(*config); err != nil {
			log.Println(err)
			return
		} else {
			cfg := &Config{}
			if err := json.Unmarshal(data, cfg); err != nil {
				log.Println(err)
				return
			} else {
				ServeConfig(cfg)
			}
		}
	} else {
		OpenRelay(*local, *remote, *u)
	}

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	msg := <-stopSig
	log.Println("exit with", msg)
}
