/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/11/2020
*/

package util

import (
	"errors"
	"flag"
	"fmt"
	"github.com/w6d-io/ci-operator/internal/config"
	"go.uber.org/zap/zapcore"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"strconv"
	"strings"
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
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// outputFormatFlag contains structure for managing zap encoding
type outputFormatFlag struct {
	zapOptions *zap.Options
	value      string
}

func (o *outputFormatFlag) String() string {
	return o.value
}

func (o *outputFormatFlag) Set(flagValue string) error {
	val := strings.ToLower(flagValue)
	switch val {
	case "json":
		o.zapOptions.Encoder = zapcore.NewJSONEncoder(JsonEncoderConfig())
	case "text":
		o.zapOptions.Encoder = zapcore.NewConsoleEncoder(TextEncoderConfig())
	default:
		return fmt.Errorf("invalid \"%s\"", flagValue)
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

// levelFlag contains structure for managing zap level
type levelFlag struct {
	zapOptions *zap.Options
	value      string
}

func (l levelFlag) String() string {
	return l.value
}

func (l levelFlag) Set(flagValue string) error {
	level, validLevel := levelStrings[strings.ToLower(flagValue)]
	if !validLevel {
		logLevel, err := strconv.Atoi(flagValue)
		if err != nil {
			return fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
		if logLevel > 0 {
			intLevel := -1 * logLevel
			l.zapOptions.Level = zapcore.Level(int8(intLevel))
		} else {
			return fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
	} else {
		l.zapOptions.Level = level
	}
	l.value = flagValue
	return nil
}

type configFlag struct {
	value string
}

func (f configFlag) String() string {
	return f.value
}

func (f configFlag) Set(flagValue string) error {
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

	var outputFormat outputFormatFlag
	outputFormat.zapOptions = o
	fs.Var(&outputFormat, "log-format", "log encoding ( 'json' or 'text')")

	var level levelFlag
	level.zapOptions = o
	fs.Var(&level, "log-level", "log level verbosity. Can be 'debug', 'info', 'error', "+
		"or any integer value > 0 which corresponds to custom debug levels of increasing verbosity")

	var c configFlag
	fs.Var(&c, "config", "config file")
}
