package tests

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTime(t *testing.T) {
	t.Skip()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"http://localhost:8080/time",
		nil,
	)
	require.NoError(t, err)

	t.Skip()

	client := http.Client{}

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}

	wg.Wait()
}
