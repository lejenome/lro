package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lejenome/lro/pkg/config"
	"github.com/lejenome/lro/services/process-api/web/server"
	"github.com/lejenome/lro/services/process-executor/lib/process"
	"github.com/lejenome/lro/services/process-executor/lib/process/db"
	"github.com/lejenome/lro/services/process-executor/lib/process/queues"
	"github.com/lejenome/lro/services/process-executor/lib/process/redis"
)

type Service struct {
	config ProcessApiConfig
	cache  process.JobCache
	queue  process.Queue
	store  process.JobStore
	server server.Server
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
	s.queue = queues.NatsPublisher(s.config.Nats.URL, "lro", s.cache)
	s.store = db.DBJobStore(s.ctx, s.config.Database.URL)
	s.server = server.New(s.ctx, s.config.App.Env, s.config.Auth)
	s.server.Setup()
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

func (s *Service) ScheduleJob(job *process.Job) {
	if err := s.store.Save(job); err != nil {
		fmt.Printf("%s\n", err)
	}
	if err := s.queue.SafeAdd(job); err != nil {
		fmt.Printf("%s\n", err)
	}
}
func (s *Service) Run() {
	defer s.Shutdown()
	s.server.ListenAndServe(s.config.Server.Address)
}
func (s *Service) Shutdown() {
	s.cancel()
}
