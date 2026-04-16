package main

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shvedko/geo-service/doc"
	"github.com/shvedko/geo-service/internal/service/geo"
)

type Config struct {
	BindAddr    string  `default:"" split_words:"true" desc:"bind address"`
	BindPort    string  `default:"8080" split_words:"true" desc:"bind port"`
	DatabaseURL url.URL `default:"postgres://postgres:postgres@postgres:5432/geo" split_words:"true" desc:"data base"`
}

type RedactedURL struct {
	u *url.URL
}

func (r RedactedURL) String() string {
	return r.u.Redacted()
}

func (r RedactedURL) Set(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*r.u = *u
	return nil
}

func (r RedactedURL) Type() string {
	return "string"
}

func URL(u *url.URL) pflag.Value {
	return RedactedURL{u: u}
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		_ = envconfig.Usage("", &cfg)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := &cobra.Command{
		Use:  "geo",
		Long: "Geo Bounding Box Service",
		RunE: func(*cobra.Command, []string) error {
			return run(ctx, cfg.BindAddr, cfg.BindPort, cfg.DatabaseURL.String())
		},
	}

	c.Flags().StringVar(&cfg.BindAddr, "addr", cfg.BindAddr, "bind address")
	c.Flags().StringVar(&cfg.BindPort, "port", cfg.BindPort, "bind port")
	c.Flags().Var(URL(&cfg.DatabaseURL), "base", "data base")

	err = c.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, addr string, port string, base string) error {
	db, err := pgxpool.New(ctx, base)
	if err != nil {
		return err
	}
	defer db.Close()

	{
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		err := db.Ping(ctx)
		if err != nil {
			return err
		}
	}

	g, err := geo.New(db)
	if err != nil {
		return err
	}

	const float = "[-+]?([0-9]+(\\.[0-9]*)?|\\.[0-9]+)"

	h := chi.NewRouter()
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)
	h.Use(middleware.RealIP)
	h.Post("/points", g.Post)
	h.Get("/points/{id:[0-9]+}", g.Get)
	h.Get("/points/{min_lon:"+float+"}/{min_lat:"+float+"}/{max_lon:"+float+"}/{max_lat:"+float+"}", g.Box)
	h.Get("/swagger.yaml", doc.Swagger.ServeHTTP)
	s := http.Server{
		Handler:     h,
		Addr:        net.JoinHostPort(addr, port),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	e := make(chan error, 1)

	defer context.AfterFunc(ctx, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		e <- s.Shutdown(ctx)
	})()

	err = s.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return <-e
}
