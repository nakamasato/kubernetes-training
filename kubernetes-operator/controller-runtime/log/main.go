package main

import (
	"errors"
	"flag"

	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	// NB: don't call SetLogger in init(), or else you'll mess up logging in the main suite.
	log = logf.Log.WithName("example")
)

func main() {
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	log.Error(errors.New("test error"), "error occurred", "check", 2)
	log.Info("log check", "key", 3)
}
