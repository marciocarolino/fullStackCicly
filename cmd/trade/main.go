package main


import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/marciocarolino/fullStackCycle/internal/infra/kafka"
	"github.com/marciocarolino/fullStackCycle/internal/market/dto"
	"github.com/marciocarolino/fullStackCycle/internal/market/entity"
	"github.com/marciocarolino/fullStackCycle/internal/market/transformer"
)

func main(){
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)

	wg := &sync.WaitGroup{}
	defer wg.Wait()


	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.configMap{
		"bootstrap.servers":"host.docker.internal:9094",
		"group.id":	"myGroup",
		"auto.offset.reset": "earliest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChan)  // T2

	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // T3

	go func(){
		for msg := range KafkaMsgChan {
			wg.Add(1)
			tradeInput := dtp.tradeInput{}
			err := json.Unmarshal(msg.Valeu, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut{
		output := transformer.TransformOutPut(res)
		outputJson, err := json.Marshal(output)
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}