package worker

import (
	"context"

	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/mail"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ShutDown()
	ProcessTaskSendVerifyEmail(context.Context, *asynq.Task) error
}

const (
	QUEUE_CRITICAL = "critical"
	QUEUE_DEFAULT  = "default"
)

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {

	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QUEUE_CRITICAL: 10,
			QUEUE_DEFAULT:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("process task failed")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {

	mux := *asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(&mux)
}

func (processor *RedisTaskProcessor) ShutDown() {
	processor.server.Shutdown()
}
