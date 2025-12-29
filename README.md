# MoleUI

**Visual interface for the Mole CLI** - A native macOS desktop application for system cleanup and optimization.

## Overview

MoleUI provides an intuitive graphical interface to the powerful [Mole CLI tool](https://github.com/tw93/Mole), making system maintenance accessible through a modern desktop application built with Wails v2 and Vue 3.

### Features

- 🧹 **Deep System Cleanup** - Remove caches, logs, and temporary files to reclaim disk space
- 🗑️ **Smart App Uninstaller** - Completely remove applications with all related files
- ⚡ **System Optimization** - Rebuild caches, refresh services, and optimize system performance
- 📊 **Disk Space Analysis** - Visualize and analyze disk usage across directories
- 📈 **System Monitoring** - Real-time metrics for CPU, memory, disk, GPU, and network
- 🔐 **Touch ID Setup** - Configure Touch ID authentication for sudo commands

## Attribution

This application uses the [Mole project](https://github.com/tw93/Mole) by [@tw93](https://github.com/tw93) as its core engine. MoleUI provides a visual interface to make the powerful command-line tool more accessible to all users.

## Quick Start

### Prerequisites

- macOS 11.0 or later
- Go 1.24+ installed
- Node.js 16+ installed
- Wails CLI installed:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd mole-wails
   ```

2. **Install dependencies**
   ```bash
   go mod download
   cd frontend && npm install && cd ..
   ```

3. **Run in development mode**
   ```bash
   wails dev
   ```

   The app will launch with hot-reload enabled for both frontend and backend changes.

### Building

Build a production-ready .app bundle:

```bash
wails build
```

The compiled application will be available in `build/bin/`.

## Tech Stack

### Backend
- **Go 1.24+** - Backend logic and system integration
- **Wails v2** - Go + Web frontend framework for desktop apps
- **macOS APIs** - System metrics, disk analysis, Touch ID integration

### Frontend
- **Vue 3** - Composition API with `<script setup>` syntax
- **Pinia** - State management
- **TypeScript** - Type safety
- **Vite** - Fast build tool and dev server

### Architecture
- **Event-Driven** - Wails runtime events for real-time updates
- **Service Layer** - Clean separation between UI and system operations
- **Reactive State** - Vue 3 reactivity for seamless UI updates

## Project Structure

```
mole-wails/
├── backend/                 # Go backend
│   ├── services/           # Service layer
│   │   ├── clean.go       # System cleanup service
│   │   ├── uninstall.go   # App uninstaller service
│   │   ├── optimize.go    # System optimization service
│   │   └── touchid.go     # Touch ID configuration service
│   ├── models/             # Shared data structures
│   ├── analyze/            # Disk analyzer (from Mole)
│   └── status/             # System monitor (from Mole)
├── frontend/               # Vue 3 frontend
│   └── src/
│       ├── components/
│       │   ├── tabs/      # Main feature tabs
│       │   ├── layout/    # Layout components (Sidebar)
│       │   └── shared/    # Reusable components (Toast, ConfirmDialog)
│       └── stores/        # Pinia stores for state management
├── scripts/                # Bash scripts from Mole CLI
└── docs/                   # Documentation and plans
```

## Features Status

### ✅ Fully Implemented

- **Clean Tab** - Scan and remove system caches, logs, and temporary files with progress tracking
- **Uninstall Tab** - Remove applications and all associated files
- **Optimize Tab** - System optimization tasks with sudo support
- **Analyze Tab** - Disk space analysis with directory breakdown and large file detection
- **Status Tab** - Real-time system monitoring (CPU, Memory, Disk, GPU, Network, Battery)
- **Touch ID Tab** - Enable/disable Touch ID for sudo commands
- **About Tab** - Application information and credits
- **Custom Dialogs** - Native-looking confirmation dialogs (browser dialogs don't work in Wails)
- **Event System** - Real-time progress updates via Wails events
- **Error Handling** - Comprehensive error handling with user-friendly messages

### 🎨 UI/UX Features

- Modern dark theme with gradient accents
- Responsive layouts
- Loading states and progress bars
- Toast notifications
- Smooth animations and transitions
- Confirmation dialogs for destructive operations

## Development Notes

### Browser APIs Not Supported

Native browser dialogs (`alert()`, `confirm()`, `prompt()`) don't work in Wails desktop applications. Use custom Vue components instead:

- Use `ConfirmDialog.vue` for confirmations
- Use `Toast.vue` for notifications
- Display errors in the UI, not via `alert()`

### Working with Wails Events

```javascript
// Listen to events
EventsOn('event-name', (data) => {
  console.log('Received:', data)
})

// Emit events (from Go)
runtime.EventsEmit(ctx, "event-name", data)
```

### State Management Pattern

Each tab has its own Pinia store:
- `useCleanStore()` - Clean tab state
- `useUninstallStore()` - Uninstall tab state
- `useOptimizeStore()` - Optimize tab state
- `useAnalyzeStore()` - Analyze tab state
- `useStatusStore()` - Status tab state
- `useTouchIDStore()` - Touch ID tab state

## Documentation

- [Implementation Plan](docs/plans/2025-12-28-complete-mole-wails.md) - Complete feature roadmap
- [Mole CLI Repository](https://github.com/tw93/Mole) - Original CLI tool

## Known Issues

- File size calculations may show 0 for some cloud files (iCloud, sparse files) - uses logical size as fallback
- Disk metrics use primary disk only to avoid APFS container double-counting

## Contributing

Contributions are welcome! This project builds upon the excellent [Mole CLI](https://github.com/tw93/Mole) tool.

## License

MIT License

---

**Note**: This is a visual interface for the Mole CLI tool. All core functionality is provided by the Mole project. MoleUI simply makes it more accessible through a graphical interface.
