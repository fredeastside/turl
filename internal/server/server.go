package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"turl/internal/url"
)

type Server struct {
	engine *gin.Engine
	sh     url.Shortener
}

func NewServer(sh url.Shortener) *Server {
	e := gin.Default()
	e.RedirectTrailingSlash = false
	s := &Server{
		sh:     sh,
		engine: e,
	}
	s.setupRoutes()

	return s
}

func (s *Server) Run(addr ...string) error {
	return s.engine.Run(addr...)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.engine.ServeHTTP(w, req)
}

func (s *Server) setupRoutes() {
	s.engine.GET("/health", s.HealthGET)
	v1 := s.engine.Group("/api/v1")
	{
		v1.POST("/", s.MakeShort)
		v1.GET("/:url", s.GetLong)
		v1.GET("/:url/report", s.Report)
	}
}
