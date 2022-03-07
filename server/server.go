package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"asynq-demo/tpl"
	"golang.org/x/sys/unix"
	"github.com/hibiken/asynq"
)

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		log.Printf("Start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			log.Fatalf("error:  %q", err.Error())
			return err
		}
		log.Printf("Finished processing %q: Elapsed Time = %v", t.Type(), time.Since(start))
		return nil
	})
}

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:36379", Password: "G62m50oigInC3111"},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(loggingMiddleware)
	mux.HandleFunc(tpl.EMAIL_TPL, emailMqHandler)

	if err := srv.Start(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		s := <-sigs
		if s == unix.SIGTSTP {
			srv.Shutdown()
			continue
		}
		break
	}

	srv.Stop()
}

func emailMqHandler(ctx context.Context, t *asynq.Task) error {
	var p tpl.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("emailMqHandler err:%+v", err)
	}
	fmt.Printf("p : %+v \n", p)

	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     "127.0.0.1:36379",
			Password: "G62m50oigInC30sf",
		},
	)

	payload, err := json.Marshal(tpl.EmailPayload{Email: "22222222222@qq.com", Content: "发邮件呀222"})
	if err != nil {
		log.Fatal(err)
	}

	task := asynq.NewTask(tpl.EMAIL_TPL2, payload)
	_, err = client.Enqueue(task)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
