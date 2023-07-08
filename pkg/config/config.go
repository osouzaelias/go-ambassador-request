package config

import "os"

type Config struct {
	Hostname          string
	MaxNumberMessages string
	Region            string
	QueueUrl          string
	Table             string
}

func NewConfig() *Config {
	return &Config{
		Hostname:          getEnv("HOSTNAME", "localhost:8888"),
		MaxNumberMessages: getEnv("MAX_NUMBER_MESSAGES", "10"),
		Region:            getEnv("REGION", "us-west-2"),
		QueueUrl:          getEnv("QUEUE_URL", "https://sqs.us-west-2.amazonaws.com/000/xx-happy-path"),
		Table:             getEnv("TABLE", "xx-happy-path"),
	}

	/*file, err := os.ReadFile(fmt.Sprintf("%s.json", environment))
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file, %s", err)
	}

	return &config*/
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
