package pool

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.SugaredLogger

type Token struct {
	Id              int64  `json:"id"`
	ExpireTimestamp int64  `json:"expireTimestamp"`
	Sig             string `json:"nc_sig"`
	NcToken         string `json:"nc_token"`
	SessionId       string `json:"nc_csessionid"`
}

func (token *Token) toJson() ([]byte, error) {

	data, err := json.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("to json error ! %w", err)
	}

	return data, nil
}

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	file, _ := os.OpenFile("./token.log", os.O_CREATE|os.O_APPEND, 0666)
	writeSyncer := zapcore.AddSync(file)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	_logger := zap.New(core, zap.AddCaller())
	logger = _logger.Sugar()
}
