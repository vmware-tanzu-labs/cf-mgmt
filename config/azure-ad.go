package config

//Config -
type AzureADConfig struct {
	Enabled		bool   `yaml:"enabled"`
	ClientId	string `yaml:"client-id"`
	Secret		string `yaml:"client-secret"`
	TennantID	string `yaml:"tennant-id"`
	UserOrigin	string `yaml:"origin"`
}
