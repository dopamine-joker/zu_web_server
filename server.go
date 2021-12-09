package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/dopamine-joker/zu_web_server/api/router"
	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	misc.Init()
	rpc.InitLogicRpcClient()
	r := router.Register()
	port := misc.Conf.Api.ListenPort

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
		log.Printf("router start run on port %d", port)
	}()

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	misc.Logger.Info("logic server exit")
	log.Println("web_server exit")
}
