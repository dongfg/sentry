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
	"github.com/dongfg/sentry/internal"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	port          int64
	webhookSecret string
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Start http server to receive webhook",
	Run: func(cmd *cobra.Command, args []string) {
		client.Clone(path)
		client.SetWebhookSecret(webhookSecret)
		internal.Start(port, client)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.PersistentFlags().Int64VarP(&port, "port", "P", 8080, "server port")
	appCmd.PersistentFlags().StringVarP(&webhookSecret, "webhookSecret", "", "", "webhook secret")

	_ = viper.BindPFlag("port", appCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("webhookSecret", appCmd.PersistentFlags().Lookup("webhookSecret"))
}
