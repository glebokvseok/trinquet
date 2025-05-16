package config

import (
	"bytes"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/pkg/errors"
	"go.uber.org/config"
	"io"
	"os"
	"regexp"
)

const (
	envConfigFileMode          = 0644
	envConfigTemplateFilePath  = "config/config.env.yaml"
	envConfigGeneratedFilePath = "config/config.env.generated.yaml"
)

func GetAppEnvConfigProvider() (returnConf *config.YAML, returnErr error) {
	defer func() {
		if err := os.Remove(envConfigGeneratedFilePath); err != nil {
			returnErr = errors.Errorf("error occured while deleting generated env config: %v", err)
		}
	}()

	yaml, err := config.NewYAML(config.File(envConfigGeneratedFilePath))
	return yaml, errors.WithStack(err)
}

func PopulateEnvVars() (returnErr error) {
	file, err := os.Open(envConfigTemplateFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			returnErr = errors.WithStack(err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return errors.WithStack(err)
	}

	re := regexp.MustCompile(`\${([a-zA-Z_]+)}`)
	matches := re.FindAllSubmatch(data, -1)

	for _, match := range matches {
		envVar := match[1]
		value := os.Getenv(string(envVar))
		if extensions.IsEmpty(value) {
			return errors.Errorf("missing value for environment variable: %s", envVar)
		}

		data = bytes.ReplaceAll(data, match[0], []byte(value))
	}

	return errors.WithStack(os.WriteFile(envConfigGeneratedFilePath, data, envConfigFileMode))
}
