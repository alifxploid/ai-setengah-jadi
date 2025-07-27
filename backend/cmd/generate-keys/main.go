package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

func main() {
	fmt.Println("🔐 Generating secure keys...")

	// Generate JWT Secret (64 bytes = 512 bits)
	jwtSecret, err := generateSecureKey(64)
	if err != nil {
		log.Fatal("Failed to generate JWT secret:", err)
	}

	// Generate Access Key (32 bytes = 256 bits)
	accessKey, err := generateSecureKey(32)
	if err != nil {
		log.Fatal("Failed to generate access key:", err)
	}

	// Update .env file
	envPath := "../../.env"
	err = updateEnvFile(envPath, jwtSecret, accessKey)
	if err != nil {
		log.Fatal("Failed to update .env file:", err)
	}

	fmt.Println("✅ Keys generated and .env updated successfully!")
	fmt.Printf("🔑 JWT_SECRET: %s\n", jwtSecret)
	fmt.Printf("🗝️  ACCESS_KEY: %s\n", accessKey)
	fmt.Println("\n⚠️  Keep these keys secure and never commit them to version control!")
}

// generateSecureKey generates a cryptographically secure random key
func generateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Use base64 URL encoding for safe usage in environment variables
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// updateEnvFile updates the .env file with new keys
func updateEnvFile(envPath, jwtSecret, accessKey string) error {
	// Read existing .env file
	content, err := ioutil.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	envContent := string(content)

	// Update JWT_SECRET
	jwtRegex := regexp.MustCompile(`JWT_SECRET=.*`)
	if jwtRegex.MatchString(envContent) {
		envContent = jwtRegex.ReplaceAllString(envContent, fmt.Sprintf("JWT_SECRET=%s", jwtSecret))
	} else {
		envContent += fmt.Sprintf("\nJWT_SECRET=%s", jwtSecret)
	}

	// Update ACCESS_KEY
	accessRegex := regexp.MustCompile(`ACCESS_KEY=.*`)
	if accessRegex.MatchString(envContent) {
		envContent = accessRegex.ReplaceAllString(envContent, fmt.Sprintf("ACCESS_KEY=%s", accessKey))
	} else {
		envContent += fmt.Sprintf("\nACCESS_KEY=%s", accessKey)
	}

	// Add timestamp comment
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	comment := fmt.Sprintf("\n# Keys generated on %s", timestamp)
	if !strings.Contains(envContent, "# Keys generated on") {
		envContent += comment
	} else {
		// Update existing timestamp
		timestampRegex := regexp.MustCompile(`# Keys generated on .*`)
		envContent = timestampRegex.ReplaceAllString(envContent, strings.TrimPrefix(comment, "\n"))
	}

	// Write back to file
	err = ioutil.WriteFile(envPath, []byte(envContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}