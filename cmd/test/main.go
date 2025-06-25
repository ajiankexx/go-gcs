package main

import (
	"net/http" // 导入 http 包，用于状态码

	"github.com/gin-gonic/gin" // 导入 gin 框架

	swaggerFiles "github.com/swaggo/files" // 用于提供 Swagger UI 的静态文件
	ginSwagger "github.com/swaggo/gin-swagger" // Gin 框架与 Swagger UI 的适配器
)

// @title Go GCS API  // API 标题
// @version 1.0  // API 版本
// @description This is a sample API for GCS project. // API 描述
// @termsOfService http://swagger.io/terms/ // 服务条款链接 (可选)

// @contact.name API Support // 联系人姓名
// @contact.url http://www.swagger.io/support // 联系人 URL (可选)
// @contact.email support@swagger.io // 联系人邮箱 (可选)

// @license.name Apache 2.0 // 许可证名称
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html // 许可证 URL

// @host localhost:8123 // API 服务的实际主机和端口
// @BasePath / // API 的基础路径 (例如，所有路由都以 /api/v1 开头，这里就写 /api/v1)

// @securityDefinitions.basic BasicAuth // 定义一个名为 BasicAuth 的基本认证方式 (可选，如果需要认证)

// @externalDocs.description Open API // 外部文档描述 (可选)
// @externalDocs.url https://swagger.io/resources/open-api/ // 外部文档 URL (可选)

func main() {
	r := gin.Default()

	// Swagger UI 路由，用于在浏览器中查看 API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// @Summary 获取测试消息 // API 摘要，简短描述
	// @Description 这是一个简单的测试接口，用于返回 "hello" 消息。 // 详细描述
	// @Tags 测试 // 标签，用于在 Swagger UI 中对 API 进行分组
	// @Accept json // 该 API 接受的请求内容类型 (如果请求体是 JSON)
	// @Produce json // 该 API 返回的响应内容类型 (如果响应是 JSON)
	// @Success 200 {object} map[string]interface{} "成功返回消息" // 成功响应的 HTTP 状态码和返回对象
	// @Router /test [get] // API 路径和 HTTP 方法
	r.GET("/test", func(c *gin.Context) {
		// 这里将状态码从 202 (Accepted) 改为 200 (OK)，更符合返回 "hello" 消息的场景
		c.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.Run(":8123") // 启动 Gin 服务器
}

