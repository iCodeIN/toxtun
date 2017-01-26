package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kitech/colog"
)

const (
	tunconv = uint32(0xaabbccdd)
	tunmtu  = 1000
	rdbufsz = 8192
)

var (
	// options
	inst_mode   string // = "server" | "client"
	kcp_mode    string // = "default" // fast
	config_file string // = "toxtun_whtun.ini"
	config      *TunnelConfig
	log_level   int = int(colog.LDebug)
)

const (
	ltracep   = "trace: "
	ldebugp   = "debug: "
	linfop    = "info: "
	lwarningp = "warning: "
	lerrorp   = "error: "
	lalertp   = "alert: "
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	colog.Register()
	/*
		log.Println("debug: ")
		log.Println("info: ", 1)
		log.Println("warning: ", 1)
		log.Println("error: ", 1)
	*/

	flag.StringVar(&kcp_mode, "kcp-mode", "fast", "default|fast")
	if !(kcp_mode == "default" || kcp_mode == "fast") {
		kcp_mode = "fast"
	}

	flag.StringVar(&config_file, "config", "", "config file .ini")
	flag.IntVar(&log_level, "log-level", int(colog.LDebug),
		fmt.Sprintf("%d - %d, default %d\n", colog.LTrace, colog.LAlert, colog.LDebug))
}

func main() {
	flag.Parse()
	if len(config_file) > 0 {
		config = NewTunnelConfig(config_file)
		log.Println(linfop, config)
	}

	argv := flag.Args()
	argc := len(argv)

	if argc > 0 {
		inst_mode = argv[argc-1]
	}

	go NewStatServer().serve()
	appevt.Trigger("appmode", inst_mode)

	switch inst_mode {
	case "client":
		if config == nil {
			log.Println("error:", "need config file")
			flag.PrintDefaults()
			os.Exit(-1)
		}
		tc := NewTunnelc()
		tc.serve()
	case "server":
		td := NewTunneld()
		td.serve()
	default:
		log.Println("Invalid mode:", inst_mode, ", server/client.")
		flag.PrintDefaults()
	}

}
