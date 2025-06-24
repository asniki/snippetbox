package main

// contextKey is a custom type for the context keys
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
