package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

type application struct {
	port        int
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		port:        5000,
		infoLogger:  infoLog,
		errorLogger: errorLog,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLogger.Println("listening on port: ", app.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	var out []byte
	out = []byte("ola")

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(out)
		if err != nil {
			app.errorLogger.Println(err)
		}
	})

	return mux
}
