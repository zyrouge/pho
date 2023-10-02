package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

var AtomicFilePrefix = "pho"

func WriteFileAtomic(name string, bytes []byte) error {
	dir := path.Dir(name)
	ext := path.Ext(name)
	temp, err := os.CreateTemp(
		dir,
		fmt.Sprintf(
			"%s-%s-*%s",
			AtomicFilePrefix,
			strings.TrimSuffix(path.Base(name), ext),
			ext,
		),
	)
	if err != nil {
		return err
	}
	tempName := temp.Name()
	closed, deleted := false, false
	defer func() {
		if !closed {
			temp.Close()
		}
		if !deleted {
			os.Remove(tempName)
		}
	}()
	if _, err = temp.Write(bytes); err != nil {
		return err
	}
	if err = temp.Sync(); err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	closed = true
	err = os.Rename(tempName, name)
	if err != nil {
		return err
	}
	deleted = true
	return nil
}

func ReadJsonFile[T any](name string) (*T, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var data T
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func WriteJsonFile[T any](name string, data *T) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(name, json, os.ModePerm)
}

func WriteJsonFileAtomic[T any](name string, data *T) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return WriteFileAtomic(name, json)
}
