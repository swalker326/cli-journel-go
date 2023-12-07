package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// JournalEntry represents a single journal entry
type JournalEntry struct {
	Date    time.Time
	Title   string
	Content string
}

// LoadEntries loads journal entries from a file
func LoadEntries(entriesFilePath string) ([]JournalEntry, error) {
	var entries []JournalEntry
	EntriesFilePath, err := ExpandHomeDir(entriesFilePath)
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(EntriesFilePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &entries)
	return entries, err
}

// SaveEntries saves journal entries to a file
func SaveEntries(entries []JournalEntry, entriesFilePath string) error {
	pathErr := PathExists(entriesFilePath)
	if pathErr {
		file, err := json.MarshalIndent(entries, "", " ")
		if err != nil {
			return err
		}
		return os.WriteFile(entriesFilePath, file, 0644)
	} else {
		return errors.New("path does not exist")
	}
}

// DeleteEntry deletes a journal entry from the journal
func DeleteEntry(entries []JournalEntry, index int) []JournalEntry {
	return append(entries[:index], entries[index+1:]...)
}

// AddEntry adds a new entry to the journal
func AddEntry(entries []JournalEntry, title, content string) []JournalEntry {
	entry := JournalEntry{
		Date:    time.Now(),
		Title:   title,
		Content: content,
	}
	return append(entries, entry)
}
