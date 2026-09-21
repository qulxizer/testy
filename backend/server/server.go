package server

import (
	"context"
	"log"
	"net/http"
	"testy/db"
	"testy/middlewares"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr   string
	db     *db.DB
	conn   *pgxpool.Pool
	client *http.Client
}

func NewServer(ctx context.Context, addr string, connStr string) (*Server, error) {
	if err := db.RunMigrations(connStr); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("database migrations applied successfully")
	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	db := db.NewDB(conn)

	return &Server{
		addr:   addr,
		conn:   conn,
		db:     db,
		client: http.DefaultClient,
	}, nil
}

func (s *Server) Serve() {
	mux := http.NewServeMux()

	middlewares.StartSessionCleaner(s.db, 30*time.Minute)
	stack := middlewares.NewStack(
		middlewares.CORSMiddleware,
		middlewares.AuthMiddleware(s.db),
	)
	mux.HandleFunc("/api/login", s.Login)
	mux.Handle("/api/test", stack.Then(http.HandlerFunc(s.SubmitHandler)))

	log.Println("listening: http://127.0.0.1:8090")

	if err := http.ListenAndServe("0.0.0.0:8090", mux); err != nil {
		log.Fatalf("err: %v", err)
	}
}
