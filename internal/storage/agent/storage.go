package storage

type (
	IAgentStorage interface {
		SaveMetric(string, float64)
		LoadMetric(string) (float64, error)
		UpdateCounter() error
		LoadCounter() (int64, error)
	}
)
