package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiReponse struct {
	Version string        `json:"version"`
	Items   []interface{} `json:"items"`
}

func ApiVersion(c *gin.Context) {
	ap := ApiReponse{}
	ap.Version = "1.0.0"
	ap.Items = []interface{}{"test1", "test2"}
	c.JSON(http.StatusOK, ap)

}
