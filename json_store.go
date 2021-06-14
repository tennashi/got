package got

import (
	"log"

	"github.com/schollz/jsonstore"
)

type JSONStoreConfig struct {
	FilePath string
	IsDebug  bool
}

type JSONStore struct {
	filePath string

	debugL *log.Logger
}

func NewJSONStore(ioStream *IOStream, cfg *JSONStoreConfig) (*JSONStore, error) {
	if cfg.FilePath == "" {
		return nil, &InvalidParamError{Message: "file path must not be empty"}
	}

	if _, err := jsonstore.Open(cfg.FilePath); err != nil {
		if err := jsonstore.Save(&jsonstore.JSONStore{}, cfg.FilePath); err != nil {
			panic(err)
		}
	}

	return &JSONStore{
		filePath: cfg.FilePath,
		debugL:   NewDebugLogger(ioStream.Err, "store", cfg.IsDebug),
	}, nil
}

func (s *JSONStore) Get(key string, v interface{}) error {
	s.debugL.Printf("start (*JSONStore).Get(%s, %v)\n", key, v)

	store, err := jsonstore.Open(s.filePath)
	if err != nil {
		s.debugL.Printf("error occurred in jsonstore.Open(): %v\n", err)

		return err
	}

	s.debugL.Printf("end (*JSONStore).Get(%s, %v)\n", key, v)

	return store.Get(key, v)
}

func (s *JSONStore) Save(key string, v interface{}) error {
	s.debugL.Printf("start (*JSONStore).Save(%s, %v)\n", key, v)

	store, err := jsonstore.Open(s.filePath)
	if err != nil {
		s.debugL.Printf("error occurred in jsonstore.Open(): %v\n", err)

		return err
	}

	if err := store.Set(key, v); err != nil {
		s.debugL.Printf("error occurred in store.Set(): %v\n", err)

		return err
	}

	s.debugL.Printf("end (*JSONStore).Save(%s, %v)\n", key, v)

	return jsonstore.Save(store, s.filePath)
}
