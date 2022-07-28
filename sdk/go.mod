module go.opentelemetry.io/otel/sdk

go 1.17

replace go.opentelemetry.io/otel => ../

require (
	github.com/go-logr/logr v1.2.3
	github.com/google/go-cmp v0.5.8
	github.com/stretchr/testify v1.7.5
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ../trace
