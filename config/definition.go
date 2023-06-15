package config

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_manager.go config.go Manager
//counterfeiter:generate -o fakes/fake_reader.go config.go Reader
//counterfeiter:generate -o fakes/fake_updater.go config.go Updater
