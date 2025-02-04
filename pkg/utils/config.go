package utils

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func (u *Utils) ViperReadConfig(env, appName, configFile string) {
	if len(env) > 0 {
		env = fmt.Sprintf(".%s", env)
	} else {
		if env == "prod" || env == "production" {
			env = ""
		} else {
			env = ".local"
		}
	}

	if configFile == "" {
		slog.Info("Reading config from application environment")
		slog.Debug("config file is empty, will try to use defaults...")
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nomnom" (without extension).
		//viper.AddConfigPath(home)
		//viper.SetConfigName(".nomnom")
		configFile = fmt.Sprintf("%s/.%s.yaml", home, appName)
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		slog.Error("bitebuddy.main.readConfig: error getting current directory. %s", cwdErr.Error())
		execPath, execPathErr := os.Executable()
		if execPathErr != nil {
			log.Fatalf("bitebuddy.main.readConfig: error getting path to the executable of the current process: %s\n", execPathErr.Error())
		}
		cwd = filepath.Dir(execPath)
	}

	if !u.FileExists(configFile) {
		slog.Error(fmt.Sprintf("bitebuddy.main.readConfig: configuration file: %s does not exist", configFile))
		configFile = fmt.Sprintf("%s/configs/app%s.yml", cwd, env)
		slog.Debug(fmt.Sprintf("bitebuddy.main.readConfig: will try to use configuration file %s", configFile))
	}

	slog.Info(fmt.Sprintf("bitebuddy.main.readConfig: loading configuration from file: %s", configFile))
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if u.FileExists(configFile) {
		// If a config file is found, read it in.
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("bitebuddy.main.readConfig: error reading config file %s. %s\n\n", configFile, err.Error())
		}

		viper.AutomaticEnv() // read in environment variables that match

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			slog.Info(fmt.Sprintf("bitebuddy.main.readConfig: config file changed: %s. Reloading...", e.Name))
			if err := viper.ReadInConfig(); err != nil {
				slog.Error(fmt.Sprintf("bitebuddy.main.readConfig: error reloading viper configuration file at %s: %s", configFile, err.Error()))
			} else {
				slog.Debug(fmt.Sprintf("bitebuddy.main.readConfig: successfully reloaded configuration file at: %s", configFile))
			}
		})

		slog.Debug(fmt.Sprintf("bitebuddy.main.readConfig: current working directory = %s", cwd))
		slog.Debug(fmt.Sprintf("bitebuddy.main.readConfig: using config file: %s", viper.ConfigFileUsed()))
		slog.Info(fmt.Sprintf("bitebuddy.main.readConfig: starting \"%s\" on \"%s\" environment...", appName, env))
	} else {
		slog.Info(fmt.Sprintf("bitebuddy.main.readConfig: configuration file %s does not exist", configFile))
	}
}
