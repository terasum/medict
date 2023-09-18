package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/terasum/medict/internal/entry"
	"github.com/terasum/medict/pkg/backserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := entry.DefaultConfig()
	if err != nil {
		panic(err)
	}

	//dictService, err := service.NewDictService(config)
	//if err != nil {
	//	panic(err)
	//}

	staticService, err := backserver.NewStaticServer(config)
	staticService.SetDebug()

	// 开始服务
	log.Info("Starting listen server...")
	staticService.Start()
	// Wait for interrupt signal to gracefully shut down the bs with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Info("Shutdown Server(quit) ...")
		staticService.GracefulStop()

	case <-staticService.StopChan:
		log.Info("Shutdown Server (stopchan) ...")
		staticService.GracefulStop()
	}
}

func startErrorSignalListen(errorChannel chan error, stopChannel chan int) {
	for {
		select {
		case err := <-errorChannel:
			fmt.Printf("err channel %s ", err.Error())
		case <-stopChannel:
			return
		}
	}
}

func startStopSignalListen(stopChannel chan int) {
	for {
		select {
		case <-stopChannel:
			return
		}
	}
}
