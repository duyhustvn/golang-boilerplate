package adapter

// Labels provides values for labels
type Labels map[string]string

// CollectorOptions stores details of collector
type CollectorOptions struct {
	Name   string
	Help   string
	Labels []string
}

// Counter goes up
type Counter interface {
	Inc(Labels)
	Add(int64, Labels)
}

// Gauge can be set to anything
type Gauge interface {
	Add(int64, Labels)
	Sub(int64, Labels)
	Inc(Labels)
	Set(int64, Labels)
}

// Timer observes trends
type Timer interface {
	Observe(int64, Labels)
}
