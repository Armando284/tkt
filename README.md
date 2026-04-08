# tkt

**A blazing-fast Go CLI for personal ticket management directly from your codebase.**

`tkt` bridges the gap between your code comments (`TODO`, `FIXME`, etc.) and real workflow:

- Automatically turns `TODO` comments into tickets
- One command to start working: switches to the right project folder, creates/checks out a git branch, marks the ticket as `in-progress`, and opens Neovim
- Automatically tracks time when you close Neovim
- Daily reports with time spent per ticket
- Supports **multiple projects** from your home directory (`~`)
- Lightweight file watcher (zero noticeable overhead) that detects new TODOs in real time

Built in **Go** for maximum performance and a single static binary.

## Features

- **Scan mode** — Parses your codebase for `TODO`/`FIXME`/`HACK` comments and creates tickets automatically
- **Multi-project support** — Register as many repositories as you want and list/search tickets across all of them from anywhere
- **Git integration** — Creates or switches to `tkt-XXXX` branches automatically
- **Time tracking** — Precise start/end timestamps + daily work reports
- **Neovim workflow** — `tkt start` → `cd` + git checkout + open Neovim → time recorded on exit
- **Efficient watcher** — Real-time detection of new TODOs without slowing down your system (using fsnotify)
- **SQLite backend** — Local, fast, zero-config database
- **Terminal-first** — Clean, colorful output and interactive selection when needed

## Quick Start

```bash
# 1. Build or download the binary
go install github.com/armando284/tkt@latest
# or
go build -o ~/bin/tkt .

# 2. Register your projects (do this once per repo)
cd /path/to/your/project
tkt register

# 3. Scan for existing TODOs
tkt scan

# 4. Start working on a ticket
tkt start          # interactive list
# or
tkt start 42       # by ticket ID
```

When you close Neovim, time is automatically recorded.

Other useful commands:
```bash
tkt list           # All tickets across projects
tkt daily          # What you worked on today + time
tkt done 42        # Mark ticket as done
tkt watch          # Start background watcher for new TODOs
```

## Why Go?

- Excellent performance when scanning many large projects
- Extremely lightweight and efficient file watcher
- Single static binary (no runtime dependencies)
- Great concurrency support for multi-project operations

## Tech Stack

- Go 1.23+
- `spf13/cobra` (CLI framework)
- `fsnotify/fsnotify` (file watching)
- `go-git/go-git` (optional pure-Go git operations)
- `modernc.org/sqlite` (embedded SQLite)
- Standard library for everything else

## Installation

### From source
```bash
git clone https://github.com/armando284/tkt.git
cd tkt
go install .
```

### Future: Pre-built binaries

(Will be added for Linux/WSL, macOS, etc.)

## Project Status

This tool is being built step-by-step as a learning + productivity project.  
Currently implementing core commands + multi-project support + watcher.

Contributions and ideas are welcome!

## License

MIT

---

Made with ❤️ for developers who live in the terminal and Neovim.
