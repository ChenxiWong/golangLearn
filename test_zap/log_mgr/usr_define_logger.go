package log_mgr

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func FieldsTransfer2Map(fields []zapcore.Field) map[string]string {
	resMap := make(map[string]string, len(fields))
	for _, v := range fields {
		resMap[v.Key] = v.String
	}
	return resMap
}

func (c consoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := UserGet()

	// We don't want the entry's metadata to be quoted and escaped (if it's
	// encoded as strings), which means that we can't use the JSON encoder. The
	// simplest option is to use the memory encoder and fmt.Fprint.
	//
	// If this ever becomes a performance bottleneck, we can implement
	// ArrayEncoder for our plain-text format.
	arr := getSliceEncoder()
	// res := FieldsTransfer2Map(fields)
	// if traceId, ok := res["traceId"]; ok {
	// 	usr_str := fmt.Sprintf("%s|", traceId)
	// 	line.AppendString(usr_str)
	// }
	if len(fields) > 0 && c.MessageKey == "message" {
		filedsMap := fields[0].Interface.(map[string]string)
		if traceId, ok := filedsMap["traceId"]; ok {
			usr_str := fmt.Sprintf("%s|", traceId)
			line.AppendString(usr_str)
		}
		// 时间
		if c.TimeKey != "" && c.EncodeTime != nil {
			c.EncodeTime(ent.Time, arr)
		}

		// 日志级别标识
		if c.LevelKey != "" && c.EncodeLevel != nil {
			c.EncodeLevel(ent.Level, arr)
		}

		// 日志名称
		if ent.LoggerName != "" && c.NameKey != "" {
			nameEncoder := c.EncodeName

			if nameEncoder == nil {
				// Fall back to FullNameEncoder for backward compatibility.
				nameEncoder = zapcore.FullNameEncoder
			}

			nameEncoder(ent.LoggerName, arr)
		}

		// 调用位置信息
		if ent.Caller.Defined {
			if c.CallerKey != "" && c.EncodeCaller != nil {
				c.EncodeCaller(ent.Caller, arr)
			}
			if c.FunctionKey != "" {
				arr.AppendString(ent.Caller.Function)
			}
		}
		for i := range arr.elems {
			if i > 0 {
				line.AppendString(c.ConsoleSeparator)
			}
			fmt.Fprint(line, arr.elems[i])
		}
		putSliceEncoder(arr)

		// Add the message itself.
		if c.MessageKey != "" {
			c.addSeparatorIfNecessary(line)
			line.AppendString(ent.Message)
		}
	}

	// 自定义的字段
	// c.writeContext(line, fields)

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	line.AppendString(c.LineEnding)
	return line, nil
}

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

func GetLogger() *zap.Logger {
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
	encoder := NewConsoleEncoder(encoderConfig)
	// file, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// fileWriteSyncer := zapcore.AddSync(file)
	fileWriteSyncer := getLogWriter()
	// core := zapcore.NewTee(
	// 	// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是Debug 及以上，生产环境记得改成Info
	// 	// zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	// 	zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	// )
	core := zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel)
	opts := []zap.Option{}
	opts = append(opts, zap.AddCaller())
	log := zap.New(core, opts...)
	return log
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	/*
		Filename: 日志文件的位置
		MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups：保留旧文件的最大个数
		MaxAges：保留旧文件的最大天数
		Compress：是否压缩/归档旧文件
	*/
	return zapcore.AddSync(lumberJackLogger)
}
