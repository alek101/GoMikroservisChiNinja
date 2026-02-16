package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/alek101/GoMikroservisChiNinja/database"
)

type App struct {
	router http.Handler

	db *sql.DB
}

func New() (*App, error) {
	db, err := database.New()
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return nil,err
	}

	app := &App{
		router: loadRoutes(db),
		db : db,
	}

	return app,nil
}

func (a *App) Start(ctx context.Context) error {
	// defer se obavlja kada se posalje signal da se app gasi
	// defer a.db.Close()

	server := http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()
	
	select {
		case err := <-ch: return err
		case <-ctx.Done(): 
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		a.db.Close()
		return server.Shutdown(shutdownCtx)
	}

	// return nil
}
