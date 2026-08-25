package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/b1g-nguyx/ADRDMS/shared/go/pkg/config"
	"github.com/b1g-nguyx/ADRDMS/shared/go/pkg/health"
	"github.com/b1g-nguyx/ADRDMS/shared/go/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	KafkaBroker string `mapstructure:"KAFKA_BROKER"`
	KafkaTopic  string `mapstructure:"KAFKA_TOPIC"`
}

func main() {
	logger.InitLogger(true)
	log := logger.Get()
	defer logger.Sync()

	cfg := Config{
		Port:        "8082",
		KafkaBroker: "localhost:9092",
		KafkaTopic:  "test-topic",
	}

	if err := config.LoadConfig(".", "config", "yaml", &cfg); err != nil {
		log.Warn("Failed to load config, using defaults", zap.Error(err))
	}

	if cfg.Port == "" {
		cfg.Port = "8082"
	}
	if cfg.KafkaBroker == "" {
		cfg.KafkaBroker = "localhost:9092"
	}
	if cfg.KafkaTopic == "" {
		cfg.KafkaTopic = "test-topic"
	}

	log.Info(fmt.Sprintf("🌊 Stream Service is starting on port %s...", cfg.Port))

	// Kafka connection
	kafkaWriter := &kafka.Writer{
		Addr:  kafka.TCP(cfg.KafkaBroker),
		Topic: cfg.KafkaTopic,
	}
	_ = kafkaWriter // mock usage
	log.Info("Kafka connection initialized (mock/stub)")

	// HTTP Server
	http.HandleFunc("/health", health.HTTPHealthHandler)
	addr := ":" + strings.TrimPrefix(cfg.Port, ":")
	log.Info("Starting HTTP server on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("Failed to start server", zap.Error(err))
	}
}
