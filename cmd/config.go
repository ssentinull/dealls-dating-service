package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase"
	restserver "github.com/ssentinull/dealls-dating-service/internal/handler/rest"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/grace"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/healthcheck"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpclient"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmiddleware"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmux"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpserver"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/swagger"

	goValidator "github.com/go-playground/validator/v10"
)

type Options struct {
	Env     string `validate:"required"`
	Version string
	Name    string `validate:"required"`
}

type Handler struct {
	REST restserver.Options
}

type Business struct {
	Domain  domain.Options
	Usecase usecase.Options
}

type Conf struct {
	// Handler Options
	Handler Handler

	// Business Options
	Business Business

	// Application Metadata Options
	Meta Options

	// Infrastructure Options
	Logger      logger.Options
	SQL         libsql.Options
	Parser      parser.Options
	HTTPMW      httpmiddleware.Options
	HTTPMux     httpmux.Options
	HTTP        httpserver.Options
	Redis       cache.Options
	Auth        auth.Options
	Swagger     swagger.Options
	HealthCheck healthcheck.Options
	HTTPClient  httpclient.Options

	// Application Options
	Grace grace.Options
}

func init() {
	initConfigYaml()
	// initConfigEnv()
}

func initConfigYaml() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config") // name of Config file (without extension)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(errors.New("can not load config"))
	}

	conf = Conf{}

	// Unmarshal the configuration into the struct
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %v", err))
	}

	// Perform validation
	validate := goValidator.New()
	if err := validate.Struct(conf); err != nil {
		if validationErrors, ok := err.(goValidator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				log.Printf("Validation error for field '%s': %s", ve.StructNamespace(), ve.Tag())
			}
		} else {
			log.Fatalf("Validation error: %v", err)
		}

		panic(errors.New("invalid config"))
	}

	// Load Parser
	conf.Parser = parser.Options{
		JSON: parser.JSONOptions{
			Config: parser.JSONConfigDefault,
		},
	}
}

func initConfigEnv() (err error) {
	defer func() {
		errRecov := recover()
		if errRecov != nil {
			err = errRecov.(error)
		}
	}()

	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	conf = Conf{}

	conf.Auth.SecretKey = viper.GetString("AUTH_SECRET_KEY")
	conf.Auth.StaticToken = viper.GetString("AUTH_STATIC_TOKEN")

	conf.Grace.Network = viper.GetString("GRACE_NETWORK")
	conf.Grace.Pidfile = viper.GetString("GRACE_PIDFILE")
	conf.Grace.ShutdownTimeout = viper.GetDuration("GRACE_SHUTDOWN_TIMEOUT")
	conf.Grace.UpgradeTimeout = viper.GetDuration("GRACE_UPGRADE_TIMEOUT")

	conf.HealthCheck.MaxWaitingTime = viper.GetDuration("HEALTH_CHECK_MAX_WAITING_TIME")
	conf.HealthCheck.Readiness.CheckTimeout = viper.GetDuration("HEALTH_CHECK_READINESS_CHECK_TIMEOUT")
	conf.HealthCheck.Readiness.Enabled = viper.GetBool("HEALTH_CHECK_READINESS_ENABLED")
	conf.HealthCheck.Readiness.FailureThreshold = viper.GetInt("HEALTH_CHECK_READINESS_FAILURE_THRESHOLD")
	conf.HealthCheck.Readiness.InitDelaySec = viper.GetDuration("HEALTH_CHECK_READINESS_INIT_DELAY_SEC")
	conf.HealthCheck.Readiness.PeriodSec = viper.GetDuration("HEALTH_CHECK_READINESS_PERIOD_SEC")
	conf.HealthCheck.Readiness.SuccessThreshold = viper.GetInt("HEALTH_CHECK_READINESS_SUCCESS_THRESHOLD")
	conf.HealthCheck.WaitBeforeContinue = viper.GetBool("HEALTH_CHECK_WAIT_BEFORE_CONTINUE")

	conf.HTTPMW.AppName = viper.GetString("HTTPMW_APP_NAME")
	conf.HTTPMW.Env = viper.GetString("HTTPMW_ENV")
	conf.HTTPMW.Signature.Secret = viper.GetString("HTTPMW_SIGNATURE_SECRET")
	conf.HTTPMW.Signature.TimeTolerance = viper.GetInt("HTTPMW_SIGNATURE_TIME_TOLERANCE")

	conf.HTTP.Address = viper.GetString("HTTP_ADDRESS")
	conf.HTTP.Port = viper.GetInt("HTTP_PORT")
	conf.HTTP.TimeoutIdle = viper.GetDuration("HTTP_TIMEOUT_IDLE")
	conf.HTTP.TimeoutRead = viper.GetDuration("HTTP_TIMEOUT_READ")
	conf.HTTP.TimeoutReadHeader = viper.GetDuration("HTTP_TIMEOUT_READ_HEADER")
	conf.HTTP.TimeoutWrite = viper.GetDuration("HTTP_TIMEOUT_WRITE")

	conf.HTTPClient.BackOffInterval = viper.GetInt("HTTP_CLIENT_BACK_OFF_INTERVAL")
	conf.HTTPClient.DefaultClientID = viper.GetString("HTTP_CLIENT_DEFAULT_CLIENT_ID")
	conf.HTTPClient.MaximumJitterInterval = viper.GetInt("HTTP_CLIENT_MAXIMUM_JITTER_INTERVAL")
	conf.HTTPClient.MaxRetryCount = viper.GetInt("HTTP_CLIENT_MAX_RETRY_COUNT")
	conf.HTTPClient.Timeout = viper.GetInt("HTTP_CLIENT_TIMEOUT")

	conf.HTTPMux.AppEnv = viper.GetString("HTTP_MUX_APP_ENV")
	conf.HTTPMux.AppName = viper.GetString("HTTP_MUX_APP_NAME")

	conf.Meta.Env = viper.GetString("META_ENV")
	conf.Meta.Name = viper.GetString("META_NAME")
	conf.Meta.Version = viper.GetString("META_VERSION")

	conf.Redis.Address = viper.GetStringSlice("REDIS_ADDRESS")
	conf.Redis.DB = viper.GetInt("REDIS_DB")
	conf.Redis.DialTimeout = viper.GetDuration("REDIS_DIAL_TIMEOUT")
	conf.Redis.Enabled = viper.GetBool("REDIS_ENABLED")
	conf.Redis.IdleCheckFrequency = viper.GetDuration("REDIS_IDLE_CHECK_FREQUENCY")
	conf.Redis.IdleTimeout = viper.GetDuration("REDIS_IDLE_TIMEOUT")
	conf.Redis.InsecureSkipVerify = viper.GetBool("REDIS_INSECURE_SKIP_VERIFY")
	conf.Redis.MaxConnAge = viper.GetDuration("REDIS_MAX_CONN_AGE")
	conf.Redis.MaxRetries = viper.GetInt("REDIS_MAX_RETRIES")
	conf.Redis.MaxRetryBackoff = viper.GetDuration("REDIS_MAX_RETRY_BACKOFF")
	conf.Redis.MinIdleConns = viper.GetInt("REDIS_MIN_IDLE_CONNS")
	conf.Redis.MinRetryBackoff = viper.GetDuration("REDIS_MIN_RETRY_BACKOFF")
	conf.Redis.Mock = viper.GetBool("REDIS_MOCK")
	conf.Redis.Password = viper.GetString("REDIS_PASSWORD")
	conf.Redis.PoolSize = viper.GetInt("REDIS_POOL_SIZE")
	conf.Redis.PoolTimeout = viper.GetDuration("REDIS_POOL_TIMEOUT")
	conf.Redis.ReadTimeout = viper.GetDuration("REDIS_READ_TIMEOUT")
	conf.Redis.WriteTimeout = viper.GetDuration("REDIS_WRITE_TIMEOUT")

	conf.SQL.Follower.DSN = viper.GetString("SQL_FOLLOWER_DSN")
	conf.SQL.Follower.Enabled = viper.GetBool("SQL_FOLLOWER_ENABLED")
	conf.SQL.Follower.Mock = viper.GetBool("SQL_FOLLOWER_MOCK")

	conf.SQL.Leader.DSN = viper.GetString("SQL_LEADER_DSN")
	conf.SQL.Leader.Enabled = viper.GetBool("SQL_LEADER_ENABLED")
	conf.SQL.Leader.Mock = viper.GetBool("SQL_LEADER_MOCK")

	conf.Swagger.DocPath = viper.GetString("SWAGGER_DOC_PATH")
	conf.Swagger.Name = viper.GetString("SWAGGER_NAME")
	conf.Swagger.Path = viper.GetString("SWAGGER_PATH")

	// Load Parser
	conf.Parser = parser.Options{
		JSON: parser.JSONOptions{
			Config: parser.JSONConfigDefault,
		},
	}

	// Perform validation
	validate := goValidator.New()
	if err := validate.Struct(conf); err != nil {
		if validationErrors, ok := err.(goValidator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				log.Printf("Validation error for field '%s': %s", ve.StructNamespace(), ve.Tag())
			}
		} else {
			log.Fatalf("Validation error: %v", err)
		}
		panic(errors.New("invalid config"))
	}

	return nil
}
