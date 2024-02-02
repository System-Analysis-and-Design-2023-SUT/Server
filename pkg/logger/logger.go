package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	*zap.Logger
	hostName string
}

// this type here implements the ObjectMarshaler interface
// this has a better performance than using the "map[string]string" without implementing ObjectMarshaler
// because it prevents using reflect in zap library
type AdditionalInfo map[string]string

func (additionalInfo AdditionalInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range additionalInfo {
		enc.AddString(k, v)
	}
	return nil
}

func (logger *Logger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = logger.addAdditionalLogInfo(keysAndValues)
	logger.SugaredLogger.Infow(msg, keysAndValues...)
}

// Debug logs with debug level
func (logger *Logger) Debug(msg string, keysAndValues ...interface{}) {
	keysAndValues = logger.addAdditionalLogInfo(keysAndValues)
	logger.SugaredLogger.Debugw(msg, keysAndValues...)
}

// Error logs with error level
func (logger *Logger) Error(msg string, keysAndValues ...interface{}) {
	keysAndValues = logger.addAdditionalLogInfo(keysAndValues)
	logger.SugaredLogger.Errorw(msg, keysAndValues...)
}

// Warn logs with warn level
func (logger *Logger) Warn(msg string, keysAndValues ...interface{}) {
	keysAndValues = logger.addAdditionalLogInfo(keysAndValues)
	logger.SugaredLogger.Warnw(msg, keysAndValues...)
}

// Fatal logs with fatal level
func (logger *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	keysAndValues = logger.addAdditionalLogInfo(keysAndValues)
	logger.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (logger *Logger) addAdditionalLogInfo(keysAndValues []interface{}) []interface{} {
	var finalKeysAndValues []interface{}
	additionalInfo := AdditionalInfo{}
	// Add hostname info
	finalKeysAndValues = append(finalKeysAndValues, "hostname", logger.hostName)

	for it := 0; it < len(keysAndValues); it += 2 {
		additionalInfo[fmt.Sprintf("%v", keysAndValues[it])] = fmt.Sprintf("%v", keysAndValues[it+1])
	}
	finalKeysAndValues = append(finalKeysAndValues, "additional_info", additionalInfo)

	return finalKeysAndValues
}

//functions with 'S' at the end, are 2-3x faster than the above functions BUT they only get strings as arguments

func (logger *Logger) InfoS(msg string, keysAndValues ...string) {
	additionalInfo := logger.GetAdditionalLogInfo(keysAndValues...)
	logger.SugaredLogger.Infow(msg, "additional_info", additionalInfo, "hostname", logger.hostName)
}

// Debug logs with debug level
func (logger *Logger) DebugS(msg string, keysAndValues ...string) {
	additionalInfo := logger.GetAdditionalLogInfo(keysAndValues...)
	logger.SugaredLogger.Debugw(msg, "additional_info", additionalInfo, "hostname", logger.hostName)
}

// Error logs with error level
func (logger *Logger) ErrorS(msg string, keysAndValues ...string) {
	additionalInfo := logger.GetAdditionalLogInfo(keysAndValues...)
	logger.SugaredLogger.Errorw(msg, "additional_info", additionalInfo, "hostname", logger.hostName)
}

// Warn logs with warn level
func (logger *Logger) WarnS(msg string, keysAndValues ...string) {
	additionalInfo := logger.GetAdditionalLogInfo(keysAndValues...)
	logger.SugaredLogger.Warnw(msg, "additional_info", additionalInfo, "hostname", logger.hostName)
}

// Fatal logs with fatal level
func (logger *Logger) FatalS(msg string, keysAndValues ...string) {
	additionalInfo := logger.GetAdditionalLogInfo(keysAndValues...)
	logger.SugaredLogger.Fatalw(msg, "additional_info", additionalInfo, "hostname", logger.hostName)
}

// here we return the type 'AdditionalInfo' which implements the ObjectMarshaller
func (logger *Logger) GetAdditionalLogInfo(keysAndValues ...string) AdditionalInfo {

	additionalInfo := AdditionalInfo{}

	for it := 0; it < len(keysAndValues); it += 2 {
		additionalInfo[keysAndValues[it]] = keysAndValues[it+1]
	}

	return additionalInfo
}

func NewLogger(name string, isProduction bool) (*Logger, error) {
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
	instance, err := config.Build()
	if err != nil {
		return nil, err
	}
	// Make the logger a sugar logger instance
	// with specified name.
	sugar := instance.Sugar().Named(name)

	// Get the host name of the machine
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Logger{
		SugaredLogger: sugar,
		Logger:        instance,
		hostName:      hostName,
	}, nil
}

func NewLoggerWithSampler(name string, isProduction bool, duration time.Duration, firt int, thereafter int) (*Logger, error) {
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
	instance, err := config.Build()
	if err != nil {
		return nil, err
	}
	instance = zap.New(
		zapcore.NewSamplerWithOptions(
			instance.Core(),
			duration, firt, thereafter))
	// Make the logger a sugar logger instance
	// with specified name.
	sugar := instance.Sugar().Named(name)

	// Get the host name of the machine
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Logger{
		SugaredLogger: sugar,
		Logger:        instance,
		hostName:      hostName,
	}, nil
}
