/*
Package cmd

Copyright © 2022 dongfg <mail@dongfg.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/dongfg/sentry/internal"
	"github.com/spf13/pflag"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//  cfgFile global config file path
var cfgFile string

var (
	repo       string
	path       string
	username   string
	password   string
	privateKey string
)

var client *internal.GitClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sentry",
	Short: "Keep local repo UP-TO-DATE",
	Long:  `Keep local repo UP-TO-DATE with webhook`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		client = internal.New(&internal.GitOptions{
			Repo:       repo,
			Username:   username,
			Password:   password,
			PrivateKey: privateKey,
		})
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.sentry.yaml or .sentry.yaml)")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "http(s) or ssh repo url")
	_ = rootCmd.MarkPersistentFlagRequired("repo")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "t", "", "git clone path")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "git", "http credential username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "http credential password")
	rootCmd.PersistentFlags().StringVarP(&privateKey, "privateKey", "k", "", "ssh privateKey")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose mode")

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	_ = viper.BindPFlag("repo", rootCmd.PersistentFlags().Lookup("repo"))
	_ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("privateKey", rootCmd.PersistentFlags().Lookup("privateKey"))
}

func initConfig() {
	viper.SetEnvPrefix("SENTRY")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".notify" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sentry")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
		postInitCommands(rootCmd.Commands())
	} else {
		cobra.CheckErr(err)
	}
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	_ = viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			_ = cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
