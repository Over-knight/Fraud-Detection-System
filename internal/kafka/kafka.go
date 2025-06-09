package kafka

import (
    "context"
    "encoding/json"
    "log"

    "frauddetector/pkg/model"
    "github.com/go-redis/redis/v8"
    "github.com/segmentio/kafka-go"
)

// Config holds Kafka settings
type Config struct {
    Brokers     []string
    GroupID     string
    Topic       string
    AlertsTopic string
}

// Consumer represents the Kafka consumer
type Consumer struct {
    reader *kafka.Reader
    writer *kafka.Writer
    rdb *redis.Client
	ml *ml.Client 
}

// NewConsumer initializes Kafka reader and writer
type NewConsumerFunc func(Config, *redis.Client) *Consumer

// NewConsumer creates a new Kafka Consumer
type NewConsumerConstructor struct{}

func NewConsumer(cfg Config, rdb *redis.Client, mlClient *ml.Client) *Consumer {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: cfg.Brokers,
        GroupID: cfg.GroupID,
        Topic:   cfg.Topic,
    })
    writer := &kafka.Writer{
        Addr:     kafka.TCP(cfg.Brokers...),
        Topic:    cfg.AlertsTopic,
        Balancer: &kafka.LeastBytes{},
    }
    return &Consumer{reader: reader, 
		writer: writer, 
		rdb: rdb,
		ml: mlClient,}
}

// Run starts consuming messages
func (c *Consumer) Run(ctx context.Context) error {
    for {
        m, err := c.reader.FetchMessage(ctx)
        if err != nil {
            return err
        }

        var tx model.Transaction
        if err := json.Unmarshal(m.Value, &tx); err != nil {
            log.Printf("invalid message payload: %v", err)
            continue
        }

        // Simple rule: flag fraud if amount > 1000
        verdict. err := c.ml.Predict(tx)
		if err != nil {
			log.Printf("ML prediction error: %v", err)
			verdict = "error"
        }

        // Store in Redis
        if err := c.rdb.Set(ctx, tx.ID, verdict, 0).Err(); err != nil {
			log.Printf("failed to store transaction in Redis: %v", err)
		}

        // Publish verdict to alerts topic
        alert := map[string]string{"id": tx.ID, "verdict": verdict}
        payload, _ := json.Marshal(alert)
        if err := c.writer.WriteMessages(ctx, kafka.Message{Value: payload}); err != nil {
            log.Printf("failed to write alert: %v", err)
        }

        // Commit offset
        if err := c.reader.CommitMessages(ctx, m); err != nil {
            log.Printf("failed to commit message: %v", err)
        }
    }
}