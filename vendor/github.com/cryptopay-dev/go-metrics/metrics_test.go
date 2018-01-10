package metrics

import (
	"math/rand"
	"os"
	"testing"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

func generateMetric() M {
	return M{
		"uint64":  uint64(rand.Int()),
		"uint32":  uint32(rand.Int()),
		"unit16":  uint16(rand.Int()),
		"int64":   int64(rand.Int()),
		"int32":   int32(rand.Int()),
		"int":     rand.Int(),
		"float64": rand.Float64(),
		"float32": rand.Float32(),
		"bool":    true,
		"string":  "username",
	}
}

func generateTags() T {
	return T{
		"username": "user" + string(rand.Intn(100)),
	}
}

var natsURL string

func TestMain(m *testing.M) {
	natsURL = os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	m.Run()
}

func TestFormat(t *testing.T) {
	var f []byte

	f = format("test", M{"m1": 1, "m2": 2, "m3": 3.02, "m4": "string", "m5": true}, T{"t1": "tag-one", "t2": "tag-two"})
	assert.Equal(t, `test,t1=tag-one,t2=tag-two m1=1,m2=2,m3=3.02,m4="string",m5=true`, string(f))

	f = format("test", M{"m1": 1, "m2": 2, "m3": 3.02, "m4": "string"}, nil)
	assert.Equal(t, `test m1=1,m2=2,m3=3.02,m4="string"`, string(f))
}

func TestMetrics(t *testing.T) {
	t.Run("No setup", func(t *testing.T) {
		var err error

		Disable()
		Send(M{"field": 1}, "test")
		SendWithTags(M{"field": 1}, T{"tag": "string"}, "test")

		err = SendAndWait(M{"field": 1}, "test")
		assert.NoError(t, err)
		err = SendWithTagsAndWait(M{"field": 1}, T{"tag": "string"}, "test")
		assert.NoError(t, err)

		err = Watch(time.Second)
		assert.NoError(t, err)
	})

	t.Run("Unknown server", func(t *testing.T) {
		metrics, err := New("1.1.1.1:1111", "metrics", "localhost")

		assert.Error(t, err)
		assert.True(t, metrics == nil)
	})

	t.Run("Unknown server setup", func(t *testing.T) {
		err := Setup("1.1.1.1:1111", "metrics", "localhost")

		assert.Error(t, err)
		assert.True(t, DefaultConn == nil)
	})

	t.Run("Empty application", func(t *testing.T) {
		metrics, err := New("1.1.1.1:1111", "", "")

		assert.Error(t, err)
		assert.True(t, metrics == nil)
	})

	t.Run("Empty hostname", func(t *testing.T) {
		metrics, err := New("1.1.1.1:1111", "metrics", "")

		assert.Error(t, err)
		assert.True(t, metrics == nil)
	})

	t.Run("Disabled metrics", func(t *testing.T) {
		metrics, err := New("", "metrics", "localhost")

		assert.NoError(t, err)
		assert.True(t, metrics != nil)

		err = metrics.SendAndWait(generateMetric(), "test")
		assert.NoError(t, err)
	})

	t.Run("Empty metrics", func(t *testing.T) {
		metrics, err := New(natsURL, "metrics", "localhost")

		if assert.NoError(t, err) {
			assert.True(t, metrics != nil)

			err = metrics.SendAndWait(map[string]interface{}{}, "test")
			assert.NoError(t, err)
		}
	})

	t.Run("Connection", func(t *testing.T) {
		metrics, err := New(natsURL, "metrics", "localhost")

		if assert.NoError(t, err) {
			assert.True(t, metrics != nil)

			t.Run("Synchronous send", func(t *testing.T) {
				var err error

				err = metrics.SendAndWait(generateMetric(), "test")
				assert.NoError(t, err)

				err = metrics.SendWithTagsAndWait(generateMetric(), generateTags(), "test")
				assert.NoError(t, err)
			})

			t.Run("Asynchronous send", func(t *testing.T) {
				metrics.SendWithTags(generateMetric(), generateTags(), "test")
				metrics.Send(generateMetric(), "test")
			})
		}
	})

	t.Run("Default connection", func(t *testing.T) {
		err := Setup(natsURL, "metrics", "localhost")

		if assert.NoError(t, err) {
			assert.True(t, DefaultConn != nil)

			t.Run("Synchronous send", func(t *testing.T) {
				var err error

				err = SendAndWait(generateMetric(), "test")
				assert.NoError(t, err)

				err = SendWithTagsAndWait(generateMetric(), generateTags(), "test")
				assert.NoError(t, err)
			})

			t.Run("Asynchronous send", func(t *testing.T) {
				SendWithTags(generateMetric(), generateTags(), "test")
				Send(generateMetric(), "test")
			})
		}
	})

	t.Run("Auto sending", func(t *testing.T) {
		metrics, err := New(natsURL, "metrics", "localhost")

		if assert.NoError(t, err) {
			assert.True(t, metrics != nil)

			done := make(chan bool, 1)
			go func() {
				metrics.Watch(time.Millisecond * 100)
				done <- true
			}()

			time.Sleep(time.Millisecond * 500)
			metrics.Disable()

			assert.True(t, <-done)
		}
	})

	t.Run("Auto sending default connection", func(t *testing.T) {
		err := Setup(natsURL, "metrics", "localhost")

		if assert.NoError(t, err) {
			assert.True(t, DefaultConn != nil)

			done := make(chan bool, 1)
			go func() {
				Watch(time.Millisecond * 100)
				done <- true
			}()

			time.Sleep(time.Millisecond * 500)
			Disable()

			assert.True(t, <-done)
		}
	})
}
