package log_mgr

// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type consoleEncoder struct {
	*jsonEncoder
}

// NewConsoleEncoder creates an encoder whose output is designed for human -
// rather than machine - consumption. It serializes the core log entry data
// (message, level, timestamp, etc.) in a plain-text format and leaves the
// structured context as JSON.
//
// Note that although the console encoder doesn't use the keys specified in the
// encoder configuration, it will omit any element whose key is set to the empty
// string.
func NewConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	if cfg.ConsoleSeparator == "" {
		// Use a default delimiter of '\t' for backwards compatibility
		cfg.ConsoleSeparator = "\t"
	}
	return consoleEncoder{newJSONEncoder(cfg, true)}
}

func (c consoleEncoder) Clone() zapcore.Encoder {
	return consoleEncoder{c.jsonEncoder.Clone().(*jsonEncoder)}
}

// func (c consoleEncoder) writeContext(line *buffer.Buffer, extra []zapcore.Field) {
// 	context := c.jsonEncoder.Clone().(*jsonEncoder)
// 	defer func() {
// 		// putJSONEncoder assumes the buffer is still used, but we write out the buffer so
// 		// we can free it.
// 		context.buf.Free()
// 		putJSONEncoder(context)
// 	}()

// 	addFields(context, extra)
// 	context.closeOpenNamespaces()
// 	if context.buf.Len() == 0 {
// 		return
// 	}

// 	c.addSeparatorIfNecessary(line)
// 	line.AppendByte('{')
// 	line.Write(context.buf.Bytes())
// 	line.AppendByte('}')
// }

func (c consoleEncoder) addSeparatorIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendString(c.ConsoleSeparator)
	}
}
