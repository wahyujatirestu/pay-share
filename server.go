package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"

	"github.com/wahyujatirestu/payshare/config"
	"github.com/wahyujatirestu/payshare/controller"
	// "github.com/wahyujatirestu/payshare/middleware"
	"github.com/wahyujatirestu/payshare/repository"
	"github.com/wahyujatirestu/payshare/routes"
	service "github.com/wahyujatirestu/payshare/service"
	repo "github.com/wahyujatirestu/payshare/utils/repo"
	utilsservice "github.com/wahyujatirestu/payshare/utils/service"
)

type Server struct {
	UserService      service.UserService
	AuthService      service.AuthenticationService
	JWTService       utilsservice.JWTService
	RefreshTokenRepo repo.RefreshTokenRepository
	UserRepo         repository.UserRepository
	DB               *sql.DB
	Engine           *gin.Engine
	Host             string
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repo.NewRefreshTokenRepository(db)

	userService := service.NewUserService(userRepo)
	jwtService := utilsservice.NewJWTService(cfg.TokenConfig)
	authService := service.NewAuthenticationService(userService, jwtService, refreshTokenRepo)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		UserService:      userService,
		AuthService:      authService,
		JWTService:       jwtService,
		RefreshTokenRepo: refreshTokenRepo,
		UserRepo:         userRepo,
		DB:               db,
		Engine:           engine,
		Host:             host,
	}
}

func (s *Server) SetupRoutes() {
	apiV1 := s.Engine.Group("/api/v1")

	userController := controller.NewUserController(s.UserService)
	authController := controller.NewAuthController(s.AuthService)

	routes.AuthRoute(apiV1, authController, userController)
	routes.UserRoute(apiV1, userController, s.JWTService)
}

func (s *Server) Run() {
	s.SetupRoutes()
	if err := s.Engine.Run(s.Host); err != nil {
		log.Fatalf("failed to run server on %s: %v", s.Host, err)
	}
}

func (s *Server) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}
