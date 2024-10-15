package grpcapp

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log    *zerolog.Logger
	server *grpc.Server
	port   string
}

type Config struct {
	Port           int
	Timeout        time.Duration
	CurrentService string
	UseTls         bool
	CertPath       string
	KeyPath        string
}

func New(logger *zerolog.Logger, conf Config) *App {
	var (
		creds credentials.TransportCredentials
		err   error
	)

	grpcLogger := logger.With().Logger()

	if conf.UseTls {
		grpcLogger.Info().
			Str("cert_path", conf.CertPath).
			Str("key_path", conf.KeyPath).
			Msg("grpc uses TLS certificate")

		creds, err = credentials.NewServerTLSFromFile(conf.CertPath, conf.KeyPath)
		if err != nil {
			panic(err)
		}
	} else {
		creds = insecure.NewCredentials()
	}

	serv := grpc.NewServer(
		grpc.ConnectionTimeout(conf.Timeout),
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(&grpcLogger),
			recoverInterceptor(&grpcLogger, conf.CurrentService),
		),
		grpc.Creds(creds),
	)

	// Register servers

	return &App{
		log:    logger,
		server: serv,
		port:   fmt.Sprintf(":%d", conf.Port),
	}
}

func (a *App) Start() error {
	lis, err := net.Listen("tcp", a.port)
	if err != nil {
		a.log.Error().
			Str("error", err.Error()).
			Msg("error while starting gRPC")

		return err
	}

	go func() {
		if err := a.server.Serve(lis); err != nil {
			a.log.Error().Err(err).Msg("error while starting gRPC")
			return
		}
	}()

	a.log.Info().Str("port", a.port).Msg("gRPC started")

	return nil
}

func (a *App) Stop() {
	a.server.GracefulStop()
	a.log.Info().Str("port", a.port).Msg("stopped gRPC")
}
