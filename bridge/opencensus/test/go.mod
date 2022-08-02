module go.opentelemetry.io/otel/bridge/opencensus/test

go 1.17

require (
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/bridge/opencensus v0.31.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.31.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/bridge/opencensus => ../

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace
