package main

import "net/http"

func addDefaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' cdnjs.cloudflare.com; fonts.googleapis.com; font-src fonts.gstatic.com;")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func addSessionManager(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		app.SessionManager.LoadAndSave(next)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
