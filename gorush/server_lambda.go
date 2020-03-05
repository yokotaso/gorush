// +build lambda

package gorush

import (
	"context"
	"errors"
	"github.com/appleboy/com/file"
	"github.com/appleboy/gorush/config"
	"os"
	"plugin"
	"github.com/apex/gateway"
)

type ConfigCustomizer interface {
   Customize(config config.ConfYaml);
}

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer(ctx context.Context) error {
	if !PushConf.Core.Enabled {
		LogAccess.Debug("httpd server is disabled.")
		return nil
	}
	pluginPath := os.Getenv("GORUSH_CONFIG_PLUGIN_PATH")
	if pluginPath != "" && file.IsFile(pluginPath) {
		if err := loadConfigCustomizer(); err != nil {
				return err;
		}
	}

	LogAccess.Info("HTTPD server is running on " + PushConf.Core.Port + " port.")

	return gateway.ListenAndServe(PushConf.Core.Address+":"+PushConf.Core.Port, routerEngine())
}

func loadConfigCustomizer(pluginPath string) error {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		LogAccess.Debug(fmt.Sprintf("Failed to load %s", pluginPath))
		return err
	}

	symCustomizer, err := plug.Lookup("ConfigCustomizer")
	if err != nil {
		LogAccess.Debug("Failed to looup ConfigCustomizer")
		return err
	}

	var configCustomizer ConfigCustomizer
	configCustomizer, ok := symCustomizer.(ConfigCustomizer)
	if !ok {
		return errors.New("Unexpectecd type from module symbol")
	}
	configCustomizer.Customize(PushConf)
	return nil
}
