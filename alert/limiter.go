package alert

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type svc struct {
	limiter    *rate.Limiter
	lastExcept time.Time
}

var services = make(map[string]*svc)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanServices()
}

// Retrieve and return the rate limiter for the current service if it
// already exists. Otherwise create a new rate limiter and add it to
// the services map, using the service as the key.
func getService(service string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	s, exists := services[service]
	if !exists {
		// Limit defines the maximum frequency of some events. Limit is represented as number of events per second. A zero Limit allows no events.
		// https://pkg.go.dev/golang.org/x/time/rate#Limit
		// limit 20 second can only triger 1 time
		limiter := rate.NewLimiter(rate.Every(20*time.Second), 1)
		// Include the current time when creating a new service.
		services[service] = &svc{limiter, time.Now()}
		return limiter
	}

	return s.limiter
}

// BackGround clean service map
func cleanServices() {
	for {
		fmt.Println(services)
		time.Sleep(30 * time.Second)
		mu.Lock()
		for service, v := range services {
			if time.Since(v.lastExcept) > 60*time.Second {
				fmt.Println(v.lastExcept)
				delete(services, service)
			}
		}
		mu.Unlock()
	}
}
