package config

// Env structure
type Env struct {
	Environment string
}

// GetKeys gets crypto keys
func (a *Env) GetKeys() *Env {
	a.Environment = GetEnv("APP_ENV")
	return a
}
