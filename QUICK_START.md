# 🎯 INICIO RÁPIDO: Arreglando `tkt`

## 👉 Empieza por aquí (orden de urgencia)

### PASO 1️⃣ - Crítico (30 minutos)
```bash
# 1. Quitar DEBUG PRINTS de cmd/tkt/start.go
#    → Busca todas las líneas con "DEBUG:"
#    → Reemplaza con utils.Dev.Debugf() o condición if isDebug

# 2. Implementar logger en internal/utils/logger.go
#    → Ver ANALYSIS.md para el código completo
#    → Asegúrate que Debugf, Infof, Warnf, Errorf funcionen
```

### PASO 2️⃣ - Alto (1 hora)
```bash
# 1. Agregar validación de tickets en start.go
#    → Verificar que el ticket existe ANTES de actualizar
#    → Verificar que ticketID > 0

# 2. Mejorar error handling en list.go y scan.go
#    → Usar utils.Dev.Warnf() en lugar de continue silencioso
```

### PASO 3️⃣ - Importante (2-3 horas)
```bash
# 1. Agregar tests básicos
#    → cmd/tkt/scan_test.go   (test regex parsing)
#    → cmd/tkt/start_test.go  (test con ticket simulado)
#    → internal/db/db_test.go (test schema y queries)

# 2. Crear wrapper bash/zsh
#    → ~/.local/bin/tkt (script que procesa CD: y BRANCH:)
#    → Permite que \"tkt start\" cambie de carpeta realmente
```

### PASO 4️⃣ - Medio (1-2 horas)
```bash
# 1. Sistema de migraciones para DB
#    → Reemplaza ALT
ER TABLE ad-hoc con migrations formales

# 2. Hacer config.yaml más flexible
#    → ignore_dirs configurable
#    → patterns configurables
#    → Validación en register.go
```

### PASO 5️⃣ - Bajo (30-60 minutos)
```bash
# 1. Agregar soporte de colores (opcional)
#    → go get github.com/fatih/color
#    → Reemplaza fmt.Printf con color.Green, color.Red, etc

# 2. Flags de verbosidad
#    → tkt scan --verbose
#    → tkt start --debug
```

---

## 📋 Checklist de Implementación

- [ ] Remover DEBUG prints (30 min)
- [ ] Implementar Logger (20 min)
- [ ] Validación de entrada (30 min)
- [ ] Mejorar error handling (30 min)
- [ ] Tests básicos (1-2 horas)
- [ ] Wrapper bash/zsh (20 min)
- [ ] Migraciones DB (30 min)
- [ ] Config flexible (30 min)
- [ ] Colores (15 min)
- [ ] Verbosity flags (20 min)

**Total: 3-4 horas de trabajo para una herramienta profesional**

---

## 🚀 Próximos Pasos Después del Cleanup

Una vez arreglado lo anterior:
1. Agregar `tkt watch` (watcher de cambios en real-time)
2. Integración con `tkt daily` (reporte de lo que hiciste)
3. Soporte para tiempo estimado vs tiempo real
4. Exportar datos (JSON, CSV para análisis)
5. Shell alias helper (auto-alias para tkt)

