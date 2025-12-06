package toml

type packageConfig struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Authors     []string `toml:"authors"`
	Description string   `toml:"description"`
	License     string   `toml:"license"`
	Repository  string   `toml:"repository"`
}

func newExamplePackageConfig() packageConfig {
	return packageConfig{
		Name:        "example-package",
		Version:     "0.1.0",
		Authors:     []string{"author1", "author2"},
		Description: "A example library for Flint programming language.",
		License:     "MIT",
		Repository:  "https://github.com/example/example-repo",
	}
}
