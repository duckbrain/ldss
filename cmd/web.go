package cmd

import (
	"github.com/duckbrain/ldss/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch the web server",
	Long:  `Launch the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetInt("port")
		web.Run(port, lang())
	},
}

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().Int("port", 1830, "The TCP port to run the server on")
	viper.BindPFlag("port", webCmd.Flags().Lookup("port"))

}
