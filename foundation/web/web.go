package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httptreemux.ContextMux // embedding = type promotion. Everythnuig from inner type promotes to the outer type
	shutdown                chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method string, group string, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		// Pre code procesing here

		if err := handler(r.Context(), w, r); err != nil {
			// error handling
			return
		}

		// Post code processing here
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	// We need to feed a function matching the signature that a ContextMux expects since that is the
	// lowest level mux that will actuall repsond to requests.
	// But by defining h within this larger function that wraps it, we can now access the function's
	// variables within h.
	a.ContextMux.Handle(method, finalPath, h)
}
