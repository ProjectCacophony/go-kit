package featureflag

import (
	"net/http"
	"time"

	unleash "github.com/Unleash/unleash-client-go/v3"
	"github.com/Unleash/unleash-client-go/v3/context"
)

// Config represents a Feature Flagger Configuration
type Config struct {
	Environment string `envconfig:"ENVIRONMENT"`

	UnleashURL        string `envconfig:"UNLEASH_URL"`
	UnleashInstanceID string `envconfig:"UNLEASH_INSTANCE_ID"`
}

// FeatureFlagger represents a service to check feature flags
type FeatureFlagger struct {
	unleashClient *unleash.Client
}

// New creates a new FeatureFlagger
func New(config *Config) (*FeatureFlagger, error) {
	if config.UnleashInstanceID == "" ||
		config.UnleashURL == "" {
		return &FeatureFlagger{}, nil
	}

	unleashClient, err := unleash.NewClient(
		unleash.WithUrl(config.UnleashURL),
		unleash.WithInstanceId(config.UnleashInstanceID),
		unleash.WithAppName(config.Environment),
		unleash.WithHttpClient(&http.Client{
			Timeout: 15 * time.Second,
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

// IsEnabled checks if a feature flag is enabled globally
func (ff *FeatureFlagger) IsEnabled(key string, fallback bool) bool {
	if ff.unleashClient == nil {
		return fallback
	}

	return ff.unleashClient.IsEnabled(key, unleash.WithFallback(fallback))
}

// IsEnabled checks if a feature flag is enabled for a specific UserID
func (ff *FeatureFlagger) IsEnabledFor(key string, fallback bool, userID string) bool {
	if ff.unleashClient == nil {
		return fallback
	}

	return ff.unleashClient.IsEnabled(key, unleash.WithFallback(fallback), unleash.WithContext(context.Context{
		UserId: userID,
	}))
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
