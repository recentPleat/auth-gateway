package authgateway

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// generateRandomString generates a cryptographically secure random string
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashPassword hashes a password using PBKDF2
func hashPassword(password string, salt string) (string, error) {
	hash := pbkdf2.Key([]byte(password), []byte(salt), 10000, 32, sha256.New)
	return base64.URLEncoding.EncodeToString(hash), nil
}

// verifyPassword verifies a password by comparing it to a hashed password
func verifyPassword(password string, hashedPassword string, salt string) (bool, error) {
	hash, err := hashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return hash == hashedPassword, nil
}

// getRemoteAddress gets the remote address from an HTTP request
func getRemoteAddress(r *http.Request) string {
	// try to get the real IP address from the 'X-Forwarded-For' header
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// if the header is present, return the first IP address
		return strings.Split(xForwardedFor, ",")[0]
	}
	// if the header is not present, return the remote address from the request
	return r.RemoteAddr
}

// logError logs an error with the current time and a message
func logError(message string, err error) {
	log.Printf("%v: %s: %v", time.Now(), message, err)
}

// setErrorResponse sets an error response for an HTTP request
func setErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}