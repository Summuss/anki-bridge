package config

var Conf *Config

type Config struct {
	MongoConnectURL string
	DBNAme          string
	AnkiAPIURL      string
}
