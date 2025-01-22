package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kenzo0107/backlog"
)

const (
	spaceName = ""
	baseUrl   = "https://" + spaceName + ".backlog.com"
	projectID = ""
	apiKey    = ""
)

func main() {
	// wikiListの取得
	c := backlog.New(apiKey, baseUrl)

	wikis, err := c.GetWikis(&backlog.GetWikisOptions{ProjectIDOrKey: projectID})
	if err != nil {
		log.Fatalf("failed to get wikis: %v", err)
	}

	// エラーログ用のファイルを作成
	errorFile, err := os.OpenFile("failed_wiki_ids.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to create error log file: %v", err)
	}
	defer errorFile.Close()

	for _, v := range wikis {
		if err := exportWiki(c, "output", *v.ID); err != nil {
			// エラーをログに出力
			log.Printf("failed to export wiki ID %d: %v", *v.ID, err)

			// エラーが発生したWiki IDをファイルに書き込む
			_, writeErr := fmt.Fprintf(errorFile, "%d\n", *v.ID)
			if writeErr != nil {
				log.Printf("failed to write error log file: %v", writeErr)
			}
		}
	}

	log.Println("backup successfully")
}

func exportWiki(c *backlog.Client, baseDir string, wikiID int) error {
	wiki, err := c.GetWiki(wikiID)
	if err != nil {
		return fmt.Errorf("failed to get wiki: %v", err)
	}

	// wiki用のディレクトリ作成
	wikiName := *wiki.Name
	wikiDir := filepath.Join(baseDir, wikiName)
	if err := os.MkdirAll(wikiDir, 0755); err != nil {
		return fmt.Errorf("failed to create wiki directory: %v", err)
	}

	// wikiの保存
	content := []byte(*wiki.Content)
	sanitizeWikiName := sanitizeFilename(*wiki.Name)
	contentPath := filepath.Join(wikiDir, sanitizeWikiName+".md")
	if err := os.WriteFile(contentPath, content, 0644); err != nil {
		return fmt.Errorf("failed to save wiki content: %v", err)
	}

	// 添付ファイルのダウンロード
	for _, attachment := range wiki.Attachments {
		fileName := *attachment.Name
		filePath := filepath.Join(wikiDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("failed to create file for attachment %s: %v", *attachment.Name, err)
			continue
		}
		defer file.Close()

		err = c.GetWikiAttachmentContent(*wiki.ID, *attachment.ID, file)
		if err != nil {
			log.Printf("failed to download attachment %s: %v", *attachment.Name, err)
			continue
		}

		log.Printf("save attachment: %s", fileName)
	}

	// wikiのレスポンスをメタデータとして保存する
	metadata, err := json.MarshalIndent(wiki, "", "  ")
	if err != nil {
		log.Printf("failed to marshal wiki metadata: %v", err)
	} else {
		metadataPath := filepath.Join(wikiDir, sanitizeWikiName+"_metadata.json")
		if err := os.WriteFile(metadataPath, metadata, 0644); err != nil {
			log.Printf("failed to save wiki metadata: %v", err)
		}
	}

	log.Printf("successfully exported wiki '%s' to directory: %s", *wiki.Name, wikiDir)
	return nil
}

// ファイル名として使用できない文字をサニタイズする関数
func sanitizeFilename(filename string) string {
	// Windowsでも使用できるよう、一般的な禁止文字をすべて置換
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}
