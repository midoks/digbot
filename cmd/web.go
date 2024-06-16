package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "This command starts web service",
	Description: `Start Web Service`,
	Action:      WebRun,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func WebRun(c *cli.Context) error {
	r := gin.Default()
	r.GET("/add_task", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8090") // 监听并在 0.0.0.0:8080 上启动服务
	return nil
}
