package express

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	*http.Server
}

type GenericResponse map[string]string

type Handler = func(http.ResponseWriter, *http.Request)

type Route map[string]Handler

type Express struct {
	routes       StringSet
	getRoutes    StringSet
	getHandlers  Route
	postRoutes   StringSet
	postHandlers Route
}

func response404Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// note to self: this goes strictly AFTER else it does not
	// get parsed as JSON lul
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(GenericResponse{"message": "Error"})
}

func New() *Express {
	return &Express{
		routes:       []string{},
		getRoutes:    []string{},
		postRoutes:   []string{},
		getHandlers:  make(Route),
		postHandlers: make(Route),
	}
}

func (app *Express) Get(pathname string, handler Handler) {
	app.getHandlers[pathname] = handler
	app.routes.Add(pathname)
	app.getRoutes.Add(pathname)
}

func (app *Express) Post(pathname string, handler Handler) {
	app.postHandlers[pathname] = handler
	app.routes.Add(pathname)
	app.postRoutes.Add(pathname)
}

func (app *Express) Listen(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentPath := r.URL.Path
		log.Printf("%s `%s`\n", r.Method, currentPath)

		isPathExists := app.routes.Contains(currentPath)
		if !isPathExists {
			response404Handler(w, r)
			return
		}

		if r.Method == http.MethodGet && app.getRoutes.Contains(currentPath) {
			getHandler := app.getHandlers[currentPath]
			getHandler(w, r)
			return
		}

		if r.Method == http.MethodPost && app.postRoutes.Contains(currentPath) {
			postHandler := app.postHandlers[currentPath]
			postHandler(w, r)
			return
		}

		response404Handler(w, r)
	})

	PORT := fmt.Sprintf(":%d", port)
	log.Printf("HTTP Server listening on PORT %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}