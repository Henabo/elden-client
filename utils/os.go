package utils

import (
	"github.com/hiro942/elden-client/constant"
	"log"
	"os"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func WriteFile(path string, data []byte) {
	err := os.WriteFile(path, data, constant.DefaultFilePerm)
	if err != nil {
		log.Panicln("failed to write file:", err)
	}
}

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicln("failed to read file:", err)
	}
	return data
}
