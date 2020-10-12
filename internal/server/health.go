package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) HealthGET(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
