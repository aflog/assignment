package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRoot = &cobra.Command{
	Use:          "assignment",
	Short:        "MessageBird assignment API",
	SilenceUsage: true,
}

// Execute will perform the execution of a given command
func Execute() {
	cobra.OnInitialize(initConfig)
	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}

// initConfig sets AutomaticEnv in viper to true.
func initConfig() {
	viper.AutomaticEnv()
}
