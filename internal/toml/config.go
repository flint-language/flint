package toml

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	ProjectConfig struct {
		PackageConfig packageConfig `toml:"package"`
		BuildConfig   buildConfig   `toml:"build"`
	}
)

var configFile = "flint.toml"

func newExampleProjectConfig() ProjectConfig {
	return ProjectConfig{
		PackageConfig: newExamplePackageConfig(),
		BuildConfig: newExampleBuildConfig(),
	}
}

func newExampleProjectConfigAt(path string) {
	saveProjectConfigAt(newExampleProjectConfig(), path)
}

func loadProjectConfig() (ProjectConfig, error) {
	return loadProjectConfigAt("")
}

func loadProjectConfigAt(path string) (ProjectConfig, error){
	var config ProjectConfig
	_, err := toml.DecodeFile(path + "/" + configFile, &config)
	return config, err 
}

func saveProjectConfig(config ProjectConfig) {
	saveProjectConfigAt(config, "")
}

func saveProjectConfigAt(config ProjectConfig, path string) error {
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

func printProjectConfig(config ProjectConfig) {
	fmt.Printf("%+v", config)
}


