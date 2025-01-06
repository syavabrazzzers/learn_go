package settings

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type PostgresSettings struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_NAME,required"`
}

func (s *PostgresSettings) PostgresDsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", s.Host, s.Port, s.User, s.Password, s.Name)
}

type RedisSettings struct {
	Host string `env:"REDIS_HOST,required"`
	Port string `env:"REDIS_PORT,required"`
	DB   int    `env:"REDIS_DB,required"`
}

func (r *RedisSettings) Url() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type SmtpSettings struct {
	Host     string `env:"SMTP_HOST"`
	Port     string `env:"SMTP_PORT"`
	From     string `env:"SMTP_FROM"`
	Password string `env:"SMTP_PASSWORD"`
}

type SettingsStruct struct {
	Postgres                   PostgresSettings
	Redis                      RedisSettings
	Smtp                       SmtpSettings
	JwtSecret                  string `env:"JWT_SECRET,required"`
	Debug                      bool   `env:"DEBUG,required"`
	VerificationCodeExpiration int    `env:"VERIFICATION_CODE_EXP,required"`
}

var Settings SettingsStruct

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	err = env.Parse(&Settings)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}
