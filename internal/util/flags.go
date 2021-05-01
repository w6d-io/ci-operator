/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 21/11/2020
*/

package util

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/w6d-io/ci-operator/internal/config"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// JsonEncoderConfig returns an opinionated EncoderConfig
func JsonEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// TextEncoderConfig returns an opinionated EncoderConfig
func TextEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    "",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// OutputFormatFlag contains structure for managing zap encoding
type OutputFormatFlag struct {
	ZapOptions *zap.Options
	value      string
}

func (o *OutputFormatFlag) String() string {
	return o.value
}

func (o *OutputFormatFlag) Set(flagValue string) error {
	flagValue = LookupEnvOrString("LOG_FORMAT", flagValue)
	val := strings.ToLower(flagValue)
	switch val {
	case "json":
		o.ZapOptions.Encoder = zapcore.NewJSONEncoder(JsonEncoderConfig())
	case "text":
		o.ZapOptions.Encoder = zapcore.NewConsoleEncoder(TextEncoderConfig())
	default:
		return fmt.Errorf(`invalid "%s"`, flagValue)
	}
	o.value = flagValue
	return nil
}

// levelStrings contains level string supported
var levelStrings = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
}

// LevelFlag contains structure for managing zap level
type LevelFlag struct {
	ZapOptions *zap.Options
	value      string
}

func (l LevelFlag) String() string {
	return l.value
}

func (l LevelFlag) Set(flagValue string) error {
	flagValue = LookupEnvOrString("LOG_LEVEL", flagValue)
	level, validLevel := levelStrings[strings.ToLower(flagValue)]
	if !validLevel {
		logLevel, err := strconv.Atoi(flagValue)
		if err != nil {
			return fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
		if logLevel > 0 {
			intLevel := -1 * logLevel
			l.ZapOptions.Level = zapcore.Level(int8(intLevel))
		} else {
			return fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
	} else {
		l.ZapOptions.Level = level
	}
	l.value = flagValue
	return nil
}

type ConfigFlag struct {
	value string
}

func (f ConfigFlag) String() string {
	return f.value
}

func (f ConfigFlag) Set(flagValue string) error {
	if flagValue == "" {
		return errors.New("config cannot be empty")
	}
	isFileExists := func(filename string) bool {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}
	if !isFileExists(flagValue) {
		return fmt.Errorf("file %s does not exist", flagValue)
	}
	if err := config.New(flagValue); err != nil {
		return fmt.Errorf("instanciate config returns %s", err)
	}
	f.value = flagValue
	return nil
}

// BindFlags custom flags
func BindFlags(o *zap.Options, fs *flag.FlagSet) {

	var outputFormat OutputFormatFlag
	outputFormat.ZapOptions = o
	fs.Var(&outputFormat, "log-format", "log encoding ( 'json' or 'text')")

	var level LevelFlag
	level.ZapOptions = o
	fs.Var(&level, "log-level", "log level verbosity. Can be 'debug', 'info', 'error', "+
		"or any integer value > 0 which corresponds to custom debug levels of increasing verbosity")

	var c ConfigFlag
	fs.Var(&c, "config", "config file")
}

// LookupEnvOrString adds the capability to get env instead of param flag
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// LookupEnvOrBool adds the capability to get env instead of param flag
func LookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		b, _ := strconv.ParseBool(val)
		return b
	}
	return defaultVal
}
