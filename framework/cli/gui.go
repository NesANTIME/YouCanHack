package main

import (
	"fmt"
	"os"
	"strings"

	"youcanhack/framework/controller"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorPanel     = lipgloss.Color("#161b22")
	colorBorder    = lipgloss.Color("#30363d")
	colorAccent    = lipgloss.Color("#58a6ff")
	colorGreen     = lipgloss.Color("#3fb950")
	colorMuted     = lipgloss.Color("#8b949e")
	colorSelected  = lipgloss.Color("#1f6feb")
	colorText      = lipgloss.Color("#e6edf3")
	colorRun       = lipgloss.Color("#238636")
	colorRunHover  = lipgloss.Color("#2ea043")
	colorTag       = lipgloss.Color("#1f2d3d")
	colorSub       = lipgloss.Color("#21262d")
	colorSubAccent = lipgloss.Color("#f78166")
)

var menu []controller.MenuItem

func init() {
	menu = controller.Get_BDLab()
}

// ── Modelo ────────────────────────────────────────────────────────────────────

type model struct {
	parentCursor int
	subCursor    int // -1 = ninguna sub seleccionada
	expanded     int // -1 = ninguno expandido

	runFocused bool
	width      int
	height     int
	runOutput  map[string]string
}

func initialModel() model {
	return model{
		parentCursor: 0,
		subCursor:    -1,
		expanded:     -1,
		runOutput:    make(map[string]string),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) currentID() string {
	if m.subCursor >= 0 && m.expanded >= 0 {
		return menu[m.expanded].SubOptions[m.subCursor].ID
	}
	return menu[m.parentCursor].ID
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			m.runFocused = !m.runFocused

		case "esc":
			if m.subCursor >= 0 {
				m.subCursor = -1
				m.expanded = -1
				m.runFocused = false
			}

		case "up", "k":
			if m.runFocused {
				break
			}
			if m.subCursor >= 0 {
				if m.subCursor > 0 {
					m.subCursor--
				}
			} else {
				if m.parentCursor > 0 {
					m.parentCursor--
					if m.expanded != m.parentCursor {
						m.expanded = -1
						m.subCursor = -1
					}
				}
			}

		case "down", "j":
			if m.runFocused {
				break
			}
			if m.subCursor >= 0 {
				subs := menu[m.expanded].SubOptions
				if m.subCursor < len(subs)-1 {
					m.subCursor++
				}
			} else {
				if m.parentCursor < len(menu)-1 {
					m.parentCursor++
					if m.expanded != m.parentCursor {
						m.expanded = -1
						m.subCursor = -1
					}
				}
			}

		case "right", "l":
			if m.runFocused {
				break
			}
			if m.subCursor < 0 && len(menu[m.parentCursor].SubOptions) > 0 {
				m.expanded = m.parentCursor
				m.subCursor = 0
			}

		case "enter", " ":
			if m.runFocused {
				id := m.currentID()
				var label, tag string
				if m.subCursor >= 0 && m.expanded >= 0 {
					sub := menu[m.expanded].SubOptions[m.subCursor]
					label, tag = sub.Label, sub.Tag
				} else {
					item := menu[m.parentCursor]
					label, tag = item.Label, item.Tag
				}
				m.runOutput[id] = fmt.Sprintf(
					"▶  Ejecutando: %s\n\n$ %s\n\n✔  Completado exitosamente.", label, tag)
			} else {
				// Enter en menú expande submenú
				if m.subCursor < 0 && len(menu[m.parentCursor].SubOptions) > 0 {
					m.expanded = m.parentCursor
					m.subCursor = 0
				}
			}

		case "left", "h":
			if m.runFocused {
				break
			}
			if m.subCursor >= 0 {
				m.subCursor = -1
				m.expanded = -1
			}
		}
	}

	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.width == 0 {
		return "Cargando..."
	}

	totalH := m.height
	leftW := 32
	rightW := m.width - leftW - 3

	// ── Panel izquierdo ───────────────────────────────────────────────────────

	var rows []string

	menuHeader := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Width(leftW-2).
		Padding(0, 1).
		Render("⬡  COMANDOS")

	rows = append(rows, menuHeader)

	for i, item := range menu {
		isParentSel := i == m.parentCursor
		isExpanded := i == m.expanded

		arrow := " "
		if len(item.SubOptions) > 0 {
			if isExpanded {
				arrow = "▾"
			} else {
				arrow = "▸"
			}
		}

		iconStyle := lipgloss.NewStyle().Foreground(colorMuted)
		labelStyle := lipgloss.NewStyle().Foreground(colorText)

		if isParentSel && m.subCursor < 0 {
			iconStyle = iconStyle.Foreground(colorAccent)
			labelStyle = labelStyle.Foreground(colorAccent).Bold(true)
		}

		arrowStyle := lipgloss.NewStyle().Foreground(colorMuted)
		if isExpanded {
			arrowStyle = arrowStyle.Foreground(colorAccent)
		}

		rowContent := iconStyle.Render(item.Icon+" ") +
			labelStyle.Render(item.Label) +
			"  " +
			arrowStyle.Render(arrow)

		rowStyle := lipgloss.NewStyle().
			Width(leftW-2).
			Padding(0, 1)

		if isParentSel && m.subCursor < 0 && !m.runFocused {
			rowStyle = rowStyle.Background(colorSelected)
		} else if isParentSel && m.subCursor < 0 {
			rowStyle = rowStyle.Foreground(colorAccent)
		}

		rows = append(rows, rowStyle.Render(rowContent))

		// sub opciones expandidas
		if isExpanded {
			for j, sub := range item.SubOptions {
				isSubSel := j == m.subCursor
				isLast := j == len(item.SubOptions)-1

				bullet := "  ├─"
				if isLast {
					bullet = "  └─"
				}

				bulletColor := colorBorder
				labelColor := colorMuted
				if isSubSel {
					bulletColor = colorSubAccent
					labelColor = colorSubAccent
				}

				subRow := lipgloss.NewStyle().Foreground(bulletColor).Render(bullet) +
					" " +
					lipgloss.NewStyle().Foreground(labelColor).Bold(isSubSel).Render(sub.Label)

				subStyle := lipgloss.NewStyle().
					Width(leftW - 2).
					Background(colorSub)

				if isSubSel && !m.runFocused {
					subStyle = subStyle.Background(lipgloss.Color("#2d1f12"))
				}

				rows = append(rows, subStyle.Render(subRow))
			}
		}
	}

	navHint := lipgloss.NewStyle().
		Foreground(colorMuted).
		Width(leftW-2).
		Padding(1, 1, 0, 1).
		Render("↑↓ nav · → expandir · ← colapsar\nEsc volver · Tab Run · q salir")

	menuContent := strings.Join(rows, "\n")

	leftPanel := lipgloss.NewStyle().
		Width(leftW).
		Height(totalH-4).
		Background(colorPanel).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 0).
		Render(menuContent + "\n\n" + navHint)

	// ── Panel derecho ─────────────────────────────────────────────────────────

	var (
		title   string
		tag     string
		desc    string
		details []string
		runID   string
	)

	if m.subCursor >= 0 && m.expanded >= 0 {
		parent := menu[m.expanded]
		sub := parent.SubOptions[m.subCursor]
		title = parent.Icon + "  " + parent.Label + "  /  " + sub.Label
		tag = sub.Tag
		desc = sub.Description
		details = sub.Details
		runID = sub.ID
	} else {
		item := menu[m.parentCursor]
		title = item.Icon + "  " + item.Label
		tag = item.Tag
		desc = item.Description
		details = item.Details
		runID = item.ID
	}

	titleRendered := lipgloss.NewStyle().Foreground(colorText).Bold(true).Render(title)

	tagRendered := lipgloss.NewStyle().
		Background(colorTag).
		Foreground(colorAccent).
		Padding(0, 1).
		Bold(true).
		Render(tag)

	titleRow := lipgloss.JoinHorizontal(lipgloss.Center, titleRendered, "  ", tagRendered)

	divider := lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(strings.Repeat("─", rightW-4))

	descRendered := lipgloss.NewStyle().
		Foreground(colorMuted).
		Width(rightW - 4).
		Render(desc)

	var detailLines []string
	for _, d := range details {
		bullet := lipgloss.NewStyle().Foreground(colorAccent).Render("›")
		text := lipgloss.NewStyle().Foreground(colorText).Render(" " + d)
		detailLines = append(detailLines, bullet+text)
	}
	detailBlock := ""
	if len(detailLines) > 0 {
		detailBlock = "\n" + strings.Join(detailLines, "\n")
	}

	outputSection := ""
	if out, ok := m.runOutput[runID]; ok {
		outputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGreen).
			Foreground(colorGreen).
			Padding(1, 2).
			Width(rightW - 4).
			Render(out)
		outputSection = "\n\n" + outputBox
	}

	// Breadcrumb
	var breadcrumb string
	if m.subCursor >= 0 && m.expanded >= 0 {
		parent := menu[m.expanded]
		sub := parent.SubOptions[m.subCursor]
		breadcrumb = lipgloss.NewStyle().Foreground(colorMuted).
			Render(parent.Label + " › " + sub.Label)
	} else {
		breadcrumb = lipgloss.NewStyle().Foreground(colorMuted).
			Render(menu[m.parentCursor].Label)
	}

	rightContent := lipgloss.JoinVertical(
		lipgloss.Left,
		breadcrumb,
		"\n",
		titleRow,
		"\n",
		divider,
		"\n",
		descRendered,
		detailBlock,
		outputSection,
	)

	// Botón Run
	runLabel := "  ▶  RUN  "
	runStyle := lipgloss.NewStyle().
		Background(colorRun).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorGreen)

	if m.runFocused {
		runStyle = runStyle.
			Background(colorRunHover).
			BorderForeground(lipgloss.Color("#56d364"))
	}

	runBtn := runStyle.Render(runLabel)

	tabHint := ""
	if m.runFocused {
		tabHint = lipgloss.NewStyle().Foreground(colorGreen).Render("  ← Tab para volver")
	} else {
		tabHint = lipgloss.NewStyle().Foreground(colorMuted).Render("  Tab para enfocar")
	}

	runRow := lipgloss.JoinHorizontal(lipgloss.Center, runBtn, tabHint)

	contentH := lipgloss.Height(rightContent)
	availH := totalH - 8
	pad := availH - contentH
	if pad < 2 {
		pad = 2
	}

	rightInner := rightContent + strings.Repeat("\n", pad) + runRow

	rightPanel := lipgloss.NewStyle().
		Width(rightW).
		Height(totalH-4).
		Background(colorPanel).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2).
		Render(rightInner)

	// ── Header + layout ───────────────────────────────────────────────────────

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, " ", rightPanel)

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Render("  ◆ Go Task Runner"),
		lipgloss.NewStyle().Foreground(colorBorder).Render("  ·  "),
		lipgloss.NewStyle().Foreground(colorMuted).Render("q para salir"),
	)

	return header + "\n" + layout + "\n"
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
