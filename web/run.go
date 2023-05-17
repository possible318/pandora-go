package web

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/pandora_go/exts/logger"
	"github.com/pandora_go/web/api/chatgpt"
	"github.com/spf13/viper"
	"net/http"
)

func checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Use-Token") == "default" {
			logger.Info("access token miss")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"errorMessage": "access token miss",
			})
		}
		c.Next()
	}
}

func Run(hosts string) {
	gin.ForceConsoleColor()

	// 开启debug模式
	if viper.GetBool("verbose") {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	// 开启sentry
	if viper.GetBool("sentry") {
		router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}

	// 加载模板
	router.LoadHTMLGlob("./web/view/local/templates/*")
	// 设置静态文件目录 ./tpl/local/static
	router.StaticFS("/fonts", http.Dir("./web/view/local/static/fonts/"))
	router.Static("/_next", "./web/view/local/static/_next/")
	router.StaticFile("/favicon-32x32.png", "./web/view/local/static/favicon-32x32.png")
	router.StaticFile("/favicon-16x16.png", "./web/view/local/static/favicon-16x16.png")

	group := router.Group("/api", checkToken())
	{
		group.GET("/models", chatgpt.ListModels)
		group.GET("/conversations", chatgpt.ListConversations)     // 获取会话列表
		group.DELETE("/conversations", chatgpt.ClearConversations) // 清空会话列表
		group.GET("/conversation/:id", chatgpt.GetConversation)
		group.DELETE("/conversation/:id", chatgpt.DelConversation) // 删除会话
		group.PATCH("/conversation/:id", chatgpt.SetConversationTitle)
		group.POST("/conversation/gen_title/:id", chatgpt.GenConversationTitle)

		group.POST("/conversation/talk", chatgpt.Talk)
		group.POST("/conversation/regenerate", chatgpt.Regenerate)
		group.POST("/conversation/goon", chatgpt.Goon)

		group.GET("/accounts/check", chatgpt.Check)
		group.GET("/auth/session", chatgpt.Session)
	}

	router.GET("/", chatgpt.Chat)
	router.GET("/chat", chatgpt.Chat)
	router.GET("/chat/:id", chatgpt.Chat)

	logger.Info("Server is running at " + hosts)
	err := router.Run(hosts)
	if err != nil {
		logger.Error("Failed to start server:" + err.Error())
	}
}
