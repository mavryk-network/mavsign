package commands

import (
	"context"
	"os"
	"strings"

	"github.com/mavryk-network/mavsign/pkg/auth"
	"github.com/mavryk-network/mavsign/pkg/config"
	"github.com/mavryk-network/mavsign/pkg/metrics"
	"github.com/mavryk-network/mavsign/pkg/mavsign"
	"github.com/mavryk-network/mavsign/pkg/mavsign/watermark"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Context represents root command context shared with its children
type Context struct {
	Context context.Context

	config    *config.Config
	mavsign *mavsign.MavSign
}

// NewRootCommand returns new root command
func NewRootCommand(c *Context, name string) *cobra.Command {
	var (
		level      string
		configFile string
		baseDir    string
		jsonLog    bool
	)

	rootCmd := cobra.Command{
		Use:   name,
		Short: "A Mavryk Remote Signer for signing block-chain operations with private keys",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if cmd.Use == "version" ||
				strings.Contains(cmd.CommandPath(), "ledger") ||
				strings.Contains(cmd.CommandPath(), "list-requests") ||
				strings.Contains(cmd.CommandPath(), "list-ops") {
				return nil
			}

			// cmd always points to the top level command!!!
			conf := config.Default()
			if configFile != "" {
				if err := conf.Read(configFile); err != nil {
					return err
				}
			}

			if baseDir != "" {
				conf.BaseDir = baseDir
			}
			conf.BaseDir = os.ExpandEnv(conf.BaseDir)
			if err := os.MkdirAll(conf.BaseDir, 0770); err != nil {
				return err
			}

			validate := config.Validator()
			if err := validate.Struct(conf); err != nil {
				return err
			}

			if jsonLog {
				log.SetFormatter(&log.JSONFormatter{})
			}

			lv, err := log.ParseLevel(level)
			if err != nil {
				return err
			}

			log.SetLevel(lv)

			pol, err := mavsign.PreparePolicy(conf.Mavryk)
			if err != nil {
				return err
			}

			watermark, err := watermark.Registry().New(cmd.Context(), conf.Watermark.Driver, &conf.Watermark.Config, conf)
			if err != nil {
				return err
			}

			sigConf := mavsign.Config{
				Policy:      pol,
				Vaults:      conf.Vaults,
				Interceptor: metrics.Interceptor,
				Watermark:   watermark,
			}

			if conf.PolicyHook != nil && conf.PolicyHook.Address != "" {
				sigConf.PolicyHook = &mavsign.PolicyHook{
					Address: conf.PolicyHook.Address,
				}
				if conf.PolicyHook.AuthorizedKeys != nil {
					ak, err := auth.StaticAuthorizedKeys(conf.PolicyHook.AuthorizedKeys.List()...)
					if err != nil {
						return err
					}
					sigConf.PolicyHook.Auth = ak
				}
			}

			sig, err := mavsign.New(c.Context, &sigConf)
			if err != nil {
				return err
			}

			if err = sig.Unlock(c.Context); err != nil {
				return err
			}

			c.config = conf
			c.mavsign = sig
			return nil
		},
	}

	f := rootCmd.PersistentFlags()

	f.StringVarP(&configFile, "config", "c", "/etc/mavsign.yaml", "Config file path")
	f.StringVar(&level, "log", "info", "Log level: [error, warn, info, debug, trace]")
	f.StringVar(&baseDir, "base-dir", "", "Base directory. Takes priority over one specified in config")
	f.BoolVar(&jsonLog, "json-log", false, "Use JSON structured logs")

	return &rootCmd
}
