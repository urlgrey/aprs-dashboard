package ingest

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"github.com/urlgrey/aprs-dashboard/models"
	"github.com/zencoder/disque-go/disque"
	"golang.org/x/net/context"
)

type IngestProcessor struct {
	pool        *disque.DisquePool
	queueName   string
	stopChannel chan bool
	waitGroup   *sync.WaitGroup
	database    *dynamodb.DynamoDB
}

func NewIngestProcessor(pool *disque.DisquePool, queueName string, database *dynamodb.DynamoDB) *IngestProcessor {
	return &IngestProcessor{
		pool:        pool,
		queueName:   queueName,
		stopChannel: make(chan bool),
		waitGroup:   &sync.WaitGroup{},
		database:    database,
	}
}

const (
	CONNECTION_POOL_ERROR_SLEEP_DURATION = 100 * time.Millisecond
)

func (i *IngestProcessor) Run() {
	i.waitGroup.Add(1)
	defer i.waitGroup.Done()

	ctx := context.Background()
	context.WithTimeout(ctx, time.Second)

	for {
		select {
		case <-i.stopChannel:
			log.Println("Stopping the consumer processes")
			return
		default:
		}

		var err error
		var conn *disque.Disque
		if conn, err = i.pool.Get(ctx); err != nil {
			log.Printf("Error while getting connection from pool: %s", err)
			time.Sleep(CONNECTION_POOL_ERROR_SLEEP_DURATION)
			continue
		}
		defer i.pool.Put(conn)

		var job *disque.Job
		if job, err = conn.Fetch(i.queueName, time.Second); err != nil {
			log.Printf("Error retrieving next job from Disque: %s", err)
			continue
		}

		if job != nil {
			log.Printf("Job ID: %s", job.JobId)
			message := &models.AprsMessage{}
			if err = json.Unmarshal([]byte(job.Message), message); err != nil {
				log.Printf("Error while binding message to model: %s", err)
				continue
			}
			log.Printf("Message location: lat=%f, long=%f", message.Latitude, message.Longitude)
		}
	}
}

func (i *IngestProcessor) Stop() {
	close(i.stopChannel) // signal to the Run goroutine to stop
	i.waitGroup.Wait()   // wait for the Run routine and its children to finish
}
