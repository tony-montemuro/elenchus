package main

type contextKey string

const isAuthenticatedKey = contextKey("isAuthenticated")
const requestIDKey = contextKey("requestID")
