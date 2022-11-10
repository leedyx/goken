package pool

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var sugarLogger *zap.SugaredLogger

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	file, _ := os.OpenFile("./token.log", os.O_CREATE|os.O_APPEND, 0666)
	writeSyncer := zapcore.AddSync(file)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}
