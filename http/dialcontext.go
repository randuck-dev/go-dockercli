package http

type DialContext struct {
	Address string
	network string
}

func TcpDialContext(address string) DialContext {
	return DialContext{
		Address: address,
		network: "tcp",
	}
}

func UnixDialContext(socket string) DialContext {
	return DialContext{
		Address: socket,
		network: "unix",
	}
}
