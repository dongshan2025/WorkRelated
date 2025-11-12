package config

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapSugaredLogger returns a Logger adapter that wraps zap.SugaredLogger.
func NewZapSugaredLogger(sugar *zap.SugaredLogger) Logger {
	return zapLogger{su: sugar}
}

type zapLogger struct {
	su *zap.SugaredLogger
}

func (z zapLogger) Infof(format string, args ...any) {
	if z.su != nil {
		z.su.Infof(format, args...)
	}
}
func (z zapLogger) Errorf(format string, args ...any) {
	if z.su != nil {
		z.su.Errorf(format, args...)
	}
}
func (z zapLogger) Debugf(format string, args ...any) {
	if z.su != nil {
		z.su.Debugf(format, args...)
	}
}

// SetZapLogger is a convenience to set zap logger (non-sugared or sugared).
// If you pass a *zap.Logger it will be sugar'ed; if you pass nil it will be ignored.
func SetZapLogger(z *zap.Logger) {
	if z == nil {
		return
	}
	su := z.Sugar()
	SetLogger(NewZapSugaredLogger(su))
}

// SetupZapFromConfig builds a zap.Logger from viper-backed configuration under
// the `log` key and sets it as the package logger. It returns the built logger
// (caller should call Sync when appropriate) or an error.
//
// Supported config keys:
//
//	log.level (string) - e.g. "debug", "info", "warn", "error"
//	log.development (bool)
//	log.encoding (string) - "json" or "console"
//	log.output_paths ([]string)
//	log.error_output_paths ([]string)
func SetupZapFromConfig() (*zap.Logger, error) {
	if v == nil {
		return nil, nil
	}

	// start from production or development config
	var cfg zap.Config
	if v.GetBool("log.development") {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	if enc := v.GetString("log.encoding"); enc != "" {
		cfg.Encoding = enc
	}
	if outs := v.GetStringSlice("log.output_paths"); len(outs) > 0 {
		cfg.OutputPaths = outs
	}
	if errs := v.GetStringSlice("log.error_output_paths"); len(errs) > 0 {
		cfg.ErrorOutputPaths = errs
	}

	// encoder config overrides
	ec := cfg.EncoderConfig
	if tk := v.GetString("log.encoder.time_key"); tk != "" {
		ec.TimeKey = tk
	}
	if lk := v.GetString("log.encoder.level_key"); lk != "" {
		ec.LevelKey = lk
	}
	if nk := v.GetString("log.encoder.name_key"); nk != "" {
		ec.NameKey = nk
	}
	if ck := v.GetString("log.encoder.caller_key"); ck != "" {
		ec.CallerKey = ck
	}
	if mk := v.GetString("log.encoder.message_key"); mk != "" {
		ec.MessageKey = mk
	}
	if sk := v.GetString("log.encoder.stacktrace_key"); sk != "" {
		ec.StacktraceKey = sk
	}
	if le := v.GetString("log.encoder.line_ending"); le != "" {
		ec.LineEnding = le
	}

	// encode level
	switch strings.ToLower(v.GetString("log.encoder.level_encode")) {
	case "lowercase":
		ec.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "lowercase_color":
		ec.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "capital":
		ec.EncodeLevel = zapcore.CapitalLevelEncoder
	case "capital_color":
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "numeric":
		// numeric encoder not provided by zapcore; fall back to lowercase
		ec.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	// encode caller
	switch strings.ToLower(v.GetString("log.encoder.caller_encode")) {
	case "short":
		ec.EncodeCaller = zapcore.ShortCallerEncoder
	case "full":
		ec.EncodeCaller = zapcore.FullCallerEncoder
	}

	// encode name
	switch strings.ToLower(v.GetString("log.encoder.name_encode")) {
	case "short":
		ec.EncodeName = func(s string, enc zapcore.PrimitiveArrayEncoder) {
			// short name: last segment after dot
			parts := strings.Split(s, ".")
			enc.AppendString(parts[len(parts)-1])
		}
	case "full":
		ec.EncodeName = func(s string, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(s)
		}
	}

	// encode duration
	switch strings.ToLower(v.GetString("log.encoder.duration_encode")) {
	case "string":
		ec.EncodeDuration = zapcore.StringDurationEncoder
	case "nanos":
		ec.EncodeDuration = zapcore.NanosDurationEncoder
	case "millis":
		ec.EncodeDuration = zapcore.MillisDurationEncoder
	case "seconds":
		ec.EncodeDuration = zapcore.SecondsDurationEncoder
	}

	// encode time: support standard encoders or custom format
	if tf := v.GetString("log.encoder.time_format"); tf != "" {
		layout := tf
		ec.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(layout))
		}
	} else {
		switch strings.ToLower(v.GetString("log.encoder.time_encode")) {
		case "iso8601":
			ec.EncodeTime = zapcore.ISO8601TimeEncoder
		case "millis":
			ec.EncodeTime = zapcore.EpochMillisTimeEncoder
		case "nanos":
			ec.EncodeTime = zapcore.EpochNanosTimeEncoder
		case "seconds":
			ec.EncodeTime = zapcore.EpochTimeEncoder
		}
	}

	cfg.EncoderConfig = ec

	// sampling
	if v.GetBool("log.sampling.enabled") {
		sc := zap.SamplingConfig{}
		if v.IsSet("log.sampling.initial") {
			sc.Initial = v.GetInt("log.sampling.initial")
		}
		if v.IsSet("log.sampling.thereafter") {
			sc.Thereafter = v.GetInt("log.sampling.thereafter")
		}
		cfg.Sampling = &sc
	}

	if lvlStr := v.GetString("log.level"); lvlStr != "" {
		if lvl, err := zapcore.ParseLevel(lvlStr); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(lvl)
		} else {
			logger.Errorf("invalid log.level %q: %v, using default", lvlStr, err)
		}
	}

	l, err := cfg.Build()
	if err != nil {
		logger.Errorf("failed to build zap logger from config: %v", err)
		return nil, err
	}
	SetZapLogger(l)
	logger.Infof("zap logger configured from config: level=%s encoding=%s", cfg.Level.String(), cfg.Encoding)
	return l, nil
}
