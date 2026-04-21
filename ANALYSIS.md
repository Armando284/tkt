# 📋 Análisis de `tkt` - Herramienta de Tickets Personales en Go

## ✅ Fortalezas del Proyecto

### 1. **Arquitectura limpia y bien organizada**
- Separación clara entre `cmd/` (CLI), `internal/db`, `internal/config`, `internal/models`
- Patrón de Cobra bien implementado
- Estructura escalable para agregar nuevos comandos

### 2. **Buena intención del proyecto**
- Resuelve un problema real: integrar TODO comments con un sistema de tickets
- Combina git, time tracking y CLI de forma coherente
- Idea de "one-command workflow" es excelente

### 3. **Stack técnico sólido**
- Go 1.25 es moderno
- SQLite con WAL mode es eficiente
- Dependencias mínimas y confiables (Cobra, Viper, fsnotify)

---

## ⚠️ Problemas y Mejoras Necesarias

### 1. **DEBUG CODE EN PRODUCCIÓN** ⚠️ CRÍTICO
**Archivo:** `cmd/tkt/start.go`

```go
fmt.Println("=== DEBUG: start command started ===")
fmt.Printf("DEBUG: Modo directo con ID = %s\n", args[0])
// ... 10+ líneas de debug prints
```

**Problema:** Los prints de DEBUG están en código de producción
**Solución:**
```go
// Usar la interfaz Logger adecuadamente
utils.Dev.Debugf("start command started")  // Solo en modo debug

// O mejor aún, usar una variable de control global
if isDebug {
    fmt.Println("=== DEBUG: start command started ===")
}
```

---

### 2. **Logger incompleto** 🔧 IMPORTANTE
**Archivo:** `internal/utils/logger.go`

```go
func (l *devLogger) Debugf(format string, args ...any) {}  // ← NO HACE NADA
func (l *devLogger) Infof(format string, args ...any)  {}  // ← NO HACE NADA
func (l *devLogger) Warnf(format string, args ...any)  {}  // ← NO HACE NADA
```

**Problema:** Los métodos de logging no implementan nada
**Solución propuesta:**
```go
package utils

import (
    "fmt"
    "os"
)

type Logger interface {
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
}

type SimpleLogger struct {
    level string
}

var Dev Logger = &SimpleLogger{level: "DEBUG"}

func (l *SimpleLogger) Debugf(format string, args ...any) {
    fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
}

func (l *SimpleLogger) Infof(format string, args ...any) {
    fmt.Fprintf(os.Stderr, "[INFO] "+format+"\n", args...)
}

func (l *SimpleLogger) Warnf(format string, args ...any) {
    fmt.Fprintf(os.Stderr, "[WARN] "+format+"\n", args...)
}

func (l *SimpleLogger) Errorf(format string, args ...any) {
    fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", args...)
}
```

---

### 3. **Manejo de errores inconsistente**

**Problemas encontrados:**

a) En `list.go`:
```go
if err := rows.Scan(&id, &title, &status, &folder, &project); err != nil {
    continue  // ← Silencia el error sin registrarlo
}
```

b) En `scan.go`:
```go
if err := rows.Scan(&root, &name); err != nil {
    continue  // ← También ignora el error
}
```

**Solución:**
```go
if err := rows.Scan(&id, &title, &status, &folder, &project); err != nil {
    utils.Dev.Warnf("failed to scan row: %v", err)
    continue
}
```

---

### 4. **SQLite schema improvements** 🔧

**Problemas en `internal/db/db.go`:**

a) Uso incorrecto de `ALTER TABLE`:
```go
// Esto falla si ya existe la columna
if _, err := DB.Exec(`ALTER TABLE tickets ADD COLUMN current_start_ts TEXT;`); err != nil {
    if !strings.Contains(err.Error(), "duplicate column name") {
        return err
    }
}
```

**Mejor aproximación (Migration system):**
```go
func applyMigrations() error {
    migrations := []struct {
        name string
        sql  string
    }{
        {
            name: "001_initial_schema",
            sql:  initialSchema,
        },
        // Futuras migraciones aquí
    }
    
    for _, m := range migrations {
        _, err := DB.Exec(m.sql)
        if err != nil {
            return err
        }
    }
    return nil
}
```

---

### 5. **Falta validación de entrada**

**En `start.go`:**
```go
ticketID, err = strconv.Atoi(args[0])
if err != nil {
    return fmt.Errorf("ID inválido: %s", args[0])
}
```

**Mejorar:**
```go
if ticketID < 1 {
    return fmt.Errorf("ID debe ser mayor a 0, recibido: %d", ticketID)
}

// Validar que el ticket exista ANTES de actualizar
var exists bool
err = db.DB.QueryRow("SELECT 1 FROM tickets WHERE id = ?", ticketID).Scan(&exists)
if err != nil {
    return fmt.Errorf("Ticket #%d no encontrado", ticketID)
}
```

---

### 6. **Hardcoding de rutas y convenciones**

**En `scan.go` - regex hardcoded:**
```go
todoRegex := regexp.MustCompile(`(?im)(?:^|[\s])(?:\/\/|/\*|#)\s*(TODO|FIXME|HACK):\s*(.+?)(?:\s*\*/|$)`)
```

**Mejor:**
- Hacer esto configurable en `config.yaml`
- Permitir patrones personalizados por proyecto

**Directorios ignorados hardcoded:**
```go
if name == ".git" || name == "node_modules" || name == ".venv" || name == "venv" ||
    name == "dist" || name == "build" || name == "target" {
```

**Solución:**
```go
// En config.go
type Config struct {
    IgnoreDirs []string `mapstructure:"ignore_dirs"`
    TodoPatterns []string `mapstructure:"todo_patterns"`
    // ...
}

// En scan.go
func shouldIgnore(name string, config *Config) bool {
    for _, ignored := range config.IgnoreDirs {
        if name == ignored {
            return true
        }
    }
    return false
}
```

---

### 7. **Falta de validación en `register.go`**

```go
// No verifica si la ruta existe
abRoot, err := filepath.Abs(root)
if err != nil {
    return err
}

// ← Debería validar
if _, err := os.Stat(abRoot); os.IsNotExist(err) {
    return fmt.Errorf("ruta no existe: %s", abRoot)
}

// ← Debería verificar si es repositorio Git
gitDir := filepath.Join(absRoot, ".git")
if _, err := os.Stat(gitDir); os.IsNotExist(err) {
    return fmt.Errorf("no es un repositorio Git: %s", absRoot)
}
```

---

### 8. **Falta el wrapper bash/zsh para `tkt start`** ⚠️

El README menciona:
```bash
tkt start          # interactive list
```

Pero el código Go solo imprime:
```go
fmt.Println("CD:" + finalFolder)
fmt.Println("BRANCH:" + branch.String)
```

**Necesitas un script wrapper:**
```bash
#!/bin/bash
# ~/.local/bin/tkt (wrapper)

output=$(tkt-bin "$@")

# Procesar comandos especiales
if echo "$output" | grep -q "^CD:"; then
    cd_path=$(echo "$output" | grep "^CD:" | cut -d: -f2)
    [ -n "$cd_path" ] && cd "$cd_path"
fi

if echo "$output" | grep -q "^BRANCH:"; then
    branch=$(echo "$output" | grep "^BRANCH:" | cut -d: -f2)
    [ -n "$branch" ] && git checkout -b "$branch" 2>/dev/null || git checkout "$branch"
fi

echo "$output"
```

---

### 9. **Namespace conflicts en imports**

```go
import (
    "github.com/go-viper/mapstructure/v2 v2.4.0"  // Indirecto
)
```

Este paquete no se usa directamente. Revisar si es necesario.

---

### 10. **Tests missing** 🧪

No hay tests en el proyecto. Para un CLI tool es crítico tener tests para:
- Parsing de TODO comments
- Git operations
- Database operations
- Time tracking

**Estructura recomendada:**
```
cmd/tkt/
├── start.go
├── start_test.go    ← Tests
├── scan.go
└── scan_test.go
```

---

## 🚀 Mejoras Prácticas para el Flujo de Trabajo

### 1. **Agregar flag de verbosidad**
```bash
tkt scan --verbose     # Muestra debug logs
tkt start 5 --debug    # Modo debug
```

### 2. **Mejorar la salida**
Usar colores ANSI correctamente (considera usar `fatih/color`)
```go
import "github.com/fatih/color"

color.Green("✅ Ticket started: #%d\n", ticketID)
color.Yellow("⚠️  Branch already exists\n")
color.Red("❌ Error: %v\n", err)
```

### 3. **Agregar validación de Git en tiempo de ejecución**
```go
// Verificar que el repositorio tiene cambios no commiteados antes de cambiar de rama
func hasUncommittedChanges(repoPath string) (bool, error) {
    // Implementar con go-git o git CLI
}
```

### 4. **Agregar configuración por proyecto**
```yaml
# ~/.local/share/tkt/projects/myproject/config.yaml
patterns:
  - "TODO"
  - "FIXME"
  - "BUG:"
ignore_dirs:
  - "vendor"
  - "dist"
  - "coverage"
```

### 5. **Mejorar el Makefile**
```makefile
.PHONY: test coverage lint fmt

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

fmt:
	go fmt ./...
```

---

## 📝 Resumen de Cambios Prioritarios

| Prioridad | Tarea | Esfuerzo |
|-----------|-------|---------|
| 🔴 CRÍTICO | Remover debug prints de producción | 10 min |
| 🔴 CRÍTICO | Implementar Logger correctamente | 20 min |
| 🟠 ALTO | Agregar validación de entradas | 30 min |
| 🟠 ALTO | Mejorar manejo de errores | 30 min |
| 🟠 ALTO | Agregar tests básicos | 1-2 horas |
| 🟡 MEDIO | Sistema de migraciones para DB | 30 min |
| 🟡 MEDIO | Crear wrapper bash/zsh | 20 min |
| 🟢 BAJO | Agregar soporte de colores | 15 min |

---

## 🎯 Conclusión

El proyecto tiene **buena base arquitectónica** pero necesita:
1. ✅ Limpiar código de debug
2. ✅ Mejorar logging y error handling
3. ✅ Agregar tests
4. ✅ Validar inputs y estados

Una vez arreglado esto, será una herramienta **muy productiva** para tu flujo de trabajo.

