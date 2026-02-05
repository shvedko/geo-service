package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shvedko/geo-service/internal/service/geo"
)

func main() {
	var addrFlag string
	var portFlag string
	var baseFlag string

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := &cobra.Command{
		Use:  "geo",
		Long: "Geo Bounding Box Service",
		RunE: func(*cobra.Command, []string) error {
			err := run(ctx, addrFlag, portFlag, baseFlag)
			switch {
			case errors.Is(err, http.ErrServerClosed):
				return nil
			default:
				return err
			}
		},
	}

	c.Flags().StringVar(&addrFlag, "addr", "", "bind address")
	c.Flags().StringVar(&portFlag, "port", "8080", "bind port")
	c.Flags().StringVar(&baseFlag, "db", "postgres://postgres:postgres@postgres:5432/geo", "data base")

	err := c.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const float = "[-+]?([0-9]+(\\.[0-9]*)?|\\.[0-9]+)"

func run(ctx context.Context, addr string, port string, base string) error {
	db, err := pgxpool.New(ctx, base)
	if err != nil {
		return err
	}
	defer db.Close()

	g, err := geo.New(db)
	if err != nil {
		return err
	}

	h := chi.NewRouter()
	h.Put("/points"+"", g.Put)
	h.Get("/points/{left:"+float+"}/{top:"+float+"}/{right:"+float+"}/{bottom:"+float+"}", g.Get)
	s := http.Server{
		Handler:     h,
		Addr:        net.JoinHostPort(addr, port),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	w := sync.WaitGroup{}
	defer w.Wait()

	context.AfterFunc(ctx, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		w.Add(1)
		defer w.Done()
		_ = s.Shutdown(ctx)
	})

	return s.ListenAndServe()
}
