# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Things is a terminal-based task management application built with Go. It uses the Bubble Tea TUI framework to provide an interactive interface for tracking tasks ("things"), with features like priority sorting, time tracking, search, and customizable types.

## Core Architecture

### Package Structure

- **main.go**: Entry point that resolves `THINGS_DIR` (defaults to `~/.things/`) and starts the UI
- **things/**: Data layer handling Things, Types, and time tracking
  - `things.go`: Main `Things` struct with filtering, sorting, and search logic
  - `thing.go`: Individual Thing model with Age() and Time() calculation methods
  - `type.go`: Type definitions loaded from markdown files
  - `timer.go`: CSV-based time tracking (records start/end timestamps)
- **ui/**: Bubble Tea TUI implementation following the Elm Architecture
  - `ui.go`: Model definition and core view logic
  - `init.go`: Bubble Tea Init command
  - `update.go`: Bubble Tea Update function with all keyboard handlers
  - `view.go`: Bubble Tea View rendering
  - `editor.go`: External editor integration (respects `$EDITOR`)

### Data Storage Model

Things uses a file-based storage approach in `THINGS_DIR` (default `~/.things/`):

- **things/**: Each Thing is a markdown file named with timestamp `YYYYMMDDHHMMSS.md`
- **types/**: Type definitions as markdown files (e.g., `chore.md`, `bug.md`)
- **time/**: CSV files tracking time spent editing each Thing

All files use frontmatter (YAML) for metadata and markdown for content.

### Thing Lifecycle

1. Created via `n` key → generates timestamped file with `priority: 5` default → opens in `$EDITOR`
2. Timer starts when Thing opens in editor (`things.Start()`)
3. Timer stops when editor closes (`things.Stop()`) → appends CSV row to time file
4. Things can be filtered (current/done/pause/today), sorted (priority/age/type), and searched
5. Things are grouped by priority (0-4 individual groups, 5+ grouped together) with blank lines as visual separators

### Key Architectural Patterns

- **Elm Architecture**: Model-Update-View cycle via Bubble Tea
- **Frontmatter Parsing**: Uses `github.com/adrg/frontmatter` to parse YAML metadata from markdown files
- **Time Tracking**: Automatic CSV logging of editor open/close times for calculating cumulative time spent
- **Dynamic Type System**: Types are user-defined markdown files loaded at runtime, not hardcoded

## Development Commands

### Building and Installing

```bash
make install          # Installs to $GOPATH/bin
go install .          # Alternative install method
go build .            # Build without installing
```

### Running

```bash
things                # Run with default ~/.things/ directory
THINGS_DIR=/path things  # Override things directory
EDITOR=nano things    # Override default editor (vim)
```

### Testing

No test files currently exist in the repository.

## Key Implementation Details

### Priority-Based Grouping and Sorting

The `Things.Sort()` method (things/things.go:136-186) implements priority-based grouping:
- Uses `priorityGroup()` helper to map priorities: 0-4 are individual groups, 5+ all map to group 5
- Groups things by priority level first (0 → 1 → 2 → 3 → 4 → 5+)
- Within each group, applies the selected sort mode (age/priority/type)
- Results in visual separation with blank lines between priority groups in the view

### Viewport Management

The viewport system (ui/ui.go and ui/view.go) dynamically manages scrolling and rendering:
- Viewport height set to terminal height minus 1 line (ui/update.go:17)
- Search input visibility reduces height by 3 when active (ui/ui.go:151-153)
- Cursor position kept in view via `setCursorInView()` which accounts for blank separator lines between priority groups
- `countSeparatorLines()` helper (ui/ui.go:105-126) counts blank lines in a range to accurately calculate visible items
- Rendering logic (ui/view.go:84-153) tracks `linesRendered` including both item lines and separator lines to maximize viewport usage

### Time Calculation

Time tracking (things/thing.go:45) reads CSV rows containing RFC3339 timestamps, calculates durations between start/end pairs, and sums them for total time spent.

### Editor Integration

When editing (ui/editor.go), the app:
1. Exits the TUI alternate screen
2. Launches `$EDITOR` (defaults to vim)
3. Returns to TUI when editor closes
4. Reloads Things/Types data to reflect changes

## Dependencies

Built with Go 1.22 and these key libraries:
- **github.com/charmbracelet/bubbletea**: TUI framework (Elm Architecture)
- **github.com/charmbracelet/lipgloss**: Styling and layout
- **github.com/charmbracelet/bubbles**: Reusable TUI components (textinput for search)
- **github.com/adrg/frontmatter**: YAML frontmatter parsing
