package utils

// Logger es la interfaz que usarás en todo el código
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Dev es la implementación que se usa en desarrollo
var Dev Logger = &devLogger{}

type devLogger struct{}

func (l *devLogger) Debugf(format string, args ...any) {}
func (l *devLogger) Infof(format string, args ...any)  {}
func (l *devLogger) Warnf(format string, args ...any)  {}
func (l *devLogger) Errorf(format string, args ...any) { l.log("ERROR", format, args...) }

func (l *devLogger) log(level, format string, args ...any) {
	// En producción este archivo ni siquiera se compila
}