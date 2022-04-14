package queues

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/nats.go"

	"github.com/ThreeDotsLabs/watermill-jetstream/pkg/jetstream"
	"github.com/google/uuid"
	"github.com/lejenome/lro/services/process-executor/lib/process"
)

const TOPIC_JOB_SCHEDULED = "process/job/scheduled"

type natsQueue struct {
	sync.RWMutex
	jobs     process.JobCache
	closed   bool
	sub      *jetstream.Subscriber
	pub      *jetstream.Publisher
	messages <-chan *message.Message
}

var _ process.Queue = (*natsQueue)(nil)

func NatsSubscriber(ctx context.Context, url string, clusterID string, jc process.JobCache) process.Queue {
	subscriber, err := jetstream.NewSubscriber(
		jetstream.SubscriberConfig{
			URL:              url,
			QueueGroup:       "lro-jobs-queue",
			DurableName:      "lro-jobs-queue",
			SubscribersCount: 4, // how many goroutines should consume messages
			NatsOptions: []nats.Option{
				nats.RetryOnFailedConnect(true),
				nats.Timeout(30 * time.Second),
				nats.ReconnectWait(1 * time.Second),
			},
			SubscribeOptions: []nats.SubOpt{
				nats.DeliverAll(),
				nats.AckExplicit(),
			},
			CloseTimeout:   time.Minute,
			AckWaitTimeout: time.Second * 30,
			AutoProvision:  true,
			AckSync:        true,

			// ClusterID: clusterID,
			// ClientID: "process-executor",
			// StanOptions: []stan.Option{
			// 	stan.NatsURL(url),
			// },
			Unmarshaler: &jetstream.ProtoMarshaler{}, // &wmpb.NATSMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err.Error())
	}
	messages, err := subscriber.Subscribe(ctx, TOPIC_JOB_SCHEDULED)
	if err != nil {
		panic(err.Error())
	}
	return &natsQueue{
		jobs:     jc,
		sub:      subscriber,
		messages: messages,
		pub:      nil,
	}
}
func NatsPublisher(url string, clusterID string, jc process.JobCache) process.Queue {
	pub, err := jetstream.NewPublisher(
		jetstream.PublisherConfig{
			URL: url,
			// QueueGroup:       "lro-jobs-queue",
			// DurableName:      "lro-jobs-queue",
			NatsOptions: []nats.Option{
				nats.RetryOnFailedConnect(true),
				nats.Timeout(30 * time.Second),
				nats.ReconnectWait(1 * time.Second),
			},
			// PublishOptions: []nats.PubOpt{},
			AutoProvision: true,
			TrackMsgId:    true,

			// ClusterID: clusterID,
			// ClientID: "process-executor",
			// StanOptions: []stan.Option{
			// 	stan.NatsURL(url),
			// },
			Marshaler: &jetstream.ProtoMarshaler{}, // &wmpb.NATSMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err.Error())
	}
	return &natsQueue{
		jobs:     jc,
		pub:      pub,
		messages: nil,
		sub:      nil,
	}
}

func (q *natsQueue) Len() (int, error) {
	return 1, errors.New("Not supported")
}
func (q *natsQueue) Add(job *process.Job) error {
	if err := q.jobs.Add(job); err != nil {
		return err
	}
	return q.AddJobId(job.ID)
}
func (q *natsQueue) SafeAdd(job *process.Job) error {
	q.jobs.Add(job)
	return q.AddJobId(job.ID)
}
func (q *natsQueue) AddJobId(jobid uuid.UUID) error {
	if q.pub == nil {
		panic("Queue type is not a publisher")
	}
	if q.closed {
		return errors.New("Queue closed")
	}
	msg := message.NewMessage(watermill.NewUUID(), jobid[:])
	err := q.pub.Publish(TOPIC_JOB_SCHEDULED, msg)
	return err
}
func (q *natsQueue) GetJobId() (uuid.UUID, error) {
	if q.sub == nil {
		panic("Queue type is not a subscriber")
	}
	if q.closed {
		return uuid.Nil, errors.New("Queue closed")
	}
	msg, ok := <-q.messages
	if !ok {
		return uuid.Nil, errors.New("Error retriving new job id from the queue")
	} else if job, err := uuid.FromBytes(msg.Payload); err != nil {
		return uuid.Nil, err
	} else {
		msg.Ack()
		return job, nil
	}
}
func (q *natsQueue) Get() (*process.Job, error) {
	id, err := q.GetJobId()
	if err != nil {
		return nil, err
	}
	return q.jobs.Get(id)
}

func (q *natsQueue) Purge() error {
	if q.closed {
		return errors.New("Queue closed")
	}
	return errors.New("Not supported")
}
func (q *natsQueue) Close() error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	q.closed = true
	return nil
}
