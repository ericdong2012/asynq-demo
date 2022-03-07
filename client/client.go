package main

import (
	"encoding/json"
	"log"
	"asynq-demo/tpl"
	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:36379"
const redisPwd = "G62m50oigInC3111"

func main() {
	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPwd,
		},
	)

	payload, err := json.Marshal(tpl.EmailPayload{Email: "11111111111@qq.com", Content: "发邮件呀111"})
	if err != nil {
		log.Fatal(err)
	}

	task := asynq.NewTask(tpl.EMAIL_TPL, payload)

	_, err = client.Enqueue(task)
	if err != nil {
		log.Fatal(err)
	}

	//_, err2 := client.Enqueue(task, asynq.ProcessIn(10*time.Second))
	//if err2 != nil {
	//	log.Fatal(err2)
	//}

	//_, err3 := client.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(10*time.Second), asynq.Deadline(time.Now().Add(40*time.Second)))
	//if err3 != nil {
	//	log.Fatal(err3)
	//}
}
