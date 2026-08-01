package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// verifica se a função Load carrega corretamente as variáveis de ambiente
func TestLoad_Success(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "api_pedidos")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword)
	assert.Equal(t, "api_pedidos", cfg.DBName)
}

// verifica se a função Load retorna erro quando uma variável de ambiente obrigatória está ausente
func TestLoad_MissingRequiredVar(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_HOST", "localhost")
	// DB_PORT, DB_USER, DB_PASSWORD, DB_NAME propositalmente ausentes

	_, err := Load()

	assert.Error(t, err)
}
