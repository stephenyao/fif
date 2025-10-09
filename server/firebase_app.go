package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func initFirebaseApp() (*firebase.App, error) {
	// Load .env for local/dev
	_ = godotenv.Load()

	b64 := os.Getenv("FIREBASE_KEY_B64")
	if b64 == "" {
		return nil, fmt.Errorf("FIREBASE_KEY_B64 is not set")
	}
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode FIREBASE_KEY_B64: %w", err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(decoded))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase App: %w", err)
	}
	return app, nil
}
