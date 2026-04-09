package mq_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
)

func TestPublish_AddsMessageToStream(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)

	err := mq.Init()
	require.NoError(t, err)

	msg := map[string]string{
		"cow_id":       "cow-1",
		"source":       "cow_agent",
		"metric_type":  "temperature",
		"metric_value": "38.5",
		"unit":         "celsius",
		"timestamp":    "1712486400",
	}

	err = mq.Publish(msg)
	require.NoError(t, err)

	info, err := mq.StreamLen()
	require.NoError(t, err)
	assert.Equal(t, int64(1), info)
}

func TestConsume_ReadsPublishedMessage(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)

	err := mq.Init()
	require.NoError(t, err)

	msg := map[string]string{
		"cow_id":       "cow-2",
		"source":       "cow_agent",
		"metric_type":  "heart_rate",
		"metric_value": "72",
		"unit":         "bpm",
		"timestamp":    "1712486400",
	}

	err = mq.Publish(msg)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var received map[string]interface{}
	done := make(chan struct{})

	go mq.Consume(ctx, mq.GroupMetricWriter, "test_consumer", func(id string, values map[string]interface{}) error {
		received = values
		close(done)
		cancel()
		return nil
	})

	select {
	case <-done:
		assert.Equal(t, "cow-2", received["cow_id"])
		assert.Equal(t, "heart_rate", received["metric_type"])
	case <-time.After(6 * time.Second):
		t.Fatal("timeout waiting for consume")
	}
}
