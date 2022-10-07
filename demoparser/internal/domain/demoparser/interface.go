package demoparser

import "github.com/Cludch/csgo-microservices/demoparser/pkg/files"

type UseCase interface {
	Parse(dir string, demoFile *files.Demo) error
}
