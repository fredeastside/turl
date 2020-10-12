package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Server) MakeShort(c *gin.Context) {
	longUrl := c.PostForm("url")
	if !isValidUrl(longUrl) {
		c.Status(http.StatusBadRequest)
		return
	}
	shortUrl, err := s.sh.Encode(longUrl)
	if err != nil {
		log.Errorf("%s encoding err: %v", longUrl, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, shortUrl)
}

func (s *Server) GetLong(c *gin.Context) {
	shortUrl := c.Param("url")
	longUrl, err := s.sh.Decode(shortUrl)
	if err != nil {
		log.Errorf("%s decoding err: %v", shortUrl, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusMovedPermanently, longUrl)
}

func (s *Server) Report(c *gin.Context) {
	shortUrl := c.Param("url")
	daily, err := s.sh.GetDailyCount(shortUrl)
	if err != nil {
		log.Errorf("%s daily count err: %v", shortUrl, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	weekly, err := s.sh.GetWeeklyCount(shortUrl)
	if err != nil {
		log.Errorf("%s weekly count err: %v", shortUrl, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	all, err := s.sh.GetCount(shortUrl)
	if err != nil {
		log.Errorf("%s count err: %v", shortUrl, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.String(
		http.StatusOK, fmt.Sprintf("%s\n%s\n%s", strconv.Itoa(daily), strconv.Itoa(weekly), strconv.Itoa(all)))
}

func isValidUrl(s string) bool {
	u, err := url.Parse(s)

	return err == nil && u.Scheme != "" && u.Host != ""
}
