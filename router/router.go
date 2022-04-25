package router

import (
	"fsfc_store/api/v1"
	"fsfc_store/logger"
	"fsfc_store/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	server := gin.Default()
	//添加全局中间件
	server.Use(Cors())
	server.Use(Recovery)

	//socket := RunSocekt

	group := server.Group("v1")
	{
		group.POST("/changedFile", v1.GetChangedFilesAndPostDataList)
		group.POST("/rebuildFile", v1.GetRsyncOpsToRebuild)
		group.POST("/multiDownload", v1.MultiDownload)
		group.POST("/getFilesInfo", v1.GetFilesInfo)
	}
	return server
}

//跨域设置中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//http请求类型
		method := c.Request.Method
		//请求头部
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			//set_key
			c.Header("Access-Control-Allow-Origin", "*") // 可将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		//异常处理
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Error("HttpError", zap.Any("HttpError", err))
			}
		}()

		c.Next()
	}
}

//异常处理中间件
func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("gin catch error: ", logger.Any("gin catch error: ", r))
			c.JSON(http.StatusOK, response.FailMsg("系统内部错误"))
		}
	}()
	c.Next()
}
