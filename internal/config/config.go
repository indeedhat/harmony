package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	. "github.com/indeedhat/harmony/internal/logger"
)

type Config struct {
	App struct {
		TransitionPollMs int `toml:"transition_poll_ms" validate:"required,min=10"`
	} `toml:"app"`

	EscapeSequence struct {
		KeyCount         int  `toml:"key_count" validate:"required,min=2"`
		TimeframeSeconds uint `toml:"time_seconds" validate:"required,min=1,max=5"`
	} `toml:"escape_sequence"`

	Discovery struct {
		MulticastAddress    string `toml:"multicast_address" validate:"required"`
		PollCaunt           int    `toml:"poll_count" validate:"required,min=1"`
		PollIntervalSeconds int    `toml:"poll_interval_seconds" validate:"required,min=1"`
	} `toml:"discovery"`

	Server struct {
		Port               int `toml:"port" validate:"required,min=1025,max=65535"`
		WsWriteWaitSecond  int `toml:"soc_write_wait_second" validate:"required,min=1,max=30"`
		WsCloseGracePeriod int `toml:"soc_close_grace_second" validate:"required,min=1,max=30"`
	}
}

// Load the config from file and validate its contents
func Load() *Config {
	var config Config

	data, err := ioutil.ReadFile("./config.toml")
	if err != nil {
		Log("config", "failed to load from file")
		return nil
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		Log("config", "failed to parse config file")
		return nil
	}

	v := validator.New()
	if err := v.Struct(config); err != nil {
		Logf("config", "invalid config: %s", err)
		return nil
	}

	return &config
}
