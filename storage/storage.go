package storage

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path"
)

type FileStore[K any] struct {
	path   string
	prefix string
	reset  func(item *K)
}

func NewFileStore[K any](path string, prefix string, reset func(*K)) *FileStore[K] {
	createFolderIfNotExists(path)
	return &FileStore[K]{
		path:   path,
		prefix: prefix,
		reset:  reset,
	}
}

func (fs *FileStore[K]) Load(id int) (*K, error) {
	filename := path.Join(fs.path, fmt.Sprintf("%s%d.gob", fs.prefix, id))
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			data := new(K)
			fs.reset(data)
			return data, nil
		}
		return nil, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	data := new(K)
	fs.reset(data)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fs *FileStore[K]) Save(id int, data *K) error {
	filename := path.Join(fs.path, fmt.Sprintf("%s%d.gob", fs.prefix, id))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(data)
}

func createFolderIfNotExists(folderName string) {
	info, err := os.Stat(folderName)
	if !os.IsNotExist(err) {
		if !info.IsDir() {
			log.Panicf("Unable to create directory '%s', a file with the same name already exists", folderName)
		}
		return
	}
	err = os.MkdirAll(folderName, 0644)
	if err != nil {
		log.Panic(err)
	}
}
