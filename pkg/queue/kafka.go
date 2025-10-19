package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Key       string                 `json:"key,omitempty"`
	Value     interface{}            `json:"value"`
	Headers   map[string]string      `json:"headers,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type Queue interface {
	Publish(ctx context.Context, topic string, message *Message) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Close() error
}

type MessageHandler func(ctx context.Context, message *Message) error

type KafkaQueue struct {
	brokers []string
	writer  *kafka.Writer
	readers map[string]*kafka.Reader
	config  *KafkaConfig
}

type KafkaConfig struct {
	Brokers       []string
	ClientID      string
	GroupID       string
	BatchSize     int
	BatchTimeout  time.Duration
	RetryAttempts int
	RetryDelay    time.Duration
	Compression   kafka.Compression
	Security      *SecurityConfig
}

type SecurityConfig struct {
	Protocol string
	Username string
	Password string
	CertFile string
	KeyFile  string
	CAFile   string
}

func NewKafkaQueue(config *KafkaConfig) *KafkaQueue {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    config.BatchSize,
		BatchTimeout: config.BatchTimeout,
		Compression:  config.Compression,
	}

	return &KafkaQueue{
		brokers: config.Brokers,
		writer:  writer,
		readers: make(map[string]*kafka.Reader),
		config:  config,
	}
}

func (k *KafkaQueue) Publish(ctx context.Context, topic string, message *Message) error {
	value, err := json.Marshal(message.Value)
	if err != nil {
		return fmt.Errorf("marshal message error: %w", err)
	}

	kafkaMessage := kafka.Message{
		Topic:     topic,
		Key:       []byte(message.Key),
		Value:     value,
		Time:      message.Timestamp,
	}

	for k, v := range message.Headers {
		kafkaMessage.Headers = append(kafkaMessage.Headers, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	return k.writer.WriteMessages(ctx, kafkaMessage)
}

func (k *KafkaQueue) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.brokers,
		Topic:    topic,
		GroupID:  k.config.GroupID,
		MinBytes: 10e3, 
		MaxBytes: 10e6, 
	})

	k.readers[topic] = reader

	go func() {
		defer reader.Close()
		
		for {
			select {
			case <-ctx.Done():
				return
			default:
				kafkaMessage, err := reader.ReadMessage(ctx)
				if err != nil {
					continue
				}

				var value interface{}
				if err := json.Unmarshal(kafkaMessage.Value, &value); err != nil {
					continue
				}

				headers := make(map[string]string)
				for _, h := range kafkaMessage.Headers {
					headers[h.Key] = string(h.Value)
				}

				message := &Message{
					Topic:     kafkaMessage.Topic,
					Key:       string(kafkaMessage.Key),
					Value:     value,
					Headers:   headers,
					Timestamp: kafkaMessage.Time,
				}

				if err := handler(ctx, message); err != nil {
					//TODO: RETRY LOGIC IMPLEMANTATION
					continue
				}
			}
		}
	}()

	return nil
}

func (k *KafkaQueue) Close() error {
	if k.writer != nil {
		k.writer.Close()
	}

	for _, reader := range k.readers {
		reader.Close()
	}

	return nil
}

type ScrapingJob struct {
	ID          string            `json:"id"`
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        string            `json:"body,omitempty"`
	Config      interface{}       `json:"config,omitempty"`
	Priority    int               `json:"priority"`
	Retry       int               `json:"retry"`
	MaxRetries  int               `json:"max_retries"`
	CreatedAt   time.Time         `json:"created_at"`
	ScheduledAt time.Time         `json:"scheduled_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type JobQueue struct {
	queue Queue
	topic string
}

func NewJobQueue(queue Queue, topic string) *JobQueue {
	return &JobQueue{
		queue: queue,
		topic: topic,
	}
}

func (j *JobQueue) Enqueue(ctx context.Context, job *ScrapingJob) error {
	message := &Message{
		ID:        job.ID,
		Topic:     j.topic,
		Key:       job.ID,
		Value:     job,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"priority": job.Priority,
			"retry":    job.Retry,
		},
	}

	return j.queue.Publish(ctx, j.topic, message)
}

func (j *JobQueue) Subscribe(ctx context.Context, handler func(ctx context.Context, job *ScrapingJob) error) error {
	return j.queue.Subscribe(ctx, j.topic, func(ctx context.Context, message *Message) error {
		jobData, err := json.Marshal(message.Value)
		if err != nil {
			return err
		}

		var job ScrapingJob
		if err := json.Unmarshal(jobData, &job); err != nil {
			return err
		}

		return handler(ctx, &job)
	})
}

type PriorityQueue struct {
	queues map[int]*JobQueue 
	topics map[int]string    
}

func NewPriorityQueue(queue Queue, priorities []int) *PriorityQueue {
	pq := &PriorityQueue{
		queues: make(map[int]*JobQueue),
		topics: make(map[int]string),
	}

	for _, priority := range priorities {
		topic := fmt.Sprintf("scraping-jobs-p%d", priority)
		pq.queues[priority] = NewJobQueue(queue, topic)
		pq.topics[priority] = topic
	}

	return pq
}

func (p *PriorityQueue) Enqueue(ctx context.Context, job *ScrapingJob) error {
	queue, exists := p.queues[job.Priority]
	if !exists {
		queue = p.queues[0]
	}

	return queue.Enqueue(ctx, job)
}

func (p *PriorityQueue) Subscribe(ctx context.Context, handler func(ctx context.Context, job *ScrapingJob) error) error {
	for priority := 10; priority >= 0; priority-- {
		if queue, exists := p.queues[priority]; exists {
			if err := queue.Subscribe(ctx, handler); err != nil {
				return err
			}
		}
	}

	return nil
}