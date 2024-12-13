package worker

import (
	"context"

	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(context.Context, *asynq.Task) error
}

const (
	QUEUE_CRITICAL = "critical"
	QUEUE_DEFAULT  = "default"
)

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {

	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QUEUE_CRITICAL: 10,
			QUEUE_DEFAULT:  5,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {

	mux := *asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(&mux)
}
