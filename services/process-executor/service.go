package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lejenome/lro/pkg/config"
	"github.com/lejenome/lro/services/process-executor/lib/process"
	"github.com/lejenome/lro/services/process-executor/lib/process/db"
	"github.com/lejenome/lro/services/process-executor/lib/process/queues"
	"github.com/lejenome/lro/services/process-executor/lib/process/redis"
)

type Service struct {
	config ProcessExecutorConfig
	cache  process.JobCache
	queue  process.Queue
	store  process.JobStore
	runner process.Runner
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *Service) Init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if err := config.Load(&s.config); err != nil {
		panic(fmt.Errorf("Config error: %w", err))
	}
	s.initCtxt()
	s.cache = redis.RedisJobCache(s.ctx, s.config.Redis.URL, s.config.Redis.Username, s.config.Redis.Password, s.config.Redis.DB)
	s.queue = queues.NatsSubscriber(s.ctx, s.config.Nats.URL, "lro", s.cache)
	s.store = db.DBJobStore(s.ctx, s.config.Database.URL)
	s.runner = process.DefaultRunner(s.queue, s.store)
}

func (s *Service) initCtxt() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.ctx = ctx
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()
}

func (s *Service) RegisterProcesses(p process.Process) {
	s.runner.Register(p)
}
func (s *Service) Run() {
	defer s.Shutdown()
	s.runner.Run(s.ctx)
}
func (s *Service) Shutdown() {
	s.cancel()
}
