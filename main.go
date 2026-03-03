package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fezcode/atlas.deck/internal/config"
	"github.com/fezcode/atlas.deck/internal/ui"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-v" || arg == "--version" {
			fmt.Printf("atlas.deck v%s\n", Version)
			return
		}
		if arg == "-h" || arg == "--help" || arg == "help" {
			showHelp()
			return
		}
	}

	deck, err := config.LoadDeck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading deck: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.NewModel(deck), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Printf("atlas.deck v%s - Interactive TUI command deck for project workflows.\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  atlas.deck                Start the interactive TUI")
	fmt.Println("  atlas.deck -h, --help     Show this help information")
	fmt.Println("  atlas.deck -v, --version  Show version information")
	fmt.Println("\nBlueprint:")
	fmt.Println("  The application looks for a 'deck.piml' file in the current directory.")
	fmt.Println("  If not found, it falls back to the global configuration in ~/.atlas/deck.piml.")
	fmt.Println("\nTUI Controls:")
	fmt.Println("  [Key]          Execute the command mapped to that key")
	fmt.Println("  q, esc, ctrl+c Quit the application")
}
