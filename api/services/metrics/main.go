package metrics

import (
	"time"

	"github.com/ardanlabs/conf/v3"
)

func main() {
	cfg := struct {
		conf.Version
		Web struct {
			DebugHost string `conf:"default:0.0.0.0:4010"`
		}
		Expvar struct {
			Host            string        `conf:"default:0.0.0.0:4000"`
			Route           string        `conf:"default:/metrics"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Prometheus struct {
			Host            string        `conf:"default:0.0.0.0:4020"`
			Route           string        `conf:"default:/metrics"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Collect struct {
			From string `conf:"default:http://localhost:3010/debug/vars"`
		}
		Publish struct {
			To       string        `conf:"default:console"`
			Interval time.Duration `conf:"default:5s"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

}
