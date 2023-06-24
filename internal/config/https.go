package config

// HTTPS
type HTTPS struct {
	Cert string
	Key  string
	Port string
}

// GetHTTPS returns https settings
func (h *HTTPS) GetHTTPSEnv() *HTTPS {
	h.Cert = GetEnv("HTTPS_CERT")
	h.Key = GetEnv("HTTPS_KEY")
	h.Port = GetEnv("HTTPS_PORT")

	return h
}
