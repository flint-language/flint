package toml

type buildConfig struct {
	Type     string `toml:"type"`
	Language string `toml:"language"`
	Entry    string `toml:"entry"`
	Output   string `toml:"output"`
}

func newExampleBuildConfig() buildConfig {
	return buildConfig{
		Type:     "native",
		Language: "c",
		Entry:    "example.c",
		Output:   "example",
	}
}
