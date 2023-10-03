package config

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	cli "github.com/urfave/cli/v2"
)

type AnalyzerConfig struct {
	LogLevel       string      `json:"log-level" mapstructure:"LOG_LEVEL"`
	InitSlot       phase0.Slot `json:"init-slot" mapstructure:"INIT_SLOT"`
	FinalSlot      phase0.Slot `json:"final-slot" mapstructure:"FINAL_SLOT"`
	BnEndpoint     string      `json:"bn-endpoint" mapstructure:"BN_ENDPOINT"`
	ElEndpoint     string      `json:"el-endpoint" mapstructure:"EL_ENDPOINT"`
	DBUrl          string      `json:"db-url" mapstructure:"DB_URL"`
	DownloadMode   string      `json:"download-mode" mapstructure:"DOWNLOAD_MODE"`
	WorkerNum      int         `json:"worker-num" mapstructure:"WORKER_NUM"`
	DbWorkerNum    int         `json:"db-worker-num" mapstructure:"DB_WORKER_NUM"`
	Metrics        string      `json:"metrics" mapstructure:"METRICS"`
	PrometheusPort int         `json:"prometheus-port" mapstructure:"PROMETHEUS_PORT"`
}

// TODO: read from config-file
func NewAnalyzerConfig() *AnalyzerConfig {
	// Return Default values for the ethereum configuration
	return &AnalyzerConfig{}
}

func (c *AnalyzerConfig) LoadFromEnv() error {
	if err := godotenv.Overload(); err != nil {
		return fmt.Errorf("can't load .env file: %w", err)
	}

	viper.AutomaticEnv()

	// NOTE: workadound to load env vars
	// see https://github.com/spf13/viper/issues/761
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(c, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return err
		}
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("can't unmarshal env vars to config: %w", err)
	}

	return nil
}

func (c *AnalyzerConfig) Apply(ctx *cli.Context) error {
	if ctx.IsSet("env-file") {
		if err := c.LoadFromEnv(); err != nil {
			return err
		}
	}
	// apply to the existing Default configuration the set flags
	// log level
	if ctx.IsSet("log-level") {
		c.LogLevel = ctx.String("log-level")
	}
	// init slot
	if ctx.IsSet("init-slot") {
		c.InitSlot = phase0.Slot(ctx.Int("init-slot"))
	}
	// final slot
	if ctx.IsSet("final-slot") {
		c.FinalSlot = phase0.Slot(ctx.Int("final-slot"))
	}
	// cl url
	if ctx.IsSet("bn-endpoint") {
		c.BnEndpoint = ctx.String("bn-endpoint")
	}
	// el url
	if ctx.IsSet("el-endpoint") {
		c.ElEndpoint = ctx.String("el-endpoint")
	}
	// db url
	if ctx.IsSet("db-url") {
		c.DBUrl = ctx.String("db-url")
	}
	// download mode
	if ctx.IsSet("download-mode") {
		c.DownloadMode = ctx.String("download-mode")
	}
	// worker num
	if ctx.IsSet("worker-num") {
		c.WorkerNum = ctx.Int("worker-num")
	}
	// db worker num
	if ctx.IsSet("db-worker-num") {
		c.DbWorkerNum = ctx.Int("db-worker-num")
	}
	// metrics
	if ctx.IsSet("metrics") {
		c.Metrics = ctx.String("metrics")
	}
	// prometheus port
	if ctx.IsSet("prometheus-port") {
		c.PrometheusPort = ctx.Int("prometheus-port")
	}
	return nil
}
