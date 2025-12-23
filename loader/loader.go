package loader

import (
	"errors"
	"fmt"
	"os"
	"scope-guardian/logger"

	"github.com/pelletier/go-toml/v2"
)

func Load(filepath string) (Config, error) {
	if _, err := os.Stat(filepath); err != nil {
		logger.Error(fmt.Sprintf(logErrorFileNotFound, filepath))
		return Config{}, errors.New(errFileNotFound)
	}

	var config Config

	fd, _ := os.ReadFile(filepath)

	if err := toml.Unmarshal(fd, &config); err != nil {
		return Config{}, errors.New(errDecodingToml)
	}

	return config, nil
}
