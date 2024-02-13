// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlptracehttp // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

import (
	"crypto/tls"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/otlpconfig"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/retry"
)

// Compression describes the compression used for payloads sent to the
// collector.
type Compression int

const (
	// NoCompression tells the driver to send payloads without
	// compression.
	NoCompression Compression = iota
	// GzipCompression tells the driver to send payloads after
	// compressing them with gzip.
	GzipCompression
)

// Option applies an option to the HTTP client.
type Option interface {
	applyHTTPOption(Config) Config
}

// RetryConfig defines configuration for retrying batches in case of export
// failure using an exponential backoff.
type RetryConfig retry.Config

func ApplyHTTPEnvConfigs(cfg Config) Config {
	opts := getOptionsFromEnv()
	for _, opt := range opts {
		cfg = opt.applyHTTPOption(cfg)
	}
	return cfg
}

// DefaultEnvOptionsReader is the default environments reader.
var DefaultEnvOptionsReader = envconfig.EnvOptionsReader{
	GetEnv:    os.Getenv,
	ReadFile:  ioutil.ReadFile,
	Namespace: "OTEL_EXPORTER_OTLP",
}

func getOptionsFromEnv() []Option {
	opts := []Option{}

	DefaultEnvOptionsReader.Apply(
		envconfig.WithURL("ENDPOINT", func(u *url.URL) {
			opts = append(opts, withEndpointScheme(u))
			opts = append(opts, option(func(cfg Config) Config {
				cfg.Traces.Endpoint = u.Host
				// For OTLP/HTTP endpoint URLs without a per-signal
				// configuration, the passed endpoint is used as a base URL
				// and the signals are sent to these paths relative to that.
				cfg.Traces.URLPath = path.Join(u.Path, DefaultTracesPath)
				return cfg
			}))
		}),
		envconfig.WithURL("TRACES_ENDPOINT", func(u *url.URL) {
			opts = append(opts, withEndpointScheme(u))
			opts = append(opts, option(func(cfg Config) Config {
				cfg.Traces.Endpoint = u.Host
				// For endpoint URLs for OTLP/HTTP per-signal variables, the
				// URL MUST be used as-is without any modification. The only
				// exception is that if an URL contains no path part, the root
				// path / MUST be used.
				path := u.Path
				if path == "" {
					path = "/"
				}
				cfg.Traces.URLPath = path
				return cfg
			}))
		}),
		envconfig.WithTLSConfig("CERTIFICATE", func(c *tls.Config) { opts = append(opts, WithTLSClientConfig(c)) }),
		envconfig.WithTLSConfig("TRACES_CERTIFICATE", func(c *tls.Config) { opts = append(opts, WithTLSClientConfig(c)) }),
		envconfig.WithHeaders("HEADERS", func(h map[string]string) { opts = append(opts, WithHeaders(h)) }),
		envconfig.WithHeaders("TRACES_HEADERS", func(h map[string]string) { opts = append(opts, WithHeaders(h)) }),
		WithEnvCompression("COMPRESSION", func(c Compression) { opts = append(opts, WithCompression(c)) }),
		WithEnvCompression("TRACES_COMPRESSION", func(c Compression) { opts = append(opts, WithCompression(c)) }),
		envconfig.WithDuration("TIMEOUT", func(d time.Duration) { opts = append(opts, WithTimeout(d)) }),
		envconfig.WithDuration("TRACES_TIMEOUT", func(d time.Duration) { opts = append(opts, WithTimeout(d)) }),
	)

	return opts
}

// WithEnvCompression retrieves the specified config and passes it to ConfigFn as a Compression.
func WithEnvCompression(n string, fn func(Compression)) func(e *envconfig.EnvOptionsReader) {
	return func(e *envconfig.EnvOptionsReader) {
		if v, ok := e.GetEnvValue(n); ok {
			cp := NoCompression
			if v == "gzip" {
				cp = GzipCompression
			}

			fn(cp)
		}
	}
}

type option func(cfg Config) Config

func (o option) applyHTTPOption(cfg Config) Config {
	return o(cfg)
}

func withEndpointScheme(u *url.URL) Option {
	switch strings.ToLower(u.Scheme) {
	case "http", "unix":
		return WithInsecure()
	default:
		return WithSecure()
	}
}

func WithSecure() Option {
	return option(func(cfg Config) Config {
		cfg.Traces.Insecure = false
		return cfg
	})
}

// WithEndpointURL sets the target endpoint URL the Exporter will connect to.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_METRICS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used. If both are set, OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
// will take precedence.
//
// If both this option and WithEndpointURL are used, the last used option will
// take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, "localhost:4317" will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithEndpoint(endpoint string) Option {
	return option(func(cfg Config) Config {
		cfg.Traces.Endpoint = endpoint
		return cfg
	})
}

// WithEndpoint sets the target endpoint URL the Exporter will connect to.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_METRICS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used. If both are set, OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
// will take precedence.
//
// If both this option and WithEndpoint are used, the last used option will
// take precedence.
//
// If an invalid URL is provided, the default value will be kept.
//
// By default, if an environment variable is not set, and this option is not
// passed, "localhost:4317" will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithEndpointURL(u string) Option {
	return wrappedOption{otlpconfig.WithEndpointURL(u)}
}

// WithCompression tells the driver to compress the sent data.
func WithCompression(compression Compression) Option {
	return option(func(cfg Config) Config {
		cfg.Traces.Compression = compression
		return cfg
	})
}

// WithURLPath allows one to override the default URL path used
// for sending traces. If unset, default ("/v1/traces") will be used.
func WithURLPath(urlPath string) Option {
	return option(func(cfg Config) Config {
		cfg.Traces.URLPath = urlPath
		return cfg
	})
}

// WithTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send payloads to the
// collector. Use it if you want to use a custom certificate.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return option(func(cfg Config) Config {
		cfg.Traces.TLSCfg = tlsCfg.Clone()
		return cfg
	})

}

// WithInsecure tells the driver to connect to the collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecure() Option {
	return option(func(cfg Config) Config {
		cfg.Traces.Insecure = true
		return cfg
	})
}

// WithHeaders allows one to tell the driver to send additional HTTP
// headers with the payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithHeaders(headers map[string]string) Option {
	return option(func(cfg Config) Config {
		cfg.Traces.Headers = headers
		return cfg
	})
}

// WithTimeout tells the driver the max waiting time for the backend to process
// each spans batch.  If unset, the default will be 10 seconds.
func WithTimeout(duration time.Duration) Option {
	return option(func(cfg Config) Config {
		cfg.Traces.Timeout = duration
		return cfg
	})
}

// WithRetry configures the retry policy for transient errors that may occurs
// when exporting traces. An exponential back-off algorithm is used to ensure
// endpoints are not overwhelmed with retries. If unset, the default retry
// policy will retry after 5 seconds and increase exponentially after each
// error for a total of 1 minute.
func WithRetry(rc RetryConfig) Option {
	return option(func(cfg Config) Config {
		cfg.RetryConfig = retry.Config(rc)
		return cfg
	})
}
