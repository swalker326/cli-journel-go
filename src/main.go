package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

const (
	Gray  = "\033[90m" //"\033[37m" // Or "\033[90m" for dark gray
	Reset = "\033[0m"
)

func main() {
	// Load settings and entries
	journal_dir, err := ExpandHomeDir("~/journal")
	if err != nil {
		panic(err)
	}
	settings_file, err := ExpandHomeDir(journal_dir + "/settings.json")
	if err != nil {
		panic(err)
	}
	entriesFilePath, err := ExpandHomeDir(journal_dir + "/entries.json")
	if err != nil {
		panic(err)
	}
	//Create journal directory if it does not exist
	if !PathExists(journal_dir) {
		dirErr := os.Mkdir(journal_dir, os.ModePerm)
		if dirErr != nil {
			panic(dirErr)
		}
		fmt.Printf("Journal directory created %s%s%s", Gray, journal_dir, Reset+"\n")
	} else {
		fmt.Printf("Journal directory exists %s%s%s", Gray, journal_dir, Reset+"\n")
	}
	//Create settings file if it does not exist
	if !PathExists(settings_file) {
		_, settingsErr := os.Create(settings_file)
		if settingsErr != nil {
			panic(settingsErr)
		}
		fmt.Printf("Settings file created %s%s%s", Gray, settings_file, Reset+"\n")
	} else {
		fmt.Printf("Settings file exists %s%s%s", Gray, settings_file, Reset+"\n")
	}
	//Create entries file if it does not exist
	if !PathExists(entriesFilePath) {
		_, entriesErr := os.Create(entriesFilePath)
		if entriesErr != nil {
			panic(entriesErr)
		}
		fmt.Printf("Entries file created %s%s%s", Gray, entriesFilePath, Reset+"\n")
	} else {
		fmt.Printf("Entries file exists %s%s%s", Gray, entriesFilePath, Reset+"\n")
	}

	settingsPath, err := LoadSettings(settings_file)
	if err != nil && os.IsNotExist(err) {
		SaveSettings(settingsPath, settings_file)
	}

	entries, err := LoadEntries(entriesFilePath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Using color to enhance CLI output
		color.HiCyan("\nJournal App")
		color.Yellow("1. Add Entry")
		color.Yellow("2. List Entries")
		color.Yellow("3. Settings")
		color.Yellow("4. Delete Entry")
		color.Yellow("5. Exit")
		color.Green("Choose an option: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Print("Enter title: ")
			scanner.Scan()
			title := scanner.Text()

			fmt.Print("Opening Editor")
			content, editorErr := EditJournalEntry()
			if editorErr != nil {
				fmt.Println("Error opening editor:", err)
				continue
			}

			entries = AddEntry(entries, title, content)
			err := SaveEntries(entries, settingsPath)
			if err != nil {
				fmt.Println("Error saving entry:", err)
			}

		case "2":
			for _, entry := range entries {
				fmt.Printf("\nDate: %s\nTitle: %s\nContent: %s\n", entry.Date.Format("2006-01-02 15:04:05"), entry.Title, entry.Content)
			}

		case "3":
			fmt.Printf("Enter new file path for entries or press Enter to keep the current path (%s%s%s): ", Gray, settingsPath, Reset)
			scanner.Scan()
			newPath := scanner.Text()

			if PathExists(newPath) {
				settingsPath = newPath
				color.Green("File path updated.")
			} else {
				color.Red("Path does not exist or is not writable.")
			}
		case "4":
			fmt.Println("Select an entry to delete:")
			for i, entry := range entries {
				fmt.Printf("%d. %s\n", i+1, entry.Title)
			}
			scanner.Scan()
			choice := scanner.Text()
			index, err := strconv.Atoi(choice)
			if err != nil {
				fmt.Println("Invalid choice, please try again.")
				continue
			}
			if index < 1 || index > len(entries) {
				fmt.Println("Invalid choice, please try again.")
				continue
			}
			entries = DeleteEntry(entries, index-1)
			err = SaveEntries(entries, settingsPath)
			if err != nil {
				fmt.Println("Error saving entries:", err)
			}
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
