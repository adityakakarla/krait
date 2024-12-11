package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openaiKey string

var setOpenAIKeyCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the OpenAI API key with flag -k",
	Run: func(cmd *cobra.Command, args []string) {
		if openaiKey == "" {
			fmt.Println("Error: OpenAI API key is required")
			return
		}

		viper.Set("openai.key", openaiKey)

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error finding home directory:", err)
			return
		}

		configPath := home + "/.krait.yaml"

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			file, err := os.Create(configPath)
			if err != nil {
				fmt.Println("Error creating config file:", err)
				return
			}
			file.Close()
		}

		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}

		fmt.Println("OpenAI API key has been set and saved.")
	},
}

func init() {
	rootCmd.AddCommand(setOpenAIKeyCmd)
	setOpenAIKeyCmd.Flags().StringVarP(&openaiKey, "key", "k", "", "OpenAI API key")
}
