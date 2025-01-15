package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

const (
	targetDir = "output"
)

func main() {
	// metadata.jsonからwikiIDとwiki名, outputの絶対パスのマップを作成する
	err := filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".json") {
			
		}

		return nil
	})
	if err != nil {
		log.Fatalf("failed to create map: %v", err)
	}

	// backlog.comの参照を探して、metadataのマップから取得したwiki名とパスをcsv形式とかでファイル出力する
}
