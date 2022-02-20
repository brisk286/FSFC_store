package main

import (
	"fsfc_store/logger"
	"fsfc_store/router"
	"net/http"
	"time"
)

func main() {
	//logger.InitLogger(config.GetConfig().Log.Path, config.GetConfig().Log.Level)

	//设置路由
	newRouter := router.NewRouter()

	//在本地开一个端口  接收信息
	s := &http.Server{
		Addr:           ":5555",
		Handler:        newRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if nil != err {
		logger.Logger.Error("server error", logger.Any("serverError", err))
	}
}
