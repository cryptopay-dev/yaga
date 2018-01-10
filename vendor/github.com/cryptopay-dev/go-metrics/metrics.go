package metrics

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/go-nats"
)

type conn struct {
	mu           sync.RWMutex
	nats         *nats.Conn
	errorHandler ErrorHandler
	enabled      bool
	queue        string
	url          string
	hostname     string
	application  string
}

type ErrorHandler func(err error)

// M metrics storage
// Example:
// m := metrics.M{
// 	"metric": 1000,
//	"gauge": 1,
// }
type M map[string]interface{}

// T tags storage
// Example:
// m := metrics.T{
//	"tag": "some_default_tag"
// }
type T map[string]string

// DefaultConn shared default metric
// connection
var DefaultConn *conn

// DefaultErrorHandler just printing all errors comes from
// async writes to Stderr
var DefaultErrorHandler = func(err error) {
	fmt.Fprintf(os.Stderr, "Error while sending metrics: %v", err)
}

// DefaultQueue is queue where we puts event into NATS
const DefaultQueue = "telegraf"

// Setup rewrites default metrics configuration
//
// Params:
// - url (in e.g. "nats://localhost:4222")
// - options nats.Option array
//
// Example:
// import (
//     "log"
//
//     "github.com/cryptopay.dev/go-metrics"
// )
//
// func main() {
//     err := metrics.Setup("nats://localhost:4222")
//     if err != nil {
//         log.Fatal(err)
//     }
//
//     for i:=0; i<10; i++ {
//         err = metrics.SendAndWait(metrics.M{
//             "counter": i,
//         })
//
//         if err != nil {
//             log.Fatal(err)
//         }
//     }
// }
func Setup(url string, application, hostname string, options ...nats.Option) error {
	metrics, err := New(url, application, hostname, options...)
	if err != nil {
		return err
	}

	DefaultConn = metrics
	return nil
}

// New creates new metrics connection
//
// Params:
// - url (in e.g. "nats://localhost:4222")
// - options nats.Option array
//
// Example:
// import (
//     "log"
//
//     "github.com/cryptopay.dev/go-metrics"
// )
//
// func main() {
//     m, err := metrics.New("nats://localhost:4222")
//     if err != nil {
//         log.Fatal(err)
//     }
//
//     for i:=0; i<10; i++ {
//         err = m.SendAndWait(metrics.M{
//             "counter": i,
//         })
//
//         if err != nil {
//             log.Fatal(err)
//         }
//     }
// }
func New(url string, application, hostname string, options ...nats.Option) (*conn, error) {
	if url == "" {
		return &conn{
			enabled: false,
		}, nil
	}

	// Getting current environment
	if application == "" {
		return nil, errors.New("Application name not set")
	}

	if hostname == "" {
		return nil, errors.New("Hostname not set")
	}

	nc, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}

	conn := &conn{
		nats:         nc,
		hostname:     hostname,
		enabled:      true,
		queue:        DefaultQueue,
		errorHandler: DefaultErrorHandler,
		application:  application,
	}

	return conn, nil
}

// Send metrics to NATS queue
//
// Example:
// m.Send(metrics.M{
// 		"counter": i,
// })
func Send(metrics M, path ...string) {
	if DefaultConn != nil {
		DefaultConn.SendWithTags(metrics, nil, path...)
	}
}

// SendWithTags metrics to NATS queue waiting for response
//
// Example:
// err = SendWithTags(metrics.M{
// 		"counter": i,
// }, metrics.T{
//	    "tag": "sometag",
// }, "metricname")
func SendWithTags(metrics M, tags T, path ...string) {
	if DefaultConn != nil {
		DefaultConn.SendWithTags(metrics, tags, path...)
	}
}

// SendAndWait metrics to NATS queue waiting for response
//
// Example:
// err = SendAndWait(metrics.M{
// 		"counter": i,
// }, "metricname")
func SendAndWait(metrics M, path ...string) error {
	if DefaultConn == nil {
		return nil
	}

	return DefaultConn.SendAndWait(metrics, path...)
}

// SendWithTagsAndWait metrics to NATS queue waiting for response
//
// Example:
// err = SendWithTags(metrics.M{
// 		"counter": i,
// }, metrics.T{
//	    "tag": "sometag",
// }, "metricname")
func SendWithTagsAndWait(metrics M, tags T, path ...string) error {
	if DefaultConn == nil {
		return nil
	}

	return DefaultConn.SendWithTagsAndWait(metrics, tags, path...)
}

// SetErrorHandler changes error handler to providen one
func SetErrorHandler(fn ErrorHandler) {
	if DefaultConn != nil {
		DefaultConn.SetErrorHandler(fn)
	}
}

// SetErrorHandler changes error handler to providen one
func (m *conn) SetErrorHandler(fn ErrorHandler) {
	m.errorHandler = fn
}

// Send metrics to NATS queue
//
// Example:
// m.Send(metrics.M{
// 		"counter": i,
// }, "metricname")
func (m *conn) Send(metrics M, path ...string) {
	go func() {
		if err := m.SendAndWait(metrics, path...); err != nil {
			m.errorHandler(err)
		}
	}()
}

// SendWithTags metrics to NATS queue waiting for response
//
// Example:
// err = m.SendWithTags(metrics.M{
// 		"counter": i,
// }, metrics.T{
//	    "tag": "sometag",
// }, "metricname")
func (m *conn) SendWithTags(metrics M, tags T, path ...string) {
	go func() {
		if err := m.SendWithTagsAndWait(metrics, tags, path...); err != nil {
			m.errorHandler(err)
		}
	}()
}

// SendWithTagsAndWait metrics to NATS queue waiting for response
//
// Example:
// err = m.SendWithTagsAndWait(metrics.M{
// 		"counter": i,
// }, metrics.T{
//	    "tag": "sometag",
// }, "metricname")
func (m *conn) SendWithTagsAndWait(metrics M, tags T, path ...string) error {
	m.mu.RLock()
	if !m.enabled {
		m.mu.RUnlock()
		return nil
	}
	m.mu.RUnlock()

	if len(metrics) == 0 {
		return nil
	}

	if tags == nil {
		tags = make(T)
	}

	m.mu.RLock()
	tags["hostname"] = m.hostname
	queue := m.queue
	m.mu.RUnlock()

	metricName := append([]string{m.application}, path...)
	buf := format(strings.Join(metricName, ":"), metrics, tags)

	return m.nats.Publish(queue, buf)
}

// SendAndWait metrics to NATS queue waiting for response
//
// Example:
// err = m.SendAndWait(metrics.M{
// 		"counter": i,
// }, "metricname")
func (m *conn) SendAndWait(metrics M, path ...string) error {
	return m.SendWithTagsAndWait(metrics, nil, path...)
}

// Disable disables watcher and disconnects
func (m *conn) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.enabled = false
	m.nats.Close()
}

// Disable disables watcher and disconnects
func Disable() {
	if DefaultConn != nil {
		DefaultConn.Disable()
	}
}

// Watch watches memory, goroutine counter
func (m *conn) Watch(interval time.Duration) error {
	var mem runtime.MemStats

	for {
		m.mu.RLock()
		enabled := m.enabled
		m.mu.RUnlock()

		if !enabled {
			break
		}

		// Getting memory stats
		runtime.ReadMemStats(&mem)
		metric := M{
			"alloc":         mem.Alloc,
			"alloc_objects": mem.HeapObjects,
			"goroutines":    runtime.NumGoroutine(),
			"gc":            mem.LastGC,
			"next_gc":       mem.NextGC,
			"pause_ns":      mem.PauseNs[(mem.NumGC+255)%256],
		}

		err := m.SendAndWait(metric, "gostats")
		if err != nil {
			return err
		}

		time.Sleep(interval)
	}

	return nil
}

// Watch watches memory, goroutine counter
func Watch(interval time.Duration) error {
	if DefaultConn == nil {
		return nil
	}

	return DefaultConn.Watch(interval)
}

func format(name string, metrics M, tags T) []byte {
	buf := bytes.NewBufferString(name)

	if len(tags) > 0 {
		var tagKeys []string
		for k := range tags {
			tagKeys = append(tagKeys, k)
		}
		sort.Strings(tagKeys)

		for _, k := range tagKeys {
			buf.WriteRune(',')
			buf.WriteString(k)
			buf.WriteRune('=')
			buf.WriteString(tags[k])
		}
	}

	buf.WriteRune(' ')
	count := 0

	var metricKeys []string
	for k := range metrics {
		metricKeys = append(metricKeys, k)
	}
	sort.Strings(metricKeys)

	for _, k := range metricKeys {
		if count > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')

		v := metrics[k]
		switch v.(type) {
		case string:
			buf.WriteRune('"')
			buf.WriteString(v.(string))
			buf.WriteRune('"')
		default:
			buf.WriteString(fmt.Sprintf("%v", v))
		}
		count++
	}

	return buf.Bytes()
}
