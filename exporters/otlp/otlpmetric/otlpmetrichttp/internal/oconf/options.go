// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlpmetric/oconf/options.go.tmpl

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

package oconf // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/retry"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	// DefaultMaxAttempts describes how many times the driver
	// should retry the sending of the payload in case of a
	// retryable error.
	DefaultMaxAttempts int = 5
	// DefaultMetricsPath is a default URL path for endpoint that
	// receives metrics.
	DefaultMetricsPath string = "/v1/metrics"
	// DefaultBackoff is a default base backoff time used in the
	// exponential backoff strategy.
	DefaultBackoff time.Duration = 300 * time.Millisecond
	// DefaultTimeout is a default max waiting time for the backend to process
	// each span or metrics batch.
	DefaultTimeout time.Duration = 10 * time.Second
)

type (
	SignalConfig struct {
		Endpoint    string
		Insecure    bool
		TLSCfg      *tls.Config
		Headers     map[string]string
		Compression Compression
		Timeout     time.Duration
		URLPath     string

		TemporalitySelector metric.TemporalitySelector
		AggregationSelector metric.AggregationSelector
	}

	Config struct {
		// Signal specific configurations
		Metrics SignalConfig

		RetryConfig retry.Config
	}
)

// NewHTTPConfig returns a new Config with all settings applied from opts and
// any unset setting using the default HTTP config values.
func NewHTTPConfig(opts ...HTTPOption) Config {
	cfg := Config{
		Metrics: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorHTTPPort),
			URLPath:     DefaultMetricsPath,
			Compression: NoCompression,
			Timeout:     DefaultTimeout,

			TemporalitySelector: metric.DefaultTemporalitySelector,
			AggregationSelector: metric.DefaultAggregationSelector,
		},
		RetryConfig: retry.DefaultConfig,
	}
	cfg = ApplyHTTPEnvConfigs(cfg)
	for _, opt := range opts {
		cfg = opt.ApplyHTTPOption(cfg)
	}
	cfg.Metrics.URLPath = cleanPath(cfg.Metrics.URLPath, DefaultMetricsPath)
	return cfg
}

// cleanPath returns a path with all spaces trimmed and all redundancies
// removed. If urlPath is empty or cleaning it results in an empty string,
// defaultPath is returned instead.
func cleanPath(urlPath string, defaultPath string) string {
	tmp := path.Clean(strings.TrimSpace(urlPath))
	if tmp == "." {
		return defaultPath
	}
	if !path.IsAbs(tmp) {
		tmp = fmt.Sprintf("/%s", tmp)
	}
	return tmp
}

type (
	// GenericOption applies an option to the HTTP or gRPC driver.
	GenericOption interface {
		ApplyHTTPOption(Config) Config

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}

	// HTTPOption applies an option to the HTTP driver.
	HTTPOption interface {
		ApplyHTTPOption(Config) Config

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}
)

// genericOption is an option that applies the same logic
// for both gRPC and HTTP.
type genericOption struct {
	fn func(Config) Config
}

func (g *genericOption) ApplyHTTPOption(cfg Config) Config {
	return g.fn(cfg)
}

func (genericOption) private() {}

func newGenericOption(fn func(cfg Config) Config) GenericOption {
	return &genericOption{fn: fn}
}

// httpOption is an option that is only applied to the HTTP driver.
type httpOption struct {
	fn func(Config) Config
}

func (h *httpOption) ApplyHTTPOption(cfg Config) Config {
	return h.fn(cfg)
}

func (httpOption) private() {}

func NewHTTPOption(fn func(cfg Config) Config) HTTPOption {
	return &httpOption{fn: fn}
}

// Generic Options

func WithEndpoint(endpoint string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.Endpoint = endpoint
		return cfg
	})
}

func WithEndpointURL(v string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		u, err := url.Parse(v)
		if err != nil {
			global.Error(err, "otlpmetric: parse endpoint url", "url", v)
			return cfg
		}

		cfg.Metrics.Endpoint = u.Host
		cfg.Metrics.URLPath = u.Path
		if u.Scheme != "https" {
			cfg.Metrics.Insecure = true
		}

		return cfg
	})
}

func WithCompression(compression Compression) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.Compression = compression
		return cfg
	})
}

func WithURLPath(urlPath string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.URLPath = urlPath
		return cfg
	})
}

func WithRetry(rc retry.Config) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.RetryConfig = rc
		return cfg
	})
}

func WithTLSClientConfig(tlsCfg *tls.Config) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.TLSCfg = tlsCfg.Clone()
		return cfg
	})
}

func WithInsecure() GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.Insecure = true
		return cfg
	})
}

func WithSecure() GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.Insecure = false
		return cfg
	})
}

func WithHeaders(headers map[string]string) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.Headers = headers
		return cfg
	})
}

func WithTimeout(duration time.Duration) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.Timeout = duration
		return cfg
	})
}

func WithTemporalitySelector(selector metric.TemporalitySelector) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.TemporalitySelector = selector
		return cfg
	})
}

func WithAggregationSelector(selector metric.AggregationSelector) GenericOption {
	return newGenericOption(func(cfg Config) Config {
		cfg.Metrics.AggregationSelector = selector
		return cfg
	})
}
