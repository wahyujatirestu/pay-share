package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/wahyujatirestu/payshare/config"
)

type Server struct {
	engine *gin.Engine
	host 	string
}

// func (s *Server) Run()  {
// 	s.
// }

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		fmt.Println("error :", err)
	}
	
	
}