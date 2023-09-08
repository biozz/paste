package cmd

import (
	"context"
	"fmt"

	"github.com/biozz/paste/internal/config"
	"github.com/biozz/paste/internal/server"
	"github.com/biozz/paste/internal/storage"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
)

// serverCmd represents the web command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a paste server",
	Run: func(cmd *cobra.Command, args []string) {
		var conf config.Server
		if err := envconfig.Process(context.Background(), &conf); err != nil {
			fmt.Printf("unable to initiate env config: %v", err)
			return
		}
		s, err := storage.New(conf.DSN)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		storageCloser, err := s.Init(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		defer storageCloser()
		w := server.New(conf, s)
		w.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
