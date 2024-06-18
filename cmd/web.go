package cmd

import (

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/files" // swagger embed files
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

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func WebRun(c *cli.Context) error {
	r := gin.Default()
	
	r.GET("/api/add_task", AddTask)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8090") // 监听并在 0.0.0.0:8080 上启动服务
	return nil
}

// @Summary 更新文章
// @Produce json
// @Param tag_id body string false "标签ID"
// @Param title body string false "文章标题"
// @Param desc body string false "文章简述"
// @Param cover_image_url body string false "封面图片地址"
// @Param content body string false "文章内容"
// @Param modified_by body string true "修改者"
// @Success 200 {object} Article "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/articles/{id} [put]
func AddTask(c *gin.Context){
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
