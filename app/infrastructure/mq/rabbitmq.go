package mq

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"

	"jade-mes/app/infrastructure/log"
	"jade-mes/app/infrastructure/util"
)

type rabbitmqSetting struct {
	HostName              string
	Port                  int
	UserName              string
	Password              string
	ExchangeName          string
	ExchangeType          string
	Prefetch              int
	PrefetchForDataReport int
	QueueRoutings         map[string][]string
	QueueFallbacks        map[string][]string
}

// RabbitMQSession ...
type RabbitMQSession struct {
	sync.Mutex

	addr            string
	name            string
	exchangeName    string
	exchangeType    string
	queueBindings   map[string][]string
	logger          *log.Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	prefetch        int
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	notifyReturn    chan amqp.Return
	isReady         bool
	isConsumer      bool
	consumers       map[string]func(<-chan amqp.Delivery, error)
	consumerReady   bool
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second

	setupTimeout = 30 * time.Second

	maxRetry = 3
)

var (
	errNotReady            = errors.New("session is not ready yet")
	errRoutingKeyNotExists = errors.New("routing key is not inited")
	errNotConnected        = errors.New("not connected to a server")
	errAlreadyClosed       = errors.New("already closed: not connected to the server")
	errShutdown            = errors.New("session is shutting down")
	mqConfig               rabbitmqSetting
	consumerSeq            uint64
)

// GetRabbitMQProducer ...
func GetRabbitMQProducer() *RabbitMQSession {
	return rabbitMQProducer
}

// GetRabbitMQConsumer ...
func GetRabbitMQConsumer() *RabbitMQSession {
	return rabbitMQConsumer
}

// GetRabbitMQConsumers ...
func GetRabbitMQConsumers() []*RabbitMQSession {
	return []*RabbitMQSession{rabbitMQConsumer}
}

func initRabbitMQSession(config *viper.Viper, isConsumer bool, forDataReport bool) *RabbitMQSession {
	config.UnmarshalKey("rabbitmq", &mqConfig)

	connStr := getRabbitMQConnectionString(&mqConfig)

	prefetch := mqConfig.Prefetch
	if forDataReport {
		prefetch = mqConfig.PrefetchForDataReport
	}
	queueRoutings := mqConfig.QueueRoutings
	if isConsumer {
		queueRoutings = make(map[string][]string)
	}

	session := NewRabbitMQSession(mqConfig.ExchangeName, mqConfig.ExchangeType, connStr, queueRoutings, prefetch)
	session.isConsumer = isConsumer
	session.addr = connStr

	return session
}

func getRabbitMQConnectionString(config *rabbitmqSetting) string {
	authStr := ""
	if config.UserName != "" {
		authStr = fmt.Sprintf("%v:%v", config.UserName, config.Password)
	}

	return fmt.Sprintf("amqp://%v@%v:%v/", authStr, config.HostName, config.Port)
}

// NewRabbitMQSession creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewRabbitMQSession(exchangeName, exchangeType, addr string, queueBindings map[string][]string, prefetch int) *RabbitMQSession {
	session := RabbitMQSession{
		logger:        log.GetLogger().With(log.String("exchangeName", exchangeName)),
		exchangeName:  exchangeName,
		exchangeType:  exchangeType,
		queueBindings: queueBindings,
		done:          make(chan bool),
		consumers:     make(map[string]func(<-chan amqp.Delivery, error)),
		prefetch:      prefetch,
	}
	//go session.handleReconnect(addr)
	return &session
}

func newConsumerTag() string {
	tagPrefix := "ctag-"
	cmdName := os.Args[0]
	hostName, _ := os.Hostname()
	tagInfix := cmdName + "@" + hostName
	tagSuffix := "-" + strconv.FormatUint(atomic.AddUint64(&consumerSeq, 1), 10)

	return tagPrefix + tagInfix + tagSuffix
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (session *RabbitMQSession) handleReconnect(addr string) {
	for {
		session.isReady = false
		session.logger.Debug("Attempting to connect")

		conn, err := session.connect(addr)

		if err != nil {
			session.logger.Warn("Failed to connect. Retrying...", log.String("mqConnectionString", addr), log.Err(err))

			select {
			case <-session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := session.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (session *RabbitMQSession) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	session.changeConnection(conn)
	session.logger.Debug("Connected!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (session *RabbitMQSession) handleReInit(conn *amqp.Connection) bool {
	for {
		session.isReady = false

		err := session.init(conn)

		if err != nil {
			session.logger.Error("Failed to initialize channel. Retrying...", log.Err(err))

			select {
			case <-session.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-session.done:
			return true
		case err := <-session.notifyConnClose:
			session.logger.Error("Connection closed. Reconnecting...", log.Err(err))
			return false
		case err := <-session.notifyChanClose:
			session.logger.Error("Channel closed. Re-running init...", log.Err(err))
		}
	}
}

// init will initialize channel & declare queue
func (session *RabbitMQSession) init(conn *amqp.Connection) error {
	session.logger.Info("initing mq channel & queue...")
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	if err = ch.Confirm(false); err != nil {
		return err
	}

	// exchange declare
	if session.exchangeName == "" {
		return errors.New("exchange name is required")
	}

	exchangeType := session.exchangeType
	if exchangeType == "" {
		exchangeType = amqp.ExchangeTopic
	}

	if err := ch.ExchangeDeclare(session.exchangeName, exchangeType, true, false, false, false, nil); err != nil {
		return err
	}

	// queue declare
	if len(session.queueBindings) > 0 {
		g, _ := errgroup.WithContext(context.Background())

		for routingKey, queueNames := range session.queueBindings {
			for _, queueName := range queueNames {
				routingKey := routingKey
				queueName := queueName

				g.Go(func() error {
					_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
					if err != nil {
						return err
					}

					return ch.QueueBind(queueName, routingKey, session.exchangeName, false, nil)
				})
			}
		}

		if err := g.Wait(); err != nil {
			return err
		}
	}

	if session.prefetch > 0 {
		if err := ch.Qos(session.prefetch, 0, false); err != nil {
			return err
		}
	}

	session.changeChannel(ch)
	session.isReady = true

	// only for consumer
	if session.isConsumer {
		session.applyConsumers()
	}

	session.logger.Info("rabbitmq chan & queue init done.")

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (session *RabbitMQSession) changeConnection(connection *amqp.Connection) {
	session.connection = connection
	session.notifyConnClose = make(chan *amqp.Error)
	session.connection.NotifyClose(session.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (session *RabbitMQSession) changeChannel(channel *amqp.Channel) {
	session.channel = channel
	session.notifyChanClose = make(chan *amqp.Error)
	session.notifyConfirm = make(chan amqp.Confirmation, 1)
	session.notifyReturn = make(chan amqp.Return, 1)

	session.channel.NotifyClose(session.notifyChanClose)
	session.channel.NotifyPublish(session.notifyConfirm)
	session.channel.NotifyReturn(session.notifyReturn)
}

func (session *RabbitMQSession) waitForReady() error {
	timeout := time.After(setupTimeout)

	for {
		if session.isReady {
			return nil
		}

		select {
		case <-timeout:
			return errors.New("time out while wait for ready")
		case <-time.After(resendDelay):
		}
	}
}

// Push will push data onto the queue, and wait for a confirm.
// If no confirms are received until within the resendTimeout,
// it continuously re-sends messages until a confirm is received.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (session *RabbitMQSession) Push(data []byte, routingKey string) error {
	logger := session.logger.With(log.String("routingKey", routingKey), log.String("rabbitMQData", string(data)))

	attempts := 0
	for {
		err := session.UnsafePush(data, routingKey)
		mutex.Lock()
		attempts++
		mutex.Unlock()

		if err != nil {
			logger.Error("Push failed. Retrying...", log.Err(err))
			select {
			case <-session.done:
				return errShutdown
			case <-time.After(resendDelay):
			}

			mutex.Lock()
			if attempts > maxRetry {
				mutex.Unlock()
				return err
			}
			mutex.Unlock()

			continue
		}

		select {
		case confirm := <-session.notifyConfirm:
			if confirm.Ack {
				logger.Debug("Push confirmed!")
				return nil
			}
			logger.Warn("Push confirm not ack")
		case returnPayload := <-session.notifyReturn:
			logger.Warn("Push returned.", log.Reflect("returnPayload", returnPayload))
			return nil
		case <-time.After(resendDelay):
		}

		session.logger.Warn("Push didn't confirm. Retrying...")
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// recieve the message.
func (session *RabbitMQSession) UnsafePush(data []byte, routingKey string) error {
	if !session.isReady {
		return errNotConnected
	}

	messageID, err := util.NewGUID()
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		MessageId:    messageID,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         data,
	}

	return session.channel.Publish(
		session.exchangeName,
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		message,
	)
}

func (session *RabbitMQSession) applyConsumers() {
	session.Lock()
	defer session.Unlock()

	for queueName, handler := range session.consumers {
		queueName := queueName
		handler := handler
		session.logger.Info("applying consumer", log.String("queueName", queueName), log.String("handlerName", util.GetFunctionName(handler)))
		go handler(session.stream(queueName))
	}

	if len(session.consumers) > 0 {
		session.consumerReady = true
	}
}

// RegisterConsumeHandlers ...
func (session *RabbitMQSession) RegisterConsumeHandlers(handlers map[string]func(<-chan amqp.Delivery, error)) {
	session.Lock()
	defer session.Unlock()

	for queueName, handler := range handlers {
		queueName := queueName
		handler := handler
		session.consumers[queueName] = handler
	}
}

// stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (session *RabbitMQSession) stream(queueName string) (<-chan amqp.Delivery, error) {
	logger := session.logger.With(log.String("queueName", queueName))
	if !session.isReady {
		logger.Warn("session is not ready, waiting ...")
		if err := session.waitForReady(); err != nil {
			return nil, err
		}
	}

	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := session.channel.Consume(
				queueName,
				newConsumerTag(), // Consumer
				false,            // Auto-Ack
				false,            // Exclusive
				false,            // No-local
				false,            // No-Wait
				nil,              // Args
			)
			if err != nil {
				logger.Error("rabbitmq consume failed", log.Err(err))

				select {
				case <-session.done:
					return
				case <-time.After(reconnectDelay):
				}
				continue
			}

			for msg := range d {
				deliveries <- msg
			}
		}
	}()

	return deliveries, nil

	// return session.channel.Consume(
	// 	queueName,
	// 	"",    // Consumer
	// 	false, // Auto-Ack
	// 	false, // Exclusive
	// 	false, // No-local
	// 	false, // No-Wait
	// 	nil,   // Args
	// )
}

// Close will cleanly shutdown the channel and connection.
func (session *RabbitMQSession) Close() error {
	if !session.isReady {
		return errAlreadyClosed
	}
	err := session.channel.Close()
	if err != nil {
		return err
	}
	err = session.connection.Close()
	if err != nil {
		return err
	}
	close(session.done)
	session.isReady = false
	return nil
}
