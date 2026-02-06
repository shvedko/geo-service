package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shvedko/geo-service/internal/service/geo"
)

type Config struct {
	BindAddr    string `default:"" split_words:"true" desc:"bind address"`
	BindPort    string `default:"8080" split_words:"true" desc:"bind port"`
	DatabaseURL string `default:"postgres://postgres:postgres@postgres:5432/geo" split_words:"true" desc:"data base url"`
}

type Redacted struct {
	s *string
}

func (r Redacted) String() string {
	u, err := url.Parse(*r.s)
	if err != nil {
		return ""
	}
	return u.Redacted()
}

func (r Redacted) Set(s string) error {
	*r.s = s
	return nil
}

func (r Redacted) Type() string {
	return "string"
}

func URL(s *string) pflag.Value {
	return Redacted{s: s}
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
			err := run(ctx, cfg.BindAddr, cfg.BindPort, cfg.DatabaseURL)
			switch {
			case errors.Is(err, http.ErrServerClosed):
				return nil
			default:
				return err
			}
		},
	}

	c.Flags().StringVar(&cfg.BindAddr, "addr", cfg.BindAddr, "bind address")
	c.Flags().StringVar(&cfg.BindPort, "port", cfg.BindPort, "bind port")
	c.Flags().Var(URL(&cfg.DatabaseURL), "base", "data base url")

	err = c.Execute()
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
