package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	// ENV  = os.Getenv("ENVIRONMENT")
	PORT = os.Getenv("PORT")

	JWT_SECRET_PRIVATE_KEY = os.Getenv("JWT_SECRET_PRIVATE_KEY")
	JWT_SECRET_PUBLIC_KEY  = os.Getenv("JWT_SECRET_PUBLIC_KEY")
)
