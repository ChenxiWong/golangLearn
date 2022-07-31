package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func MyISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "2006-01-02T15:04:05.000", enc)
}

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

func init() {
	// encoderConfig := zap.NewProductionEncoderConfig()
	// go.uber.org/zap/zapcore.EncoderConfig
	// {MessageKey: "msg", LevelKey: "level", TimeKey: "ts",
	// NameKey: "logger", CallerKey: "caller", FunctionKey: "",
	// StacktraceKey: "stacktrace", SkipLineEnding: false,
	// LineEnding: "\n", EncodeLevel: go.uber.org/zap/zapcore.LowercaseLevelEncoder,
	// EncodeTime: main.MyISO8601TimeEncoder,
	// EncodeDuration: go.uber.org/zap/zapcore.SecondsDurationEncoder,
	// EncodeCaller: go.uber.org/zap/zapcore.ShortCallerEncoder,
	// EncodeName: nil, NewReflectedEncoder: nil, ConsoleSeparator: "|"}
	// json.Encoder
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "loglevel",
		TimeKey:             "time",
		NameKey:             "nameKey",
		CallerKey:           "caller",
		FunctionKey:         "functionKey",
		StacktraceKey:       "stacktraceKey",
		SkipLineEnding:      false,
		LineEnding:          "\n",
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          MyISO8601TimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          nil,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "|",
	}
	// 设置日志记录中时间的格式
	encoderConfig.EncodeTime = MyISO8601TimeEncoder
	// 日志Encoder 还是JSONEncoder，把日志行格式化成JSON格式的
	// encoder := zapcore.NewJSONEncoder(encoderConfig)
	// encoderConfig.ConsoleSeparator = "|"
	encoder := NewConsoleEncoder(encoderConfig)
	// encoder := zapcore.NewConsoleEncoder(encoderConfig)
	file, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	fileWriteSyncer := zapcore.AddSync(file)
	core := zapcore.NewTee(
		// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是Debug 及以上，生产环境记得改成Info
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	opts := []zap.Option{}
	opts = append(opts, zap.AddCaller())
	log = zap.New(core, opts...)
}
func main() {
	// log := zap.NewExample()
	// log, _ := zap.NewProduction()
	// log, _ := zap.New()
	// log, _ := zap.NewDevelopment()
	// log, _ := zap.NewExample()
	// log, _ := zap.NewStdLog()
	log.Info("this is info message with fileds", zap.String("traceId", "111111111"))
	// log.Panic("this is panic message")

	// 我的诉求
	// 定义日志表示符号
	// 日志分割

}
