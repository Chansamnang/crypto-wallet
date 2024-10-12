package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	mysqlConfig "wallet/config"
	"wallet/internal/router"
	"wallet/pkg/common/config"
	"wallet/pkg/hdwallet/eth"
	"wallet/pkg/hdwallet/tron"
	"wallet/pkg/network"
	"wallet/pkg/zlogger"
)

func main() {
	configFile, logFile, _, _, _, err := config.FlagParse("wallet")
	if err != nil {
		panic(err)
	}
	err = config.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// Init log
	zlogger.InitLogConfig(logFile)

	// Init DB
	mysqlConfig.InitDB()

	engine := gin.Default()
	router.InitRouter(engine)

	// Connect Tron Grpc
	tron.Init()

	// Connect Eth
	eth.Init()

	ip, err := network.GetLocalIP()
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("%s:%d", ip, config.Config.App.Port)
	server := &http.Server{
		Addr:    host,
		Handler: engine,
	}

	zlogger.Infof("Server listening on %s", server.Addr)

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlogger.Errorw("listen", zap.Error(err))
		}
	}()

	zlogger.Infof("wallet api started, host %s", host)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zlogger.Errorw("listen", zap.Error(err))
	}

	zlogger.Infow("wallet api exiting")
}
