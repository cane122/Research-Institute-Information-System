package config

type Config struct {
	UploadPath string
}

func LoadConfig() Config {
	// You can expand this to read from a file or environment variables
	return Config{
		UploadPath: "./uploads", // Default upload path
	}
}
