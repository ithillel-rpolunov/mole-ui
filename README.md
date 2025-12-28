# Mole Wails

**Native macOS desktop application** for system cleanup and optimization, built with Wails and Vue.

## Overview

Mole Wails is a GUI wrapper around the powerful Mole CLI tool, providing an intuitive interface for:

- 🧹 **Deep System Cleanup** - Remove caches, logs, and temporary files
- 🗑️ **Smart App Uninstaller** - Remove apps with all related files
- ⚡ **System Optimization** - Rebuild caches and refresh services
- 📊 **Disk Space Analysis** - Visualize disk usage (coming soon)
- 📈 **System Monitoring** - Real-time metrics (coming soon)
- 🔥 **Project Purge** - Clean old build artifacts
- 🔐 **Touch ID Setup** - Configure Touch ID for sudo

## Quick Start

### Prerequisites

- macOS 11.0 or later
- Go 1.24+ installed
- Node.js 16+ installed
- Wails CLI installed (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Development

1. **Install dependencies**
   ```bash
   go mod download
   cd frontend && npm install && cd ..
   ```

2. **Run in development mode**
   ```bash
   wails dev
   ```

### Building

Build a production-ready .app bundle:

```bash
wails build
```

## Project Structure

```
mole-wails/
├── backend/               # Go backend
│   ├── services/         # Service layer (wraps bash scripts)
│   ├── models/           # Shared data structures
│   ├── analyze/          # Disk analyzer (from Mole)
│   └── status/           # System monitor (from Mole)
├── frontend/             # Vue 3 frontend
│   └── src/
│       ├── components/   # Vue components
│       └── stores/      # State management
├── scripts/              # Bash scripts from Mole
└── docs/                 # Documentation
```

## Features Status

### ✅ Implemented

- Project structure and architecture
- Backend Go services
- Frontend Vue components
- Clean tab with UI
- Wails bindings generated

### 🚧 Next Steps

- Complete bash script integration
- Implement remaining tabs
- Add Analyze and Status modules
- Error handling and validation
- Testing and polish

## Documentation

- [AGENT.md](AGENT.md) - AI agent development guidelines
- [Implementation Plan](docs/plans/2025-12-28-implementation-plan.md) - Complete feature roadmap

## License

MIT License
