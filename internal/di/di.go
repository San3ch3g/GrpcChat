package di

import (
	"ModuleForChat/internal/pkg/rpc"
	"ModuleForChat/internal/pkg/storage/pg"
	"ModuleForChat/internal/utils/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type Container struct {
	cfg         *config.Config
	netListener net.Listener
	grpcServer  *grpc.Server
	storage     *pg.Storage
	db          *gorm.DB
	rpcServer   *rpc.Server
}

func New(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) GetNetListener() net.Listener {
	if c.netListener == nil {
		listener, err := net.Listen("tcp", c.cfg.ServerPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		c.netListener = listener
	}
	return c.netListener
}

func (c *Container) GetGRPCServer() *grpc.Server {
	if c.grpcServer == nil {
		c.grpcServer = grpc.NewServer()
	}
	return c.grpcServer
}

func (c *Container) GetDB() *gorm.DB {
	if c.db == nil {
		c.db = pg.MustNewPostgresDB(c.cfg)
	}
	return c.db
}

func (c *Container) GetSQLStorage() *pg.Storage {
	if c.storage == nil {
		c.storage = pg.New(c.GetDB())
	}
	return c.storage
}
