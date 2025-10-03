package handlers

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"parentscontactform/internal/auth"
	"parentscontactform/internal/middleware"
	"parentscontactform/internal/models"
	"parentscontactform/internal/session"
	"parentscontactform/internal/util"
)

func (h *Handler) HandleLoginGet(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		if _, exists := session.Get(cookie.Value); exists {
			http.Redirect(w, r, "/form", http.StatusFound)
		}
	}

	loginHtml, err := h.staticFS.ReadFile("static/login.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(loginHtml)
}

func (h *Handler) HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		middleware.LogAndError(r, w, "Error parsing form", err.Error(), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	sessionId := util.GenerateRandomString(16)
	sessionData, err := auth.LoginUser(email, password)
	if err != nil {
		middleware.LogAndError(r, w, "Server Error", err.Error(), http.StatusInternalServerError)
		return
	}

	session.Set(sessionId, &sessionData)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   24 * 7 * 60 * 60,
		HttpOnly: true,
		Secure:   os.Getenv("PROD") != "",
	})

	state := util.GenerateRandomString(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "oidc_state",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		Secure:   os.Getenv("PROD") != "",
	})

	authURL := auth.Cfg.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("response_mode", "form_post"),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) HandleLogoutGet(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Delete(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   os.Getenv("PROD") != "",
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	sessionData, exists := session.Get(cookie.Value)

	if !exists {
		// Session not found (e.g., server was restarted)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Retrieve cookies
	stateCookie, err := r.Cookie("oidc_state")
	if err != nil {
		middleware.LogAndError(r, w, "state missing", err.Error(), http.StatusBadRequest)
		return
	}

	// Verify state
	stateReturned := r.FormValue("state")
	if stateReturned != stateCookie.Value {
		middleware.LogAndError(r, w, "invalid state", "invalid state", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oidc_state",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Secure:   os.Getenv("PROD") != "",
		MaxAge:   -1,
	})

	// Exchange code for tokens
	code := r.FormValue("code")
	token, err := auth.Cfg.Exchange(context.Background(), code)
	if err != nil {
		middleware.LogAndError(r, w, "token exchange failed", err.Error(), http.StatusInternalServerError)
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		middleware.LogAndError(r, w, "failed to retrieve id_token", "malformed token", http.StatusInternalServerError)
		return
	}

	parsedToken, err := auth.Verifier.Verify(ctx, rawIdToken)
	if err != nil {
		middleware.LogAndError(r, w, "failed to verify id_token", err.Error(), http.StatusInternalServerError)
		return
	}

	var claims models.ClaimExtract
	if err = parsedToken.Claims(&claims); err != nil {
		middleware.LogAndError(r, w, "failed to parse claims", err.Error(), http.StatusInternalServerError)
		return
	}

	isamsTokenData := models.ISAMSToken{
		UserClaims: &claims,
		Token:      token,
	}

	sessionData.ISAMSToken = isamsTokenData

	session.Set(cookie.Value, sessionData)

	http.Redirect(w, r, "/form", http.StatusFound)
}
