package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/VadimGossip/platform_common/pkg/db/postgres"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/pg"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/transaction"
	"github.com/VadimGossip/platform_common/pkg/db/redis"
	"github.com/VadimGossip/platform_common/pkg/db/redis/rdb"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka"
	kafkaConsumer "github.com/VadimGossip/concoleChat-auth/internal/client/kafka/consumer"
	kafkaProducer "github.com/VadimGossip/concoleChat-auth/internal/client/kafka/producer"
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	db_cfg "github.com/VadimGossip/concoleChat-auth/internal/config/db"
	kafka_cfg "github.com/VadimGossip/concoleChat-auth/internal/config/kafka"
	server_cfg "github.com/VadimGossip/concoleChat-auth/internal/config/server"
	service_cfg "github.com/VadimGossip/concoleChat-auth/internal/config/service"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	auditRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/audit"
	userRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg"
	userCacheRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/redis"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	auditService "github.com/VadimGossip/concoleChat-auth/internal/service/audit"
	consumerService "github.com/VadimGossip/concoleChat-auth/internal/service/consumer"
	userConsumerService "github.com/VadimGossip/concoleChat-auth/internal/service/consumer/user"
	producerService "github.com/VadimGossip/concoleChat-auth/internal/service/producer"
	userProducerService "github.com/VadimGossip/concoleChat-auth/internal/service/producer/user"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
	userCacheService "github.com/VadimGossip/concoleChat-auth/internal/service/usercache"
)

type serviceProvider struct {
	grpcConfig             config.GRPCConfig
	httpConfig             config.HTTPConfig
	swaggerConfig          config.SwaggerConfig
	pgConfig               config.PgConfig
	redisConfig            config.RedisConfig
	userKafkaServiceConfig config.UserKafkaServiceConfig
	userCacheConfig        config.UserCacheConfig
	kafkaProducerConfig    config.KafkaProducerConfig
	kafkaConsumerConfig    config.KafkaConsumerConfig

	pgDbClient    postgres.Client
	txManager     postgres.TxManager
	redisDbClient redis.Client
	auditRepo     repository.AuditRepository
	userCacheRepo repository.UserCacheRepository
	userRepo      repository.UserRepository

	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
	consumer             kafka.Consumer

	saramaSyncProducer sarama.SyncProducer
	producer           kafka.Producer

	auditService        service.AuditService
	userCacheService    service.UserCacheService
	userService         service.UserService
	userProducerService producerService.UserProducerService
	userConsumerService consumerService.UserConsumerService

	userImpl *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := server_cfg.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpcConfig: %s", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := server_cfg.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get httpConfig: %s", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := server_cfg.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swaggerConfig: %s", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PGConfig() config.PgConfig {
	if s.pgConfig == nil {
		cfg, err := db_cfg.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pgConfig: %s", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := db_cfg.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redisConfig: %s", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) UserKafkaServiceConfig() config.UserKafkaServiceConfig {
	if s.userKafkaServiceConfig == nil {
		cfg, err := service_cfg.NewUserKafkaServiceConfig()
		if err != nil {
			log.Fatalf("failed to get userKafkaServiceConfig: %s", err)
		}

		s.userKafkaServiceConfig = cfg
	}

	return s.userKafkaServiceConfig
}

func (s *serviceProvider) UserCacheConfig() config.UserCacheConfig {
	if s.userCacheConfig == nil {
		cfg, err := service_cfg.NewUserCacheConfig()
		if err != nil {
			log.Fatalf("failed to get userCacheConfig: %s", err)
		}

		s.userCacheConfig = cfg
	}

	return s.userCacheConfig
}

func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		cfg, err := kafka_cfg.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafkaProducerConfig: %s", err)
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := kafka_cfg.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafkaConsumerConfig: %s", err)
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) PgDbClient(ctx context.Context) postgres.Client {
	if s.pgDbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logrus.Fatalf("failed to create db client: %s", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err)
		}
		closer.Add(cl.Close)
		s.pgDbClient = cl
	}

	return s.pgDbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) postgres.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.PgDbClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) RedisDbClient(ctx context.Context) redis.Client {
	if s.redisDbClient == nil {
		cl := rdb.New(&rdb.ClientOptions{
			Addr:         s.RedisConfig().Address(),
			Username:     s.RedisConfig().Username(),
			Password:     s.RedisConfig().Password(),
			DB:           s.RedisConfig().DB(),
			ReadTimeout:  s.RedisConfig().ReadTimeoutSec(),
			WriteTimeout: s.RedisConfig().WriteTimeoutSec(),
		})

		if err := cl.DB().Ping(ctx); err != nil {
			log.Fatalf("kdb ping error: %s", err)
		}

		closer.Add(cl.Close)
		s.redisDbClient = cl
	}

	return s.redisDbClient
}

func (s *serviceProvider) AuditRepository(ctx context.Context) repository.AuditRepository {
	if s.auditRepo == nil {
		s.auditRepo = auditRepo.NewRepository(s.PgDbClient(ctx))
	}
	return s.auditRepo
}

func (s *serviceProvider) AuditService(ctx context.Context) service.AuditService {
	if s.auditService == nil {
		s.auditService = auditService.NewService(s.AuditRepository(ctx))
	}
	return s.auditService
}

func (s *serviceProvider) UserCacheRepository(ctx context.Context) repository.UserCacheRepository {
	if s.userCacheRepo == nil {
		s.userCacheRepo = userCacheRepo.NewRepository(s.RedisDbClient(ctx))
	}
	return s.userCacheRepo
}

func (s *serviceProvider) UserCacheService(ctx context.Context) service.UserCacheService {
	if s.userCacheService == nil {
		s.userCacheService = userCacheService.NewService(s.UserCacheConfig(), s.UserCacheRepository(ctx))
	}
	return s.userCacheService
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepo.NewRepository(s.PgDbClient(ctx))
	}
	return s.userRepo
}

func (s *serviceProvider) SaramaSyncProducer() sarama.SyncProducer {
	if s.saramaSyncProducer == nil {
		producer, err := sarama.NewSyncProducer(s.KafkaProducerConfig().Brokers(), s.kafkaProducerConfig.Config())
		if err != nil {
			log.Fatalf("failed to create sarama sync producer: %v", err)
		}

		s.saramaSyncProducer = producer
	}

	return s.saramaSyncProducer
}

func (s *serviceProvider) Producer() kafka.Producer {
	if s.producer == nil {
		s.producer = kafkaProducer.NewProducer(s.SaramaSyncProducer())
	}

	return s.producer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.UserCacheService(ctx), s.AuditService(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserProducerService() producerService.UserProducerService {
	if s.userProducerService == nil {
		s.userProducerService = userProducerService.NewService(s.UserKafkaServiceConfig(), s.Producer())
	}

	return s.userProducerService
}

func (s *serviceProvider) UserConsumerService(ctx context.Context) consumerService.UserConsumerService {
	if s.userConsumerService == nil {
		s.userConsumerService = userConsumerService.NewService(s.UserKafkaServiceConfig(), s.Consumer(), s.UserService(ctx))
	}

	return s.userConsumerService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx), s.UserProducerService())
	}

	return s.userImpl
}
