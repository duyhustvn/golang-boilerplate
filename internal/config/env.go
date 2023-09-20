package config

// Env structure
type Env struct {
	Environment string
	ServiceName string
}

// GetKeys gets crypto keys
func (a *Env) GetKeys() *Env {
	a.Environment = GetEnv("SERVICE_ENV")
	a.ServiceName = GetEnv("SERVICE_NAME")
	return a
}
