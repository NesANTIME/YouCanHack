package main

import (
	"fmt"
	"time"
	"youcanhack/framework/core"

	"github.com/charmbracelet/lipgloss"
)

// funciones auxiliares ---

// ---

func main() {
	title_running := lipgloss.NewStyle().Bold(true)
	fmt.Println(title_running.Render("YouCanHack tip:"), " ")

	time.Sleep(2 * time.Second)

	core.Running("*")
}
