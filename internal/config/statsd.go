package config

type Statds struct {
	Addr string
}

func (s *Statds) GetStatdsEnv() *Statds {
	s.Addr = GetEnv("STATDS_ADDR")
	return s
}
