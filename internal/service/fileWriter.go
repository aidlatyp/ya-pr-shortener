package service

import (
	"encoding/json"
	"os"
)

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Writer) Write(r *Record) error {
	return p.encoder.Encode(&r)
}
func (p *Writer) Close() error {
	return p.file.Close()
}
