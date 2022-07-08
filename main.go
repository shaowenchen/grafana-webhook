package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/shaowenchen/grafana-webhook/config"
	"github.com/shaowenchen/grafana-webhook/pkg/notification"
)

func init() {
	configpath := flag.String("c", "", "")
	flag.Parse()
	config.ReadConfig(*configpath)
}

func main() {
	gin.SetMode(config.Config.Gin.RunMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "")
	})
	router.POST("/", func(c *gin.Context) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Errorf("read request body error,%v", err)
		}
		notification.SendXieZuo(jsonData)
		c.String(200, "")
	})
	router.Run(":8000")
}
