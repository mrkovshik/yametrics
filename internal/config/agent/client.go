package config

type ClientConfig struct {
	Key            string
	Address        string
	ReportInterval int
	RateLimit      int
	CryptoKey      string
}
