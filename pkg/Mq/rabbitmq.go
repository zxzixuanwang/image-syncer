package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMq struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queueName     string
	exchange      string
	routingKey    string
	l             *logrus.Logger
	ctx           context.Context
	nodeList      []string
	lock          *sync.Mutex
	notifyClose   chan *amqp.Error
	cancel        context.CancelFunc
	content       chan amqp.Delivery
	reconnectFlag chan struct{}
	first         bool
}

// Handle implements MqHandler.
func (r *RabbitMq) Handle(f func(data []byte) error) {
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				r.l.Info("stop handle mq")
				return
			case data := <-r.content:
				err := f(data.Body)
				if err != nil {
					r.l.Error("hande data err", err)
				} else {
					err = data.Ack(true)
					if err != nil {
						r.l.Error("ack err", err)
					}
				}
			}
		}
	}()
}

func (r *RabbitMq) Receive() error {
	var err error
	go func() {
		for {
			var msgs <-chan amqp.Delivery
			r.lock.Lock()
			channel := r.channel
			r.lock.Unlock()
			msgs, err = channel.Consume(r.queueName,
				"",    // consumer
				false, // auto-ack
				false, // exclusive
				false, // no-local
				false, // no-wait
				nil,   // args)
			)
			if err != nil {
				r.l.Error("receive channel err", err)
				return
			}
		outer:
			for {
				select {
				case <-r.ctx.Done():
					r.l.Info("stop listen mq")
					return
				case m := <-msgs:
					r.l.Debug("receiver content", string(m.Body))
					r.content <- m
				case <-r.reconnectFlag:
					r.l.Info("receiver reconnect")
					break outer
				}
			}
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (r *RabbitMq) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.conn.Close(); err != nil {
		return err
	}
	r.cancel()

	return nil
}

func NewRabbitMqSvc(queueName string, url []string, log *logrus.Logger) (*RabbitMq, error) {
	subCtx, cancel := context.WithCancel(context.Background())

	rabbitmq := &RabbitMq{
		queueName:     queueName,
		l:             log,
		ctx:           subCtx,
		nodeList:      url,
		lock:          &sync.Mutex{},
		cancel:        cancel,
		notifyClose:   make(chan *amqp.Error),
		content:       make(chan amqp.Delivery),
		reconnectFlag: make(chan struct{}),
		first:         true,
	}

	//获取connection
	err := rabbitmq.connect()
	if err != nil {
		log.Error("connect rabbitmq err")
		return nil, err
	}

	return rabbitmq, nil
}

func (r *RabbitMq) connect() error {
	var err error

	if r.conn != nil {
		r.conn.Close()
	}
	newConn, err := r.connectCluster()
	if err != nil {
		return err
	}

	newChannel, err := newConn.Channel()
	if err != nil {
		return err
	}

	r.lock.Lock()
	r.conn = newConn
	r.channel = newChannel

	r.notifyClose = make(chan *amqp.Error)
	r.lock.Unlock()

	r.channel.NotifyClose(r.notifyClose)
	if !r.first {
		r.reconnectFlag <- struct{}{}
	}

	go r.listenForReconnect()
	r.first = false

	return nil
}

func (r *RabbitMq) connectCluster() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for _, url := range r.nodeList {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("All node is down!")
}

func (r *RabbitMq) listenForReconnect() {
	clo := <-r.notifyClose
	r.l.Info("closed rabbitmq reason", clo.Reason, "code", clo.Code, "from server", clo.Server)
	for {
		if err := r.connect(); err == nil {
			r.l.Info("connect ok")
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func RabbitMqUrl(username, password, url string) string {
	return fmt.Sprintf("amqp://%s:%s@%s/", username, password, url)
}

type RabbitMqData struct {
	Name string
	Tag  string
}
