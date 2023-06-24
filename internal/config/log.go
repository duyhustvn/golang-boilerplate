package config

// Logger struct
type Logger struct {
	Level string
	Path  string
}

// GetLoggerEnv get keys
func (l *Logger) GetLoggerEnv() *Logger {
	l.Level = GetEnv("LOG_LEVEL")
	l.Path = GetEnv("LOG_PATH")
	return l
}
