package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubUser struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handlers) getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     h.config.GitHubClientID,
		ClientSecret: h.config.GitHubSecret,
		RedirectURL:  h.config.BaseURL + "/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func (h *Handlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title      string
		Breadcrumb string
	}{
		Title:      "Login",
		Breadcrumb: "Login",
	}

	h.renderTemplate(w, "login.html", data)
}

func (h *Handlers) handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Generate state token
	state := generateStateToken()

	// Store state in session/cookie for verification
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   600, // 10 minutes
	})

	config := h.getOAuthConfig()
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handlers) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	config := h.getOAuthConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from GitHub
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Check if user is allowed
	if !h.isUserAllowed(user.Login) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    user.Login,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   86400 * 7, // 7 days
	})

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handlers) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) isUserAllowed(username string) bool {
	if len(h.config.AllowedUsers) == 0 {
		return true // Allow all if no whitelist is configured
	}

	for _, allowedUser := range h.config.AllowedUsers {
		if allowedUser == username {
			return true
		}
	}
	return false
}

func (h *Handlers) getSessionUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func generateStateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
