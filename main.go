package main

import (
	"fsfc_store/logger"
	"fsfc_store/router"
	"fsfc_store/rpc/data_rpc/protocol"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", ":8008")
	if err != nil {
		log.Fatal(err)
	}
	//起一个rpc server
	server := router.NewServer()
	//server注册器
	err = server.Register(new(protocol.RsyncService))
	if err != nil {
		log.Fatal(err)
	}
	go server.Serve(lis)

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
	log.Printf("restful api started on: %s", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		logger.Logger.Error("server error", logger.Any("serverError", err))
	}

}
