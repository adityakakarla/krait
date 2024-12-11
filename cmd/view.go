package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewOpenAIKeyCmd = &cobra.Command{
	Use:   "view",
	Short: "View the OpenAI API key",
	Run: func(cmd *cobra.Command, args []string) {
		openai := viper.GetString("openai.key")

		if openai == "" {
			fmt.Println("OpenAI API key not set")
		} else {
			fmt.Println("OpenAI API key:", openai)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewOpenAIKeyCmd)
}
