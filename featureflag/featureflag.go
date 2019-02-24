package featureflag

import (
	"net/http"
	"time"

	unleash "github.com/Unleash/unleash-client-go"
)

// FeatureFlaggerConfig represents a Feature Flagger Configuration
type FeatureFlaggerConfig struct {
	Environment string `envconfig:"ENVIRONMENT"`

	UnleashURL        string `envconfig:"UNLEASH_URL"`
	UnleashInstanceID string `envconfig:"UNLEASH_INSTANCE_ID"`
}

// FeatureFlagger represents a service to check feature flags
type FeatureFlagger struct {
	unleashClient *unleash.Client
}

// NewFeatureFlagger creates a new FeatureFlagger
func NewFeatureFlagger(config *FeatureFlaggerConfig) (*FeatureFlagger, error) {
	if config.Environment == "development" {
		return &FeatureFlagger{}, nil
	}

	unleashClient, err := unleash.NewClient(
		unleash.WithUrl(config.UnleashURL),
		unleash.WithInstanceId(config.UnleashInstanceID),
		unleash.WithAppName(config.Environment),
		unleash.WithHttpClient(&http.Client{
			Timeout: time.Second * 10,
		}),
		unleash.WithListener(&UnleashListener{}),
	)
	if err != nil {
		return nil, err
	}
	return &FeatureFlagger{
		unleashClient: unleashClient,
	}, nil
}

// IsEnabled checks if a feature flag is enabled
func (ff *FeatureFlagger) IsEnabled(key string, fallback bool) bool {
	if ff.unleashClient == nil {
		return fallback
	}

	return ff.unleashClient.IsEnabled(key, unleash.WithFallback(fallback))
}

// UnleashListener is our listener for Unleash events
type UnleashListener struct{}

// OnError logs errors
func (l UnleashListener) OnError(err error) {
}

// OnWarning logs warnings
func (l UnleashListener) OnWarning(warning error) {
}

// OnReady prints to the console when the repository is ready.
func (l UnleashListener) OnReady() {
}

// OnCount prints to the console when the feature is queried.
func (l UnleashListener) OnCount(name string, enabled bool) {
}

// OnSent prints to the console when the server has uploaded metrics.
func (l UnleashListener) OnSent(payload unleash.MetricsData) {
}

// OnRegistered prints to the console when the client has registered.
func (l UnleashListener) OnRegistered(payload unleash.ClientData) {
}
