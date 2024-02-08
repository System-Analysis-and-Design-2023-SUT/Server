package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
)

// NOTE: To make logger usage more easy, we use singleton pattern
// which may not be a good pattern, by its a useful one.
// The codes here may seem a bit odd.
//var instance *zap.Logger
//var sugar *zap.SugaredLogger

var instance *zap.Logger
var sugar *Logger
var once sync.Once

// getOrCreateLogger is a singleton "Create of not created yet" method
// that creates the logger if it is not created yet.
// It supports logger naming, which is highly advised and two modes
// for development and production usage.
// Production mode creates the json output and submits logs
// with info level and upper.
// Development mode logs in more human friendly format and
// submits logs with all levels.
// In underlying logger there are more detailed differences,
// so if you need more power take a look at zap and zapcore
// and feel free to change the codes based on your needs.
//
// Example: logger.Info("A simple info", "key", "value", "foo", "bar")
func getOrCreateLogger(name string, isProduction bool) (*zap.Logger, error) {
	if instance == nil {
		var err error
		var config zap.Config
		if isProduction {
			// Create production logger configs
			encoder := zap.NewProductionEncoderConfig()
			// Disable the logger caller field
			// since it has no usage for us here.
			encoder.CallerKey = ""
			encoder.TimeKey = "timestamp"
			encoder.MessageKey = "message"
			config = zap.NewProductionConfig()
			// Disable error stack traces
			// since the error message should
			// be readable and traceable.
			config.OutputPaths = []string{"stdout"}
			config.DisableStacktrace = true
			config.EncoderConfig = encoder
		} else {
			// Create development logger configs
			encoder := zap.NewDevelopmentEncoderConfig()
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig = encoder
		}
		// Create the logger instance
		instance, err = config.Build()
		if err != nil {
			return nil, err
		}
		// Make the logger a sugar logger instance
		// with specified name.
		sugarLogger := instance.Sugar().Named(name)
		sugar = &Logger{
			SugaredLogger: sugarLogger,
			hostName:      "",
		}

		// Get the host name of the machine
		hostName, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		sugar.hostName = hostName

	}

	return instance, nil
}

// CheckAndCreateLogger does a thread-safe call to getOrCreateLogger
// and ensures that it's called only once.
func CheckAndCreateLogger(name string, isProduction bool) (*zap.Logger, error) {
	var err error
	once.Do(func() {
		instance, err = getOrCreateLogger(name, isProduction)
	})
	return instance, err
}

// Info logs with info level
func Info(msg string, keysAndValues ...interface{}) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.Info(msg, keysAndValues...)
}

// Debug logs with debug level
func Debug(msg string, keysAndValues ...interface{}) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.Debug(msg, keysAndValues...)
}

// Error logs with error level
func Error(msg string, keysAndValues ...interface{}) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.Error(msg, keysAndValues...)
}

// Warn logs with warn level
func Warn(msg string, keysAndValues ...interface{}) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.Warn(msg, keysAndValues...)
}

// Fatal logs with fatal level
func Fatal(msg string, keysAndValues ...interface{}) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.Fatal(msg, keysAndValues...)
}

//functions with 'S' at the end, are 2-3x faster than the above functions BUT they only get strings as arguments

// Info logs with info level
func InfoS(msg string, keysAndValues ...string) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.InfoS(msg, keysAndValues...)
}

// Debug logs with debug level
func DebugS(msg string, keysAndValues ...string) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.DebugS(msg, keysAndValues...)
}

// Error logs with error level
func ErrorS(msg string, keysAndValues ...string) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.ErrorS(msg, keysAndValues...)
}

// Warn logs with warn level
func WarnS(msg string, keysAndValues ...string) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.WarnS(msg, keysAndValues...)
}

// Fatal logs with fatal level
func FatalS(msg string, keysAndValues ...string) {
	_, _ = CheckAndCreateLogger("", true)
	sugar.FatalS(msg, keysAndValues...)
}
