package config

import (
	"fmt"
	"github.com/vrischmann/envconfig"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var appConfig Config

type Config struct {
	Sentry struct {
		Dsn string `envconfig:"default=local"`
	}
	Image struct {
		Name string `envconfig:"default=image-name"`
		Tag  string `envconfig:"default=tag-release"`
	}
	Region      string `envconfig:"default=region"`
	Environment string `envconfig:"default=env"`
}

func InitConfig() error {
	appConfig = Config{}
	err := envconfig.Init(&appConfig)
	return err
}

// AppConfig returns the current AppConfig
func AppConfig() Config {
	return appConfig
}

// NewDynamicClient initialize k8s dynamic client use for k8s communication/operations
// Path param is path to Kubeconfig
// There is also TypedClient
func NewDynamicClient(path string) (dynamic.Interface, error) {
	conf, err := NewConfig(path)
	if err != nil {
		return nil, fmt.Errorf("dynamic client config creation failed, path: %s, error: %w", path, err)
	}

	dynC, err := dynamic.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("dynamic client creation failed, path: %s, error: %w", path, err)
	}
	return dynC, nil
}

func NewConfig(path string) (*rest.Config, error) {
	if len(path) > 0 {
		return clientcmd.BuildConfigFromFlags("", path)
	}
	return rest.InClusterConfig()
}
