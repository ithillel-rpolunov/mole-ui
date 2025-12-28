# Mole Wails AI Agent Documentation

> **READ THIS FIRST**: This file serves as the single source of truth for any AI agent working on the Mole Wails repository. It defines architectural context, development workflows, and behavioral guidelines.

## 1. Philosophy & Guidelines

### Core Philosophy

- **Safety First**: Never risk user data. Always validate before executing destructive operations.
- **User Experience**: Prioritize intuitive UI/UX. The app should be self-explanatory.
- **Code Quality**: Write clean, maintainable code. Prefer clarity over cleverness.
- **Leverage Existing Code**: Reuse tested scripts from Mole project. Don't reinvent the wheel.

### Development Principles

- **Honor** in careful research, **Shame** in guessing APIs
- **Honor** in seeking confirmation, **Shame** in vague execution
- **Honor** in human verification, **Shame** in assuming business logic
- **Honor** in reusing existing code, **Shame** in creating unnecessary abstractions
- **Honor** in proactive testing, **Shame** in skipping validation
- **Honor** in following specifications, **Shame** in breaking architecture

### Quality Standards

- **English Only**: All code, comments, and documentation in English
- **Type Safety**: Use TypeScript in frontend, leverage Go's type system in backend
- **No Unnecessary Comments**: Code should be self-explanatory
- **Consistent Formatting**: Run `go fmt` for Go, use ESLint/Prettier for Vue

## 2. Project Identity

- **Name**: Mole Wails
- **Purpose**: Native macOS desktop application for system cleanup and optimization
- **Core Value**: User-friendly GUI wrapper around the powerful Mole CLI tool
- **Mechanism**:
  - **Frontend**: Vue 3 + TypeScript + Vite for reactive UI
  - **Backend**: Go services that orchestrate bash scripts and native Go modules
  - **Scripts**: Reused bash scripts from original Mole project
  - **Bridge**: Wails runtime for seamless Go-JavaScript communication

## 3. Technology Stack

### Frontend
- **Framework**: Vue 3 with Composition API
- **Language**: TypeScript
- **Build Tool**: Vite
- **State Management**: Pinia
- **Styling**: CSS with scoped styles (optional: Tailwind CSS)
- **UI Components**: Custom components + potential UI library (e.g., Element Plus, Naive UI)

### Backend
- **Language**: Go 1.21+
- **Framework**: Wails v2
- **Script Execution**: `os/exec` package
- **Concurrency**: Goroutines for async operations
- **Event System**: Wails Events for real-time updates

### Scripts
- **Shell**: Bash 3.2+ (macOS compatible)
- **Source**: Copied from original Mole project (`scripts/bin/`, `scripts/lib/`)

## 4. Repository Architecture

### Directory Structure

```
mole-wails/
├── frontend/                 # Vue 3 application
│   ├── src/
│   │   ├── components/
│   │   │   ├── tabs/        # CleanTab, UninstallTab, etc.
│   │   │   ├── shared/      # Reusable components (Button, ProgressBar, etc.)
│   │   │   └── layout/      # Sidebar, Header, Footer
│   │   ├── stores/          # Pinia stores for state management
│   │   ├── services/        # Wails runtime bindings
│   │   ├── types/           # TypeScript interfaces
│   │   ├── App.vue          # Root component
│   │   └── main.ts          # Entry point
│   ├── package.json
│   └── vite.config.ts
│
├── backend/                  # Go backend
│   ├── services/            # Service layer
│   │   ├── clean.go         # Executes scripts/bin/clean.sh
│   │   ├── uninstall.go     # Executes scripts/bin/uninstall.sh
│   │   ├── optimize.go      # Executes scripts/bin/optimize.sh
│   │   ├── purge.go         # Executes scripts/bin/purge.sh
│   │   └── touchid.go       # Executes scripts/bin/touchid.sh
│   ├── analyze/             # Disk analyzer (from Mole cmd/analyze)
│   │   ├── main.go          # Entry point (refactored for library use)
│   │   ├── scanner.go       # Directory scanning logic
│   │   ├── heap.go          # Large file tracking
│   │   └── ...
│   ├── status/              # System monitor (from Mole cmd/status)
│   │   ├── main.go          # Entry point (refactored for library use)
│   │   ├── metrics.go       # Metrics collection orchestration
│   │   ├── metrics_*.go     # Individual metric collectors
│   │   └── ...
│   └── models/              # Shared data structures
│       └── types.go
│
├── scripts/                  # Bash scripts from Mole
│   ├── bin/                 # Executable scripts
│   │   ├── clean.sh
│   │   ├── uninstall.sh
│   │   ├── optimize.sh
│   │   ├── purge.sh
│   │   └── touchid.sh
│   └── lib/                 # Script libraries
│       ├── core/            # Core utilities (logging, file ops)
│       ├── clean/           # Cleaning tasks
│       ├── ui/              # TUI components (not used in Wails)
│       └── ...
│
├── build/                    # Build configuration
│   ├── appicon.png
│   └── darwin/
│
├── app.go                    # Main Wails app struct
├── main.go                   # Application entry point
├── wails.json                # Wails configuration
├── go.mod
├── go.sum
├── AGENT.md                  # This file
└── README.md                 # User-facing documentation
```

## 5. Key Workflows

### Development

1. **Setup**:
   ```bash
   cd mole-wails
   wails dev
   ```
   This starts both Go backend and Vue frontend with hot-reload.

2. **Adding a New Feature**:
   - **Backend**: Create/edit service in `backend/services/`
   - **Frontend**: Create component in `frontend/src/components/tabs/`
   - **Integration**: Bind Go method to Vue via Wails runtime

3. **Testing**:
   - **Go**: `go test ./...`
   - **Frontend**: `cd frontend && npm test`
   - **E2E**: Manual testing in dev mode

### Building

```bash
wails build
```

Produces a native macOS `.app` bundle in `build/bin/` directory.

### Architecture Patterns

#### Backend Service Pattern

All backend services follow this pattern:

```go
type ServiceName struct {
    scriptsPath string
    // other fields
}

func NewServiceName(scriptsPath string) *ServiceName {
    return &ServiceName{scriptsPath: scriptsPath}
}

// Exposed to frontend
func (s *ServiceName) DoSomething(params Params) (Result, error) {
    // 1. Validate input
    // 2. Execute bash script or Go code
    // 3. Parse output
    // 4. Return structured result
}

// For long-running operations
func (s *ServiceName) DoSomethingAsync(params Params) error {
    go func() {
        // Emit progress events via Wails Events
        runtime.EventsEmit(ctx, "progress", data)
    }()
    return nil
}
```

#### Frontend Component Pattern

All Vue components follow this structure:

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useServiceStore } from '@/stores/service'

const store = useServiceStore()
const loading = ref(false)

async function handleAction() {
    loading.value = true
    try {
        await store.performAction()
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    // Initialize component
})
</script>

<template>
    <div class="component-name">
        <!-- Component content -->
    </div>
</template>

<style scoped>
/* Component styles */
</style>
```

## 6. Implementation Details

### Script Execution from Go

**CRITICAL**: Always use absolute paths when executing scripts.

```go
func (s *Service) executeScript(scriptName string, args ...string) (string, error) {
    scriptPath := filepath.Join(s.scriptsPath, "bin", scriptName)

    // Verify script exists
    if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
        return "", fmt.Errorf("script not found: %s", scriptPath)
    }

    cmd := exec.Command("/bin/bash", scriptPath)
    cmd.Args = append(cmd.Args, args...)

    // Capture output
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

### Real-time Progress Updates

For long-running operations, use Wails Events:

**Backend**:
```go
func (s *Service) CleanAsync(ctx context.Context) error {
    go func() {
        // ... do work
        runtime.EventsEmit(ctx, "clean:progress", ProgressData{
            Percent: 50,
            Message: "Cleaning caches...",
        })
        // ... more work
        runtime.EventsEmit(ctx, "clean:complete", ResultData{})
    }()
    return nil
}
```

**Frontend**:
```typescript
import { EventsOn } from '../../wailsjs/runtime/runtime'

EventsOn('clean:progress', (data: ProgressData) => {
    // Update UI
})

EventsOn('clean:complete', (data: ResultData) => {
    // Show completion message
})
```

### Adapting Analyze/Status Modules

The `analyze` and `status` modules were originally CLI tools with TUI. They need adaptation:

1. **Remove `main()` function**: Convert to library packages
2. **Extract core logic**: Create public functions/methods
3. **Return structured data**: Replace TUI rendering with data structures
4. **Add context support**: For cancellation and timeout

**Example** (analyze):
```go
// Original: renders TUI
// Adapted: returns data

func ScanDirectory(ctx context.Context, path string) (*ScanResult, error) {
    // Reuse scanning logic
    // Return structured data instead of rendering
    return &ScanResult{
        Entries: entries,
        TotalSize: totalSize,
    }, nil
}
```

## 7. Common AI Tasks

### Adding a New Tab/Feature

1. **Backend**:
   - Create service in `backend/services/new_feature.go`
   - Implement methods that call scripts or native Go code
   - Bind service to app in `app.go`

2. **Frontend**:
   - Create component in `frontend/src/components/tabs/NewFeatureTab.vue`
   - Create Pinia store in `frontend/src/stores/newFeature.ts`
   - Add tab to sidebar navigation
   - Call backend methods via Wails bindings

3. **Testing**:
   - Write unit tests for Go service
   - Test component in dev mode
   - Verify end-to-end flow

### Debugging Script Execution

If a bash script fails:

1. **Check script path**: Ensure absolute path is correct
2. **Verify permissions**: Scripts must be executable (`chmod +x`)
3. **Check dependencies**: Scripts depend on `lib/` directory structure
4. **Test manually**: Run script directly in terminal to see raw output
5. **Add logging**: Log command being executed and full output

### Styling Components

- Use scoped styles to avoid conflicts
- Follow existing color scheme (Purple #8B5CF6 primary)
- Ensure dark mode compatibility
- Use CSS variables for theming

### Handling Errors

- **Backend**: Return descriptive errors, don't panic
- **Frontend**: Show user-friendly error messages via toast/modal
- **Logging**: Log errors to console in dev, to file in production

## 8. Security Considerations

- **Sudo Operations**: Scripts may require sudo. Handle password prompts gracefully.
- **Path Validation**: Always validate user-provided paths before deletion.
- **Dry-run Mode**: Implement preview mode for destructive operations.
- **User Confirmation**: Require explicit confirmation for dangerous actions.

## 9. Performance Guidelines

- **Async Operations**: Use goroutines for long-running tasks
- **UI Responsiveness**: Never block the main thread
- **Progress Feedback**: Always show progress for operations > 2 seconds
- **Caching**: Cache scan results where appropriate (analyze)
- **Throttling**: Throttle rapid UI updates (status monitoring)

## 10. Next Steps for Development

1. ✅ Project structure created
2. ✅ Scripts and modules copied
3. ⏳ Refactor `analyze` and `status` for library use
4. ⏳ Implement backend services
5. ⏳ Build frontend components
6. ⏳ Integrate and test
7. ⏳ Polish UI/UX
8. ⏳ Build and distribute

---

**Remember**: This project stands on the shoulders of the excellent Mole CLI tool. Respect its architecture, reuse its battle-tested code, and focus on building an intuitive GUI wrapper.
