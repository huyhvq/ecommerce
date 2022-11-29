package api

import (
	"github.com/huyhvq/ecommerce/internal/database"
	"github.com/huyhvq/ecommerce/internal/server"
	"github.com/huyhvq/ecommerce/pkg/leveledlog"
	"github.com/huyhvq/ecommerce/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve APIs",
	Long:  "serve APIs server",
	Run:   serve,
}

func init() {
	ServeCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")
}

type application struct {
	config config
	db     *database.DB
	logger *leveledlog.Logger
}

func serve(cmd *cobra.Command, args []string) {
	cfg := newConfig(cfgFile)
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)
	db, err := database.New(cfg.DB.Driver, cfg.DB.DSN, cfg.DB.AutoMigrate)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
	}
	logger.Info("starting server on %s (version %s)", cfg.Addr, version.Get())

	if err := server.Run(cfg.Addr, app.routes()); err != nil {
		logger.Fatal(err)
	}
	logger.Info("server stopped")
}
