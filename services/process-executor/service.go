package process

import (
	"fmt"
	"log"

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
}

func (s *Service) Init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if err := config.Load(&s.config); err != nil {
		panic(fmt.Errorf("Config error: %w", err))
	}
	s.cache = redis.RedisJobCache(s.config.Redis.URL, s.config.Redis.Username, s.config.Redis.Password, s.config.Redis.DB)
	s.queue = queues.NatsSubscriber(s.config.Nats.URL, "lro", s.cache)
	s.store = db.DBJobStore(s.config.Database.URL)
	s.runner = process.DefaultRunner(s.queue, s.store)
}

func (s *Service) RegisterProcesses(p process.Process) {
	s.runner.Register(p)
}
func (s *Service) Run() {
	s.runner.Run()
}
