package tokens

import (
	"encoding/json"
	"github.com/pandora_go/exts/logger"
	"os"
)

var TokenList = make(map[string]string)

func InitAccessToken(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			logger.Error("Error creating access_tokens.json")
		}
		defer file.Close()
	}
	// Read the file
	file, err := os.Open(path)
	if err != nil {
		logger.Error("Error opening access_tokens.json")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&TokenList)
	if err != nil {
		logger.Error("Error decoding access_tokens.json")
	}
}

func GetToken(key string) string {
	return TokenList[key]
}
