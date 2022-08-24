# [log (controller-runtime)](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/log)

Package log contains utilities for fetching a new logger when one is not already available.

1. This package contains a root logr.Logger Log.
1. The sub-package zap provides helpers for setting up logr backed by Zap (go.uber.org/zap).

## [logr](https://pkg.go.dev/github.com/go-logr/logr)



## [zap](https://pkg.go.dev/go.uber.org/zap)

1. Use `Zap.Options`

    ```go
	opts := zap.Options{
		Development: true,
		TimeEncoder:  zapcore.ISO8601TimeEncoder,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
    ```

1. Give arguments

    ```
    go run main.go -h
    Usage of /var/folders/c2/hjlk2kcn63s4kds9k2_ctdhc0000gp/T/go-build3844265895/b001/exe/main:
      -kubeconfig string
            Paths to a kubeconfig. Only required if out-of-cluster.
      -zap-devel
            Development Mode defaults(encoder=consoleEncoder,logLevel=Debug,stackTraceLevel=Warn). Production Mode defaults(encoder=jsonEncoder,logLevel=Info,stackTraceLevel=Error) (default true)
      -zap-encoder value
            Zap log encoding (one of 'json' or 'console')
      -zap-log-level value
            Zap Level to configure the verbosity of logging. Can be one of 'debug', 'info', 'error', or any integer value > 0 which corresponds to custom debug levels of increasing verbosity
      -zap-stacktrace-level value
            Zap Level at and above which stacktraces are captured (one of 'info', 'error', 'panic').
      -zap-time-encoding value
            Zap time encoding (one of 'epoch', 'millis', 'nano', 'iso8601', 'rfc3339' or 'rfc3339nano'). Defaults to 'epoch'.
    ```

    ```
    go run main.go -zap-time-encoding epoch
    ```
