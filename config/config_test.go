package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл конфигурации
	configContent := `
server:
  host: "localhost"
  port: "8080"
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "testpass"
  dbname: "testdb"
`
	// Используем strings.NewReader для избежания проблем с файлами на Windows
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader(configContent)); err != nil {
		t.Fatalf("Ошибка чтения конфигурации из строки: %v", err)
	}

	// Тестируем загрузку конфигурации
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Ожидался Server.Host 'localhost', получено '%s'", cfg.Server.Host)
	}
	if cfg.Database.DBName != "testdb" {
		t.Errorf("Ожидался Database.DBName 'testdb', получено '%s'", cfg.Database.DBName)
	}
}
