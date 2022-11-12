package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

type application struct {
	port        int
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

type user struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type tweet struct {
	Username string `json:"username"`
	Tweet    string `json:"tweet"`
	Avatar   string `json:"avatar,omitempty"`
}

var users []user
var tweets []tweet

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
	mux.Use(middleware.Recoverer)

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

	mux.Post("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		var newUser user

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&newUser)
		if err != nil {
			app.errorLogger.Println(err)
		}

		users = append(users, newUser)

		app.infoLogger.Println(newUser.Avatar, newUser.Username)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.Post("/tweets", func(w http.ResponseWriter, r *http.Request) {
		var newTweet tweet

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&newTweet)
		if err != nil {
			app.errorLogger.Println(err)
		}

		newTweet.Avatar = users[len(users)-1].Avatar

		tweets = append(tweets, newTweet)
		app.infoLogger.Println(tweets)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.Get("/tweets", func(w http.ResponseWriter, r *http.Request) {

		var lastTen []tweet

		var payload []byte
		var _ error

		if len(tweets) > 10 {
			for i := len(tweets) - 1; i > len(tweets)-11; i-- {
				lastTen = append(lastTen, tweets[i])
			}
			payload, _ = json.MarshalIndent(lastTen, "", "\t")
		} else {
			payload, _ = json.MarshalIndent(tweets, "", "\t")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(payload)

	})

	mux.With(UsernameCtx).Get("/tweets/{username}", func(w http.ResponseWriter, r *http.Request) {
		app.infoLogger.Println(r.Context().Value("username"))

		username := r.Context().Value("username")

		var t []tweet

		for i := 0; i < len(tweets); i++ {
			if tweets[i].Username == username {
				t = append(t, tweets[i])
			}
		}

		if len(t) == 0 {

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("no matching user"))
			return
		}

		app.infoLogger.Println(len(t))

		payload, _ := json.MarshalIndent(t, "", "\t")

		w.WriteHeader(http.StatusOK)
		w.Write(payload)

	})

	return mux
}

func UsernameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := chi.URLParam(r, "username")

		ctx := context.WithValue(r.Context(), "username", userName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
