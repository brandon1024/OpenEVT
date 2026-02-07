package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/brandon1024/cmder"
	"golang.org/x/sync/errgroup"

	"github.com/brandon1024/OpenEVT/internal/evt"
	"github.com/brandon1024/OpenEVT/internal/web"
)

const desc = `OpenEVT - Envertec EVT400/EVT800 Client

OpenEVT is tool that empowers you to directly communicate with your Envertec solar microinverter, giving you complete
ownership of your solar performance data - without being restricted to the Envertec APIs. OpenEVT connects directly
to your microinverter on your LAN and can be used to integrate and monitor your inverter in Home Assistant or
Prometheus.

OpenEVT is known to work with the following Envertec inverters:

  - EVT800B
  - EVT400R

To connect to your microinverter, you need:

  - the LAN address and port number (e.g. '192.0.2.1:14889'),
  - the inverter serial number (e.g. '31583078')

Before connecting to your inverter, you must set up the inverter following the instructions provided by the
manufacturer. The inverter must be connected to your LAN and must be configured in 'TCP-Server' mode in the 'Network
Parameter Settings'.

The inverter enters a low-power standby mode when there's no sunlight, so OpenEVT won't be able to connect during the
night.
`

const examples = `
# connect to inverter and listen on port 9090
openevt --addr 192.168.2.54:14889 --serial-number 31583078

# connect to inverter and listen on another port
openevt --addr 192.168.2.54:14889 --serial-number 31583078 --web.listen-address :8080
`

var (
	loggerLevel = new(slog.LevelVar)
)

var (
	cmd = &Command{
		BaseCommand: cmder.BaseCommand{
			CommandName: "openevt",
			Usage:       "openevt --addr <addr> --serial-number <num>",
			ShortHelp:   "Envertec EVT400/EVT800 Client",
			Help:        desc,
			Examples:    examples,
		},
	}
)

type Command struct {
	cmder.BaseCommand

	client evt.Client

	webListenAddress       string
	telemetryPath          string
	disableExporterMetrics bool
	reconnectInverval      time.Duration
}

func (c *Command) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.client.InverterID, "serial-number", "", "`serial` number of your microinverter (e.g. 31583078)")
	fs.Var(alias(fs.Lookup("serial-number"), "s"))
	fs.StringVar(&c.client.Address, "addr", "", "`address` and port of the microinverter (e.g. 192.0.2.1:14889)")
	fs.Var(alias(fs.Lookup("addr"), "a"))

	fs.DurationVar(&c.client.ReadTimeout, "poll-interval", time.Duration(0), "attempt to poll the inverter status more frequently than advertised")
	fs.DurationVar(&c.reconnectInverval, "reconnect-interval", time.Minute, "interval between connection attempts (e.g. 1m)")

	fs.StringVar(&c.webListenAddress, "web.listen-address", ":9090", "`address` on which to expose metrics")
	fs.StringVar(&c.telemetryPath, "web.telemetry-path", "/metrics", "`path` under which to expose metrics")
	fs.BoolVar(&c.disableExporterMetrics, "web.disable-exporter-metrics", false, "exclude metrics about the exporter itself (go_*)")

	fs.TextVar(loggerLevel, "log.level", new(slog.LevelVar), "log `level` (e.g. debug, info, warn, error)")
}

func (c *Command) Run(ctx context.Context, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unexpected args: %v", args)
	}
	if c.client.InverterID == "" {
		return fmt.Errorf("serial number required")
	}
	if c.client.Address == "" {
		return fmt.Errorf("inverter address required")
	}

	// setup logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: loggerLevel,
	})))

	grp, ctx := errgroup.WithContext(ctx)

	// launch inverter client
	grp.Go(func() error {
		return inverterConnect(ctx, &c.client, c.reconnectInverval)
	})

	// launch web server
	grp.Go(func() error {
		return web.ListenAndServe(ctx, c.webListenAddress, c.telemetryPath, c.disableExporterMetrics)
	})

	return grp.Wait()
}

func alias(flg *flag.Flag, name string) (flag.Value, string, string) {
	return flg.Value, name, flg.Usage
}
