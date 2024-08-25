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

type Options struct {
	Port           uint16
	Timeout        time.Duration
	CurrentService string
	UseTLS         bool
	CertPath       string
	KeyPath        string
}

type Dependecies struct{}

func NewGRPC(logger *zerolog.Logger, opts Options, _ Dependecies) *App {
	var (
		creds credentials.TransportCredentials
		err   error
	)

	grpcLogger := logger.With().Logger()

	if opts.UseTLS {
		grpcLogger.Info().
			Str("cert_path", opts.CertPath).
			Str("key_path", opts.KeyPath).
			Msg("grpc uses TLS certificate")

		creds, err = credentials.NewServerTLSFromFile(opts.CertPath, opts.KeyPath)
		if err != nil {
			panic(err)
		}
	} else {
		creds = insecure.NewCredentials()
	}

	serv := grpc.NewServer(
		grpc.ConnectionTimeout(opts.Timeout),
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(&grpcLogger),
			recoverInterceptor(&grpcLogger, opts.CurrentService),
		),
		grpc.Creds(creds),
	)

	// register

	return &App{
		log:    logger,
		server: serv,
		port:   fmt.Sprintf(":%d", opts.Port),
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
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
