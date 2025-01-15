package main

import (
	"encoding/json"
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

// wikiを取得し、添付ファイルを取得する
func main() {
	c := backlog.New(apiKey, baseUrl)

	wiki, err := c.GetWiki(3080262)
	if err != nil {
		log.Fatal(err)
	}

	// Create directory for wiki content
	wikiDir := sanitizeFileName(*wiki.Name)
	if err := os.MkdirAll(wikiDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Save wiki content
	content := []byte(*wiki.Content)
	contentPath := filepath.Join(wikiDir, "content.md")
	if err := os.WriteFile(contentPath, content, 0644); err != nil {
		log.Fatal(err)
	}

	// Create attachments directory
	attachmentsDir := filepath.Join(wikiDir, "attachments")
	if err := os.MkdirAll(attachmentsDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Download attachments
	for _, attachment := range wiki.Attachments {
		// Create file
		fileName := sanitizeFileName(*attachment.Name)
		filePath := filepath.Join(attachmentsDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file for attachment %s: %v", *attachment.Name, err)
			continue
		}
		defer file.Close()

		// Download attachment content directly to file
		err = c.GetWikiAttachmentContent(*wiki.ID, *attachment.ID, file)
		if err != nil {
			log.Printf("Error downloading attachment %s: %v", *attachment.Name, err)
			continue
		}

		log.Printf("Saved attachment: %s", fileName)
	}

	// Save wiki metadata as JSON
	metadata, err := json.MarshalIndent(wiki, "", "  ")
	if err != nil {
		log.Printf("Error marshaling wiki metadata: %v", err)
	} else {
		metadataPath := filepath.Join(wikiDir, "metadata.json")
		if err := os.WriteFile(metadataPath, metadata, 0644); err != nil {
			log.Printf("Error saving wiki metadata: %v", err)
		}
	}

	log.Printf("Successfully exported wiki '%s' to directory: %s", *wiki.Name, wikiDir)
}

// sanitizeFileName removes or replaces characters that are invalid in file names
func sanitizeFileName(name string) string {
	// Replace invalid characters with underscore
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name

	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}

	return result
}
