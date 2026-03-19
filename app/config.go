package app

type Config struct {
	DataDir string
}

func DefaultConfig() Config {
	return Config{DataDir: ".data/storage"}
}
