package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"net/http"
)

// cookie handling
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler
type login_handler struct {
	next SideCarHandler
}

func (h *login_handler) SetNext(next SideCarHandler) SideCarHandler {
	h.next = next
	return next
}

// validation check for satisfying interface
var _ SideCarHandler = &login_handler{}

func (h *login_handler) Handle(s *SideCarServer) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {

		name := request.FormValue("name")
		pass := request.FormValue("password")
		//redirectTarget := "/"
		if name != "" && pass != "" {
			if name == "test123" {
				ddosMap[name] = "blocked"
				//zeit einfügen als Zeitstempel
				fmt.Println(ddosMap)
				rw.Write([]byte("login failed"))
				return
			}
			// .. check credentials ..
			setSession(name, rw)
			//redirectTarget = "/internal"

		}
		tmpl.ExecuteTemplate(rw, "login_new.html", nil)

		//http.Redirect(rw, request, redirectTarget, 302)

	})

}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// internal page
/*
const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Header)
	fmt.Println(request.RemoteAddr)
	userName := getUserName(request)
	if userName != "" {

		fmt.Fprintf(response, , userName)

	} else {
		http.Redirect(response, request, "/", 302)

	}
}
*/

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "internalPage.html", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login_new.html", nil)
}

func main() {
	fs := http.FileServer((http.Dir("assets")))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	s := SideCarServer{}
	ddos := ddosHandler{}
	ddos.SetNext(&login_handler{})
	//router.HandleFunc("/", indexPageHandler)
	/*router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/", homeHandler)
	*/
	// server main method
	router := mux.NewRouter()

	router.PathPrefix("/login").Handler(ddos.Handle(&s))
	http.Handle("/", router)
	http.ListenAndServe(":80", nil)

}

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
			rw.WriteHeader(401)
			rw.Write([]byte("verpiss dich"))
		}

	})

}

type SideCarHandler interface {
	Handle(s *SideCarServer) http.Handler
	SetNext(SideCarHandler) SideCarHandler
}

type SideCarServer struct {
}
