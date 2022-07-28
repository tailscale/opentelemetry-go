module go.opentelemetry.io/otel/exporters/stdout/stdoutmetric

go 1.17

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require (
	github.com/stretchr/testify v1.7.5
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/metric v0.31.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/sdk/metric v0.31.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/trace v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace
