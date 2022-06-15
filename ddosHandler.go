package main

import (
	"fmt"
	"net/http"
)

var ddosMap = make(map[string]string)

type ddosHandler struct {
	next SideCarHandler
}

// validation check for satisfying interface
var _ SideCarHandler = &ddosHandler{}

func (h *ddosHandler) SetNext(next SideCarHandler) SideCarHandler {
	h.next = next
	return next
}

func (h *ddosHandler) Handle(s *SideCarServer) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		name := request.FormValue("name")
		fmt.Println(ddosMap[name])
		if ddosMap[name] == "" {
			h.next.Handle(s).ServeHTTP(rw, request)

		} else {
			Response("verpiss dich", 401, rw)
		}

	})

}
