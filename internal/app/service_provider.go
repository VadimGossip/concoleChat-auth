package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/VadimGossip/concoleChat-auth/internal/api/access"
	"github.com/VadimGossip/concoleChat-auth/internal/api/auth"
	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka"
	kafkaConsumer "github.com/VadimGossip/concoleChat-auth/internal/client/kafka/consumer"
	kafkaProducer "github.com/VadimGossip/concoleChat-auth/internal/client/kafka/producer"
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	dbCfg "github.com/VadimGossip/concoleChat-auth/internal/config/db"
	kafkaCfg "github.com/VadimGossip/concoleChat-auth/internal/config/kafka"
	serverCfg "github.com/VadimGossip/concoleChat-auth/internal/config/server"
	serviceCfg "github.com/VadimGossip/concoleChat-auth/internal/config/service"
	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	accessRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/access"
	auditRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/audit"
	userRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg"
	userCacheRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/redis"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	accessService "github.com/VadimGossip/concoleChat-auth/internal/service/access"
	auditService "github.com/VadimGossip/concoleChat-auth/internal/service/audit"
	authService "github.com/VadimGossip/concoleChat-auth/internal/service/auth"
	consumerService "github.com/VadimGossip/concoleChat-auth/internal/service/consumer"
	userConsumerService "github.com/VadimGossip/concoleChat-auth/internal/service/consumer/user"
	passwordService "github.com/VadimGossip/concoleChat-auth/internal/service/password"
	producerService "github.com/VadimGossip/concoleChat-auth/internal/service/producer"
	userProducerService "github.com/VadimGossip/concoleChat-auth/internal/service/producer/user"
	tokenService "github.com/VadimGossip/concoleChat-auth/internal/service/token"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
	userCacheService "github.com/VadimGossip/concoleChat-auth/internal/service/usercache"
	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/VadimGossip/platform_common/pkg/db/postgres"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/pg"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/transaction"
	"github.com/VadimGossip/platform_common/pkg/db/redis"
	"github.com/VadimGossip/platform_common/pkg/db/redis/rdb"
)

type serviceProvider struct {
	grpcConfig             config.GRPCConfig
	httpConfig             config.HTTPConfig
	swaggerConfig          config.SwaggerConfig
	pgConfig               config.PgConfig
	redisConfig            config.RedisConfig
	userKafkaServiceConfig config.UserKafkaServiceConfig
	userCacheConfig        config.UserCacheConfig
	tokenConfig            config.TokenConfig
	kafkaProducerConfig    config.KafkaProducerConfig
	kafkaConsumerConfig    config.KafkaConsumerConfig

	pgDbClient    postgres.Client
	txManager     postgres.TxManager
	redisDbClient redis.Client
	auditRepo     repository.AuditRepository
	userCacheRepo repository.UserCacheRepository
	userRepo      repository.UserRepository
	accessRepo    repository.AccessRepository

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
	passwordService     service.PasswordService
	tokenService        service.TokenService
	accessService       service.AccessService
	authService         service.AuthService

	accessImpl *access.Implementation
	authImpl   *auth.Implementation
	userImpl   *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := serverCfg.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpcConfig: %s", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := serverCfg.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get httpConfig: %s", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := serverCfg.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swaggerConfig: %s", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PGConfig() config.PgConfig {
	if s.pgConfig == nil {
		cfg, err := dbCfg.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pgConfig: %s", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := dbCfg.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redisConfig: %s", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) UserKafkaServiceConfig() config.UserKafkaServiceConfig {
	if s.userKafkaServiceConfig == nil {
		cfg, err := serviceCfg.NewUserKafkaServiceConfig()
		if err != nil {
			log.Fatalf("failed to get userKafkaServiceConfig: %s", err)
		}

		s.userKafkaServiceConfig = cfg
	}

	return s.userKafkaServiceConfig
}

func (s *serviceProvider) UserCacheConfig() config.UserCacheConfig {
	if s.userCacheConfig == nil {
		cfg, err := serviceCfg.NewUserCacheConfig()
		if err != nil {
			log.Fatalf("failed to get userCacheConfig: %s", err)
		}

		s.userCacheConfig = cfg
	}

	return s.userCacheConfig
}

func (s *serviceProvider) TokenConfig() config.TokenConfig {
	if s.tokenConfig == nil {
		cfg, err := serviceCfg.NewTokenConfig()
		if err != nil {
			log.Fatalf("failed to get tokenConfig: %s", err)
		}

		s.tokenConfig = cfg
	}

	return s.tokenConfig
}

func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		cfg, err := kafkaCfg.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafkaProducerConfig: %s", err)
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := kafkaCfg.NewKafkaConsumerConfig()
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
			logger.Fatalf("failed to create db client: %s", err)
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
			log.Fatalf("redis ping error: %s", err)
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

func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepo == nil {
		s.accessRepo = accessRepo.NewRepository(s.PgDbClient(ctx))
	}
	return s.accessRepo
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
		s.userService = userService.NewService(s.UserRepository(ctx), s.PasswordService(), s.UserCacheService(ctx), s.AuditService(ctx), s.TxManager(ctx))
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

func (s *serviceProvider) PasswordService() service.PasswordService {
	if s.passwordService == nil {
		s.passwordService = passwordService.NewService()
	}

	return s.passwordService
}

func (s *serviceProvider) TokenService() service.TokenService {
	if s.tokenService == nil {
		s.tokenService = tokenService.NewService()
	}

	return s.tokenService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.TokenConfig(),
			s.AccessRepository(ctx),
			s.UserService(ctx),
			s.PasswordService(),
			s.TokenService())
	}

	return s.accessService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.TokenConfig(),
			s.UserService(ctx),
			s.PasswordService(),
			s.TokenService())
	}

	return s.authService
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx), s.UserProducerService())
	}

	return s.userImpl
}
