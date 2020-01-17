package cmd

import (
	"errors"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logrus.New()
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "nutanix-image",
	Version: "v1.0.0",
	Short:   "A tool for those who want to interact with the nutanix API",
	Long: `Nutanix-image enables the user to be able to interact with
the nutanix images API for the following:

  * downloading images
  * uploading images
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var endpoint []string
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nutanix-image.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "set verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "set debug output")
	rootCmd.PersistentFlags().StringSliceVarP(&endpoint, "endpoint", "", []string{}, "Nutanix endpoint. Can be a single endpoint, or multiple (e.g. --endpoint=\"nutanix.example.org\" --endpoint=\"nutanix1.example.org,nutanix2.example.org\")")
	rootCmd.PersistentFlags().String("username", "", "the username for Nutanix API authentication")
	rootCmd.PersistentFlags().String("password", "", "the password for Nutanix API authentication")
	rootCmd.PersistentFlags().Bool("insecure", false, "enable or disable server certificate validation")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
	})
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nutanix-image" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nutanix-image")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if viper.GetBool("debug") {
		log.Level = logrus.DebugLevel
	} else if viper.GetBool("verbose") {
		log.Level = logrus.InfoLevel
	} else {
		log.Level = logrus.ErrorLevel
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}

func validateArgs(args []string) error {
	if len(args) > 1 {
		return errors.New("Too many arguments. Only enter 1 image at a time.")
	}
	if len(args) < 1 {
		return errors.New("Not enough arguments. Need to enter at least 1")
	}
	return nil
}

func requireStringFlag(name string) {
	if viper.GetString(name) == "" {
		log.Fatalf("required flag(s) \"--%s\" not set\nset this flag or set %s in the config file\n", name, name)
	}
}

func requireSliceFlag(name string) {
	if len(viper.GetStringSlice(name)) == 0 {
		log.Fatalf("required flag(s) \"--%s\" not set\nset this flag or set %s in the config file\n", name, name)
	}
}