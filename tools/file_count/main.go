package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	mdExtension = ".md"
	separator   = "/"
)

func main() {
	// コマンドライン引数でルートディレクトリを受け取る
	var rootDir string
	flag.StringVar(&rootDir, "root", "output", "Root directory to scan")
	flag.Parse()

	// フラグで指定されていない場合は、最初の引数を使用
	if args := flag.Args(); len(args) > 0 {
		rootDir = args[0]
	}

	// 末尾のセパレータを削除
	rootDir = strings.TrimSuffix(rootDir, separator)

	// サブディレクトリごとの.mdファイル数を保持するマップ
	dirCounts := make(map[string]int)
	totalCount := 0

	err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), mdExtension) {
			// パスからrootDirを除去してサブディレクトリ名を取得
			relativePath := strings.TrimPrefix(path, rootDir+separator)
			if parts := strings.Split(relativePath, separator); len(parts) > 0 {
				dirCounts[parts[0]]++
				totalCount++
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	// 結果を出力
	fmt.Printf("%sディレクトリ内の.mdファイル数: %d\n", rootDir, totalCount)
	fmt.Println("\nサブディレクトリ別の内訳:")
	for dir, count := range dirCounts {
		fmt.Printf("%s: %d files\n", dir, count)
	}
}
