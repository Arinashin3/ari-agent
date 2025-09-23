package flag

import (
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	promslogflag "github.com/prometheus/common/promslog/flag"
)

var (
	ConfigFile *string
	Logger     *slog.Logger
)

func init() {
	ConfigFile = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.yml").String()
	promslogConfig := &promslog.Config{}
	promslogflag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	Logger = promslog.New(promslogConfig)
}
