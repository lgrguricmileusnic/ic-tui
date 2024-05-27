package args

import "flag"

type Args struct {
	Addr string
	Ssh  bool
}

func ParseArgs() Args {
	addr := flag.String("a", ":8080", "Address and port for HTTP API in format <ip_address>:<port>")
	ssh := flag.Bool("ssh", false, "Set to true to run as a Wish SSH server")
	flag.Parse()
	return Args{Addr: *addr, Ssh: *ssh}
}
