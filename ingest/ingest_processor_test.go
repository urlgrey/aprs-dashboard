package ingest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zencoder/disque-go/disque"
	"golang.org/x/net/context"
)

func TestIngestProcessorStop(t *testing.T) {
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewIngestProcessor(pool, "testqueue")
	go h.Run()
	h.Stop()
}

func TestIngestProcessorRun(t *testing.T) {
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewIngestProcessor(pool, "testqueue")

	ctx := context.Background()
	context.WithTimeout(ctx, time.Second)
	var err error
	var conn *disque.Disque
	if conn, err = pool.Get(ctx); err != nil {
		t.Fail()
	}
	defer pool.Put(conn)

	conn.Push("testqueue", "asdfjob", time.Second)

	var queueLength int
	queueLength, err = conn.QueueLength("testqueue")
	assert.Equal(t, 1, queueLength)
	assert.Nil(t, err)

	go h.Run()
	time.Sleep(time.Second)
	h.Stop()

	queueLength, err = conn.QueueLength("testqueue")
	assert.Equal(t, 0, queueLength)
	assert.Nil(t, err)
}

func TestIngestProcessorRunWithBrokenDisqueConnection(t *testing.T) {
	pool := createDisquePool("127.0.0.1:8811")
	h := NewIngestProcessor(pool, "testqueue")

	go h.Run()
	time.Sleep(time.Second)
	h.Stop()
}

func createDisquePool(server string) (pool *disque.DisquePool) {
	hosts := []string{server}       // array of 1 or more Disque servers
	cycle := 1000                   // check connection stats every 1000 Fetch's
	capacity := 10                  // initial capacity of the pool
	maxCapacity := 10               // max capacity that the pool can be resized to
	idleTimeout := 15 * time.Minute // timeout for idle connections
	return disque.NewDisquePool(hosts, cycle, capacity, maxCapacity, idleTimeout)
}
