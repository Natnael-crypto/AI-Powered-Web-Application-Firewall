package captcha

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	captchaPositions = make(map[string]int)
	captchaTokens    = make(map[string]time.Time)
	mu               sync.Mutex
)

// GenerateRandomPosition returns a random number (slider position)
func GenerateRandomPosition() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(200) + 50 // Slider range: 50 to 250
}

// GenerateToken creates a secure token and stores it
func GenerateToken() string {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	mu.Lock()
	captchaTokens[token] = time.Now().Add(1 * time.Hour)
	mu.Unlock()

	return token
}

// ValidateCaptchaToken checks if a token is valid
func ValidateCaptchaToken(token string) bool {
	mu.Lock()
	defer mu.Unlock()

	expiry, exists := captchaTokens[token]
	if !exists || time.Now().After(expiry) {
		delete(captchaTokens, token) // Remove expired tokens
		return false
	}
	return true
}

// ChallengePage serves the CAPTCHA challenge page
func ChallengePage(w http.ResponseWriter, r *http.Request) {
	position := GenerateRandomPosition()

	// Store the expected position
	sessionID := strconv.FormatInt(time.Now().UnixNano(), 10) // Unique session ID
	mu.Lock()
	captchaPositions[sessionID] = position
	mu.Unlock()

	// Send HTML
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Slider CAPTCHA</title>
			<style>
				body { text-align: center; font-family: Arial, sans-serif; }
				#slider-container { width: 300px; margin: 50px auto; position: relative; }
				#slider { width: 100%%; }
				#box { position: absolute; top: 10px; width: 30px; height: 30px; background: #3498db; }
			</style>
			<script>
				let correctPos = %d;
				function validate() {
					let userPos = document.getElementById('slider').value;
					fetch('/validate-captcha', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ position: parseInt(userPos), session: '%s' })
					}).then(res => res.json()).then(data => {
						if (data.success) {
							alert('CAPTCHA passed!');
							window.location.href = "/";
						} else {
							alert('Try again!');
							window.location.reload();
						}
					});
				}
			</script>
		</head>
		<body>
			<h2>Move the slider to the correct position</h2>
			<div id="slider-container">
				<input type="range" id="slider" min="50" max="250">
				<div id="box"></div>
			</div>
			<button onclick="validate()">Submit</button>
		</body>
		</html>
	`, position, sessionID)
}

// ValidateCaptcha checks the user's slider position
func ValidateCaptcha(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Position int    `json:"position"`
		Session  string `json:"session"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	expectedPos, exists := captchaPositions[data.Session]
	delete(captchaPositions, data.Session) // Remove session after validation
	mu.Unlock()

	if !exists || abs(expectedPos-data.Position) > 10 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	// Generate and set CAPTCHA token
	token := GenerateToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "captcha_token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Middleware to enforce CAPTCHA verification
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("captcha_token")
		if err != nil || !ValidateCaptchaToken(cookie.Value) {
			http.Redirect(w, r, "/captcha", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Helper function to get absolute difference
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
