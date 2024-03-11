package storage

type (
	IAgentStorage interface {
		SaveMetric(string, string)
		LoadMetric(string) string
		UpdateCounter() error
	}
)
