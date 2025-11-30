package toml

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	ProjectConfig struct {
		PackageConfig PackageConfig `toml:"package"`
		BuildConfig   BuildConfig   `toml:"build"`
	}
)

var configFile = "flint.toml"

func newExampleProjectConfig() ProjectConfig {
	return ProjectConfig{
		PackageConfig: newExamplePackageConfig(),
		BuildConfig: newExampleBuildConfig(),
	}
}

func NewExampleProjectConfigAt(path string) {
	SaveProjectConfigAt(newExampleProjectConfig(), path)
}

func LoadProjectConfig() (ProjectConfig, error) {
	return LoadProjectConfigAt("")
}

func LoadProjectConfigAt(path string) (ProjectConfig, error){
	var config ProjectConfig
	_, err := toml.DecodeFile(path + "/" + configFile, &config)
	return config, err 
}

func SaveProjectConfig(config ProjectConfig) {
	SaveProjectConfigAt(config, "")
}

func SaveProjectConfigAt(config ProjectConfig, path string) error {
	os.MkdirAll(path, os.ModePerm)

	file, fileErr := os.Create(path + "/" + configFile)

	if fileErr != nil {
		return fileErr
	}

	err := toml.NewEncoder(file).Encode(config)

	if err != nil {
		return err
	}
	return nil
}

func PrintProjectConfig(config ProjectConfig) {
	fmt.Printf("%+v", config)
}


