# Totion 🧠

A minimalist, keyboard-centric terminal note-taking app. Built with **Go** and the **Charmbracelet TUI** stack (`bubbletea`, `bubbles`, `lipgloss`).

## 🛠 Features

* **File Management**: Create (`Ctrl+N`), List (`Ctrl+L`), and Delete (`Ctrl+D`) notes directly from the terminal.
* **Smart Editing**: A dedicated text area for writing with `Ctrl+S` to save.
* **Auto-Vault**: Automatically creates and manages a hidden storage directory at `~/.totion`.
* **Fuzzy Search**: Built-in filtering for your notes list.

## ⌨️ Keybindings

| Key | Action |
| --- | --- |
| `Ctrl + N` | Create a new note |
| `Ctrl + L` | List all notes |
| `Enter` | Open/Edit selected note |
| `Ctrl + S` | Save and close current note |
| `Ctrl + D` | Delete selected note |
| `Esc` | Back to main menu |
| `Ctrl + Q` | Quit |

## 🚀 Installation & Build

### Prerequisites

* Go 1.25+
* A terminal with Unicode support (recommended)

### Using Makefile

```bash
# Build and run immediately
make run

# Just build the binary
make build

```

## 📂 Project Structure

* **`main.go`**: Contains the Bubble Tea `Model`, `Update`, and `View` logic.
* **`~/.totion`**: The default "Vault" directory where your `.md` notes are stored.



