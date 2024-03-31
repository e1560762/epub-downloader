package pkg

import (
	"encoding/json"
	"os"
)

type Config struct {
	HostConf map[string]domElement
}

type domElement struct {
	Fields []domField `json:"fields"`
	Text   string     `json:"text"`
}

type domField struct {
	Key      string `json:"key"`
	Contains string `json:"contains"`
}

func NewConfig(configFile string) (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	parser := json.NewDecoder(f)
	var conf = make(map[string]domElement)
	if err = parser.Decode(&conf); err != nil {
		return nil, err
	}
	return &Config{HostConf: conf}, nil
}

type filePaths []string

func NewFilePathsArg() filePaths {
	return filePaths{}
}

func (*filePaths) String() string {
	return "List of paths to download the files"
}

func (fp *filePaths) Set(value string) error {
	*fp = append(*fp, value)
	return nil
}

type originURLs []string

func NewOriginURLsArg() originURLs {
	return originURLs{}
}

func (*originURLs) String() string {
	return "List of source URLs to download from"
}

func (ou *originURLs) Set(value string) error {
	*ou = append(*ou, value)
	return nil
}
