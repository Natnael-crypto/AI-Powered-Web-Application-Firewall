package config

import "os"

// JWTSecretKey is used to sign JWT tokens
var JWTSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
