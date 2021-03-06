package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zencoder/disque-go/disque"
	"golang.org/x/net/context"
)

type HealthCheckHandler struct {
	pool *disque.DisquePool
}

func NewHealthCheckHandler(pool *disque.DisquePool) *HealthCheckHandler {
	return &HealthCheckHandler{pool: pool}
}

// Add routes to router
func InitializeRouterForHealthCheckHandler(r *mux.Router, pool *disque.DisquePool) {
	m := NewHealthCheckHandler(pool)
	r.HandleFunc("/healthcheck", m.HealthCheck).Methods("GET")
}

// Examine and report the health of the component and dependencies
func (h *HealthCheckHandler) HealthCheck(resp http.ResponseWriter, req *http.Request) {
	var c *disque.Disque
	var err error
	ctx := context.Background()
	context.WithTimeout(ctx, 5*time.Second)
	if c, err = h.pool.Get(ctx); err != nil {
		http.Error(resp,
			fmt.Sprintf("{\"error\": \"Error getting Disque connection %+v\"}", err),
			http.StatusInternalServerError)
		return
	}
	defer h.pool.Put(c)

	if _, err = c.QueueLength("queueName"); err != nil {
		http.Error(resp,
			fmt.Sprintf("{\"error\": \"Error querying Disque for queue length %+v\"}", err),
			http.StatusInternalServerError)
		return
	}
}
