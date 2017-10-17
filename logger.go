package tamework

import (
	"log"
	"time"
)

func Logger() Handler {
	return func(c *Context) {
		started := time.Now()
		log.Printf("start [%s] %s", c.Method, c.Text)
		c.Next()
		log.Printf("done for %s [%s] %s", time.Since(started), c.Method, c.Text)
	}
}
