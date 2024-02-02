package settings

import (
	"time"
)

const (
	Debug   string = "debug"
	Release string = "release"
	Test    string = "test"
)

type Settings struct {
	Global struct {
		Name              string        `yaml:"name" env:"GLOBAL_NAME" env-default:"sad-server" env-description:"Instance Name"`
		ReadTimeout       time.Duration `yaml:"readTimeout" env:"GLOBAL_READ_TIMEOUT" env-default:"2m" env-description:"Read timeout of http server"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" env:"GLOBAL_READ_HEADER_TIMEOUT" env-default:"2m" env-description:"Read header timeout of http server"`
		WriteTimeout      time.Duration `yaml:"writeTimeout" env:"GLOBAL_WRITE_TIMEOUT" env-default:"2m" env-description:"Write timeout of http server"`
		IdleTimeout       time.Duration `yaml:"idleTimeout" env:"GLOBAL_IDLE_TIMEOUT" env-default:"2m" env-description:"Idle timeout of http server"`
		MaxHeaderBytes    int           `yaml:"maxHeaderBytes" env:"GLOBAL_MAX_HEADER_BYTES" env-default:"8196" env-description:"Max header bytes of http server"`
		APIPort           int           `yaml:"apiPort" env:"GLOBAL_API_PORT" env-default:"8080" env-description:"Default Port of API server"`
		MemberlistPort    int           `yaml:"memberlistPort" env:"GLOBAL_MEMBER_LIST_PORT" env-default:"8081" env-description:"Default Port of Memberlist server"`
		GossopingPort     int           `yaml:"gossopingPort" env:"GLOBAL_GOSSOPING_PORT" env-default:"8082" env-description:"Default Port of Gossoping server"`
		Environment       string        `yaml:"environment" env:"CONFIG_MODE" env-default:"file" env-description:"Execution mode of Gin framework"`
	} `yaml:"global"`
	Replica struct {
		Hostname []string `yaml:"hostname" env:"HOSTNAME" env-default:"localhost" env-description:"Base hostname of replicas"`
	} `yaml:"replica"`
}

func (settings Settings) IsValid() (bool, error) {
	if settings.Global.Name == "" {
		return false, ErrSettingNameEmpty
	}

	// Check Ports duplications
	var hasDuplicatedPort bool
	duplicatedPorts := make(map[int]bool)
	for _, item := range []int{
		settings.Global.APIPort,
		settings.Global.MemberlistPort,
		settings.Global.GossopingPort,
	} {
		_, exist := duplicatedPorts[item]
		if exist {
			hasDuplicatedPort = true
			duplicatedPorts[item] = true
		} else {
			duplicatedPorts[item] = false
		}
	}
	if hasDuplicatedPort {
		return false, ErrSettingDuplicatedServerPorts
	}

	if settings.Global.Environment != Debug && settings.Global.Environment != Release && settings.Global.Environment != Test {
		return false, ErrSettingInvalidEnvironment
	}
	return true, nil
}
