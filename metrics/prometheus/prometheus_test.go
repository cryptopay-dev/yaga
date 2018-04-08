package prometheus

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	t.Run("no address", func(t *testing.T) {
		if err := os.Setenv("PROMETHEUS", ""); err != nil {
			t.Fatal(err)
		}

		assert.Error(t, Serve())
	})

	t.Run("ok address", func(t *testing.T) {
		server := httptest.NewServer(nil)
		bind := server.Listener.Addr().String()
		server.Close()

		if err := os.Setenv("PROMETHEUS", bind); err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, Serve())
	})
}
