package config

//go:generate counterfeiter -o fakes/fake_manager.go config.go Manager
//go:generate counterfeiter -o fakes/fake_reader.go config.go Reader
//go:generate counterfeiter -o fakes/fake_updater.go config.go Updater
