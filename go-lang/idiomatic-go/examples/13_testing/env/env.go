// Код, читающий переменную окружения. Тест на него — в env_test.go.
package env

import "os"

func ProcessEnvVars() Config {
	env, _ := os.LookupEnv("OUTPUT_FORMAT")
	// зависимость от окружения — в тесте ее подменяет t.Setenv
	return Config{
		OutputFormat: env,
	}
}

type Config struct {
	OutputFormat string
}
