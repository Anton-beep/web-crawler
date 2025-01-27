package config

import "go.uber.org/zap"

// InitLogger is a function that initializes the zap sugar logger
func InitLogger(debug bool) {
	cfg := zap.NewDevelopmentConfig()
	if debug {
		cfg.Level.SetLevel(zap.DebugLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	zap.S().Info("Start")
}
