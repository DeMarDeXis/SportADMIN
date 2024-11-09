package customjsonexp

import (
	"encoding/json"
	"os"
	"sync"
)

type CustomJSONExporter struct {
	file  *os.File
	mutex sync.Mutex
	enc   *json.Encoder
}

func NewCustomJSONExporter(filename string) (*CustomJSONExporter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)

	//_, err = file.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &CustomJSONExporter{file: file, enc: enc}, nil
}

func (e *CustomJSONExporter) Export(exports chan interface{}) error {
	e.file.WriteString("[\n")
	defer e.file.WriteString("\n]")

	first := true
	for data := range exports {
		if !first {
			e.file.WriteString(",\n")
		}
		first = false

		bytes, err := json.MarshalIndent(data, "  ", "  ")
		if err != nil {
			return err
		}
		e.file.Write(bytes)
	}

	return nil
}

//func (e *CustomJSONExporter) Close() error {
//	e.file.Seek(-1, io.SeekEnd)
//	_, err := e.file.WriteString("]")
//	if err != nil {
//		return err
//	}
//	return e.file.Close()
//}
