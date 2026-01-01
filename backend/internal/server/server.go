package server

import (
	"database/sql"
	"fmt"
	"mini-app-backend/internal/config"
	"mini-app-backend/internal/handlers"
	"mini-app-backend/internal/handlers/avito"
	"mini-app-backend/internal/logger"
	"mini-app-backend/internal/middleware"
	"mini-app-backend/internal/user"
	"net/http"

	_ "github.com/lib/pq"
)

type Server struct {
	config      *config.Config
	db          *sql.DB
	authHandler *handlers.AuthHandler
	userService *user.UserService
	userRepo    *user.SQLRepository
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) initDB() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.config.PostgresUser, s.config.PostgresPassword, s.config.PostgresHost, s.config.PostgresPort, s.config.PostgresDB)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	s.db = db

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	logger.GetLogger().Info("✅ Connected to database")

	s.userRepo = user.NewSQLRepository(db)

	err = s.userRepo.CreateTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	logger.GetLogger().Info("✅ Database tables created")

	return nil
}

func (s *Server) initServices() {
	s.userService = user.NewUserService(s.userRepo)

	s.authHandler = handlers.NewAuthHandler(s.userService, s.config.TelegramBotToken, s.db)
}

func (s *Server) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.rootHandler)
	mux.HandleFunc("/health", s.healthHandler)

	mux.HandleFunc("/api/auth/telegram", s.authHandler.TelegramAuth)

	mux.HandleFunc("/api/user/me", s.authHandler.GetUser)
	mux.HandleFunc("GET /api/user/data", s.authHandler.GetUserData)
	mux.HandleFunc("POST /api/user/data", s.authHandler.SaveUserData)
	mux.HandleFunc("GET /api/avito/items/", avito.GetItems)
	mux.HandleFunc("GET /api/avito/messenger/chats/", avito.GetMesseges)
}

func (s *Server) Start() error {
	err := s.initDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	defer s.db.Close()

	s.initServices()

	mux := http.NewServeMux()

	s.setupRoutes(mux)

	// Apply middleware chain
	handler := middleware.Logging(middleware.CORS(middleware.RecoverPanic(middleware.ContentTypeJSON(mux))))

	port := s.config.ServerPort

	addr := fmt.Sprintf(":%s", port)
	logger.GetLogger().Infof("🚀 Starting HTTP server on PORT with %s", port)

	return http.ListenAndServe(addr, handler)
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "🤖 Server start! PORT: %s", s.config.ServerPort)
	logger.GetLogger().Infof("📝 Request to the root path: %s", r.RemoteAddr)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "port": "%s"}`, s.config.ServerPort)
	logger.GetLogger().Infof("💓 Check health от %s", r.RemoteAddr)
}
