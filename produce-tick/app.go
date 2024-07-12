package main

import (
	"log/slog"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/segmentio/kafka-go"
)

const TOPIC_NAME = "all-ticks"

type AppConfig struct {
	KafkaBootstratServers string `env:"KAFKA_BOOTSTRAP_SERVERS" env-default:"my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"`
}

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func main() {
	ac := AppConfig{}
	cleanenv.ReadEnv(&ac)
	appBuilder := NewAppBuilder(ac)
	appBuilder.WithRandomTickName()
	app := appBuilder.Build()

	createTopicIfNotExists(ac)

	app.Run()
}

type AppBuilder struct {
	AppConfig
	TickName string
}

func NewAppBuilder(ac AppConfig) *AppBuilder {
	return &AppBuilder{AppConfig: ac}
}

func (ab *AppBuilder) WithRandomTickName() *AppBuilder {
	ab.TickName = random()
	return ab
}

// generate a random string of length 3
func random() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 3)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (ab *AppBuilder) Build() *App {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(ab.KafkaBootstratServers),
		Topic:    TOPIC_NAME,
		Balancer: &kafka.Hash{},
	}

	tickServer := NewTickServer(ab.TickName, writer)
	return &App{
		AppConfig:  ab.AppConfig,
		writer:     writer,
		tickServer: tickServer,
	}
}

type App struct {
	AppConfig
	writer     *kafka.Writer
	tickServer *TickServer
}

func (a *App) Run() {
	a.tickServer.Start()
}

func createTopicIfNotExists(ac AppConfig) {
	topic := TOPIC_NAME

	conn, err := kafka.Dial("tcp", ac.KafkaBootstratServers)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
