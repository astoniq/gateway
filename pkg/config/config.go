package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	FileName = "janus.yaml"
	Prefix   = "janus"
)

type Config struct {
	Port                 int             `mapstructure:"port"`
	GraceTimeOut         int64           `mapstructure:"graceTimeOut"`
	MaxIdleConsPerHost   int             `mapstructure:"maxIdleConsPerHost"`
	BackendFlushInterval time.Duration   `mapstructure:"backendFlushInterval"`
	IdleConnTimeout      time.Duration   `mapstructure:"idleConnTimeout"`
	ConnPurgeInterval    time.Duration   `mapstructure:"connPurgeInterval"`
	RequestId            bool            `mapstructure:"requestId"`
	LogLevel             int8            `mapstructure:"logLevel"`
	Web                  *WebConfig      `mapstructure:"web"`
	Cluster              *ClusterConfig  `mapstructure:"cluster"`
	Database             *DatabaseConfig `mapstructure:"database"`
}

type ClusterConfig struct {
	UpdateFrequency time.Duration `mapstructure:"updateFrequency"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type WebConfig struct {
	Port int `mapstructure:"port"`
}

func NewConfig() Config {

	viper.SetConfigFile(FileName)

	viper.SetDefault("port", "8080")
	viper.SetDefault("logLevel", zerolog.DebugLevel)
	viper.SetDefault("backendFlushInterval", "20ms")
	viper.SetDefault("requestId", true)

	viper.SetDefault("web.port", "8090")

	_, err := os.ReadFile(FileName)
	if err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Msg("init: error while reading imperium.yaml. Shutting down.")
		}
	} else {
		if os.IsNotExist(err) {
			log.Info().Msg("init: could not find imperium.yaml. Attempting to use environment variables.")
		} else {
			log.Fatal().Err(err).Msg("init: error while reading imperium.yaml. Shutting down.")
		}
	}

	var config Config
	for _, fieldName := range getFlattenedStructFields(reflect.TypeOf(config)) {
		envKey := strings.ToUpper(fmt.Sprintf("%s_%s", Prefix, strings.ReplaceAll(fieldName, ".", "-")))
		envVar := os.Getenv(envKey)
		if envVar != "" {
			viper.Set(fieldName, envVar)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("init: error while creating config. Shutting down.")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return config
}

func getFlattenedStructFields(t reflect.Type) []string {
	return getFlattenedStructFieldsHelper(t, []string{})
}

func getFlattenedStructFieldsHelper(t reflect.Type, prefixes []string) []string {
	unwrappedT := t
	if t.Kind() == reflect.Pointer {
		unwrappedT = t.Elem()
	}

	flattenedFields := make([]string, 0)
	for i := 0; i < unwrappedT.NumField(); i++ {
		field := unwrappedT.Field(i)
		fieldName := field.Tag.Get("mapstructure")
		switch field.Type.Kind() {
		case reflect.Struct, reflect.Pointer:
			flattenedFields = append(flattenedFields, getFlattenedStructFieldsHelper(field.Type, append(prefixes, fieldName))...)
		default:
			flattenedField := fieldName
			if len(prefixes) > 0 {
				flattenedField = fmt.Sprintf("%s.%s", strings.Join(prefixes, "."), fieldName)
			}
			flattenedFields = append(flattenedFields, flattenedField)
		}
	}
	return flattenedFields
}
