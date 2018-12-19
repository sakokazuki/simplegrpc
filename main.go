package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sakokazuki/simplegrpc/config"
	"github.com/sakokazuki/simplegrpc/pubsub"
	"github.com/sakokazuki/simplegrpc/server"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func grpcListener(config config.Config) net.Listener {
	l, err := net.Listen("tcp", ":"+config.GrpcPort)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to listen")
	}
	return l
}

func main() {
	config := config.New()

	gl := grpcListener(config)
	defer gl.Close()

	pubsuber := pubsub.NewPubSub()

	grpcServerOption := server.Option{
		Config:   config,
		PubSuber: pubsuber,
	}
	grpcServer, err := server.NewGRPCServer(grpcServerOption)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failded to create gRPC server")
	}

	// for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	go func() {
		<-sigCh

		eg := errgroup.Group{}
		eg.Go(func() error {
			log.Fatal().Err(err).Msg("shutdown gRPC Server gracefully...")
			grpcServer.GracefulStop()
			return nil
		})

		if err := eg.Wait(); err != nil {
			opErr, ok := err.(*net.OpError)

			// NOTE: Ignore errors that occur when closing the file descriptor because it is an assumed error.
			if ok && opErr.Op == "close" {
				return
			}
			log.Fatal().Err(err).Msg("failed to shutdown gracefully")
		}
	}()

	log.Info().Msgf("gPRC server start at localhost:%v", config.GrpcPort)
	if err := grpcServer.Serve(gl); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to gRPC serve")
	}

}
