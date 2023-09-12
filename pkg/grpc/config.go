package grpc

type ServerConfig struct {
	AppID                string
	HostAddress          string
	Port                 int
	APIListenAddresses   []string
	NameSpace            string
	TrustDomain          string
	MaxRequestBodySizeMB int
	UnixDomainSocket     string
	ReadBufferSizeKB     int
	EnableAPILogging     bool
}
