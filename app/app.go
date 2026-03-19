package app

import (
	"encoding/json"
	"log"

	contract "github.com/slidebolt/sb-contract"
	messenger "github.com/slidebolt/sb-messenger-sdk"
	server "github.com/slidebolt/sb-storage-server"
)

type App struct {
	cfg     Config
	msg     messenger.Messenger
	handler *server.Handler
}

func New(cfg Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Hello() contract.HelloResponse {
	return contract.HelloResponse{
		ID:              "storage",
		Kind:            contract.KindService,
		ContractVersion: contract.ContractVersion,
		DependsOn:       []string{"messenger"},
	}
}

func (a *App) OnStart(deps map[string]json.RawMessage) (json.RawMessage, error) {
	msg, err := messenger.Connect(deps)
	if err != nil {
		return nil, err
	}
	a.msg = msg

	a.handler, err = server.NewHandlerWithDir(a.cfg.DataDir)
	if err != nil {
		return nil, err
	}

	if n, err := a.handler.LoadFromDir(); err != nil {
		log.Printf("storage: no existing data: %v", err)
	} else {
		log.Printf("storage: loaded %d entries from disk", n)
	}

	if err := a.handler.StartWatcher(); err != nil {
		return nil, err
	}

	if err := a.handler.Register(msg); err != nil {
		return nil, err
	}

	log.Printf("storage: ready, serving on storage.>")
	return nil, nil
}

func (a *App) OnShutdown() error {
	if a.msg != nil {
		a.msg.Close()
	}
	return nil
}
