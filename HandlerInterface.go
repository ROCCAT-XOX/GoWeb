package main

import "net/http"

type SideCarHandler interface {
	Handle(s *SideCarServer) http.Handler
	SetNext(SideCarHandler) SideCarHandler
}

type SideCarServer struct {
}
