package consumer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server/consumer"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestMetricWriter_WritesToDB(t *testing.T) {
	testhelper.SetupTestDB(t)

	testhelper.WithTx(t, func(tx *gorm.DB) {
		// init mq after WithTx flushes redis so the consumer group exists
		require.NoError(t, mq.Init())

		cow := testhelper.SeedCow(t, tx, "Echo", model.CowStatusInFarm)

		// publish a metric
		err := mq.Publish(map[string]string{
			"cow_id":       cow.ID,
			"source":       "cow_agent",
			"metric_type":  "temperature",
			"metric_value": "38.50",
			"unit":         "celsius",
			"timestamp":    "1712486400",
		})
		require.NoError(t, err)

		// run consumer briefly
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			consumer.StartMetricWriter(ctx)
			close(done)
		}()

		// wait for consumer to process
		time.Sleep(2 * time.Second)
		cancel()
		<-done

		// check metric in db
		var metric model.Metric
		err = pg.DB.Where("cow_id = ?", cow.ID).First(&metric).Error
		require.NoError(t, err)
		assert.Equal(t, cow.ID, metric.CowID)
		assert.Equal(t, model.MetricSourceCowAgent, metric.Source)
		assert.Equal(t, model.MetricTypeTemperature, metric.MetricType)
		assert.InDelta(t, 38.50, metric.MetricValue, 0.01)
		assert.Equal(t, model.MetricUnitCelsius, metric.Unit)
	})
}
