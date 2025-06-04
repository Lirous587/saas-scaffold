package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"os"
	"sass-scaffold/internal/common/logger"
	"sass-scaffold/internal/common/metrics"
	"sass-scaffold/internal/common/server"
	"sass-scaffold/internal/user"
)

func main() {
	var err error

	if err = logger.Init(); err != nil {
		panic(errors.WithMessage(err, "logger模块初始化失败"))
	}

	//ctx := context.Background()
	//metricsClient := metrics.NewPrometheusClient()
	//metrics.StartPrometheusServer()

	server.RunHttpServer(os.Getenv("SERVER_PORT"), metrics.NoOp{}, func(r *gin.RouterGroup) {
		user.InitV1(r)
	})
}
