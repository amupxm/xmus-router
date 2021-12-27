package router

type (
	//LeveledLoggerInterface is the interface that defines leveled logger
	LeveledLoggerInterface interface {
		// Debugf logs a debug message using Printf conventions.
		Debugf(format string, v ...interface{})

		// Errorf logs a warning message using Printf conventions.
		Errorf(format string, v ...interface{})

		// Infof logs an informational message using Printf conventions.
		Infof(format string, v ...interface{})

		// Warnf logs a warning message using Printf conventions.
		Warnf(format string, v ...interface{})
	}
)
