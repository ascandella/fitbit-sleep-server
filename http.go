package main

import "net/http"

type myHandler struct{}

func registerServeMux(handler myHandler) {
	http.Handle("/foo", handler)
}
