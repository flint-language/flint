package toml

type BuildConfig struct {
	Type     string `toml:"type"`
	Language string `toml:"language"`
	Entry    string `toml:"entry"`
	Output   string `toml:"output"`
}

func newExampleBuildConfig() BuildConfig {
	return BuildConfig{
		Type:     "native",
		Language: "c",
		Entry:    "example.c",
		Output:   "example",
	}
}
