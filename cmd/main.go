package main

import (
	"context"
	"fmt"
	"okto-nft-ms/api"
	"okto-nft-ms/initializer"
	"okto-nft-ms/internal/appconfig"
	"okto-nft-ms/pkg/cache"
	cacheClient "okto-nft-ms/pkg/cache/client"
	"okto-nft-ms/pkg/cache/drivers"
	kafkaclient "okto-nft-ms/pkg/kafka/client"
	"okto-nft-ms/pkg/logger"
	"okto-nft-ms/pkg/queue"
	"okto-nft-ms/pkg/storage"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"strconv"
)

var (
	cachewstring        *cache.CacheWrapper[string, string]
	cachewint           *cache.CacheWrapper[string, int]
	kafkaProducerClient *kafkaclient.ProducerClient
	kafkaConsumerClient *kafkaclient.ConsumerClient
	sqsqw               *queue.QueueWrapper
)

func init() {
	//if env := os.Getenv("ENV"); env == "DEVELOPMENT" {
	initializer.LoadEnvVariables()
	//	}

	appconfig.SetEnvVariables()
	logger.SetAppName("OKTO-NFT-MS")
	logger.Info(fmt.Sprintf("Current Env is %s", appconfig.ENV))
	logger.Info(fmt.Sprintf("Is Production %t", appconfig.IS_PRODUCTION))
}

func main() {
	tracer.Start(
		tracer.WithEnv(appconfig.DD_ENV),
		tracer.WithService(appconfig.DD_SERVICE),
		tracer.WithRuntimeMetrics(),
		tracer.WithAnalytics(true),
		tracer.WithService("OKTO-NFT-MS"),
		tracer.WithTraceEnabled(true),
		tracer.WithDogstatsdAddress(appconfig.DD_AGENT_HOST),
	)
	defer tracer.Stop()
	initcache()
	initProducer()
	initConsumer()
	databaseUrl := appconfig.DATABASE_URL

	if len(databaseUrl) == 0 {
		panic("env variable DATABASE_URL is missing")
	}

	db := &storage.Store{}
	db.InitPostgresClient(databaseUrl)

	sqsInterface := initializer.InitSQS()
	repositories := initializer.InitRepositories(db)
	services := initializer.InitServices(repositories, cachewstring, cachewint, kafkaProducerClient, kafkaConsumerClient, sqsInterface)

	api.InitServer(services)
}

func initProducer() {
	kafkaProducerClient, _ = kafkaclient.NewProducerClient()
}

func initConsumer() {
	kafkaConsumerClient, _ = kafkaclient.NewConsumerClient(appconfig.KAFKA_HOST, appconfig.KAFKA_GROUP_ID, appconfig.FM_CONSUMER_POLL_INTERVAL)
}

func initcache() {
	maxRedirects, err := strconv.Atoi(appconfig.MAX_REDIRECTS)
	if err != nil {
		logger.Error("MAX_REDIRECTS ENV variable not present : ", err)
		maxRedirects = 3
	}
	redisPoolSize, err := strconv.Atoi(appconfig.REDIS_POOL_SIZE)
	if err != nil {
		logger.Error("REDIS_POOL_SIZE ENV variable not present : ", err)
		redisPoolSize = 32
	}
	config := cacheClient.Config{
		Addr:         appconfig.REDIS_URL,
		MaxRedirects: maxRedirects,
		PoolSize:     redisPoolSize,
	}
	if redis, _, err := cacheClient.CreateNewRedisClient(config); err != nil {
		panic(err)
	} else {
		cachewstring = cache.NewCacheWrapper[string, string]("redis", drivers.NewRedisClientString(redis, context.Background()))
		cachewint = cache.NewCacheWrapper[string, int]("redis", drivers.NewRedisClientInt(redis, context.Background()))
	}
}
