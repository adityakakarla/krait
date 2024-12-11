package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generate = &cobra.Command{
	Use:   "generate",
	Short: "Generate a CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		openAIKey := getOpenAIKey()

		if openAIKey == "" {
			return
		}

		var toolName string
		var toolDescription string

		fmt.Println("Enter the name of your tool: ")
		fmt.Scanln(&toolName)

		toolDescription = readFullLine("Enter a brief description of your tool: ")
		if !createMain(toolName) {
			return
		}

		rootContent := generateRootContent(toolName, toolDescription, openAIKey)
		if !createRoot(toolName, rootContent) {
			return
		}

		modContent := generateModContent(toolName, rootContent, openAIKey)
		if !createMod(toolName, modContent) {
			return
		}

		modules := generateModules(modContent, openAIKey)
		if !installModules(toolName, modules) {
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(generate)
}

func getOpenAIKey() string {
	openai := viper.GetString("openai.key")

	if openai == "" {
		fmt.Println("OpenAI API key not set")
		return ""
	}

	return openai
}

func createMain(toolName string) bool {
	mainContent := []byte(fmt.Sprintf(`package main

	import "%s/cmd"
	
	func main() {
		cmd.Execute()
	}
	
	`, toolName))

	osCmd := exec.Command("mkdir", toolName)
	if err := osCmd.Run(); err != nil {
		fmt.Println("Error: creating CLI folder")
		return false
	}

	err := os.WriteFile(toolName+"/main.go", mainContent, 0644)
	if err != nil {
		fmt.Println("Error: creating main.go", err)
		return false
	}

	return true
}

func createMod(toolName string, modContent string) bool {
	byteContent := []byte(modContent)

	err := os.WriteFile(toolName+"/go.mod", byteContent, 0644)
	if err != nil {
		fmt.Println("Error: creating go.mod", err)
		return false
	}
	return true
}

func createRoot(toolName string, rootContent string) bool {
	osCmd := exec.Command("mkdir", toolName+"/cmd")
	if err := osCmd.Run(); err != nil {
		fmt.Println("Error: making command directory", err)
	}
	byteContent := []byte(rootContent)

	err := os.WriteFile(toolName+"/cmd/root.go", byteContent, 0644)
	if err != nil {
		fmt.Println("Error: writing to root", err)
		return false
	}
	return true
}

func installModules(toolName string, modules []string) bool {
	var stderr bytes.Buffer

	err := os.Chdir(toolName)
	if err != nil {
		fmt.Printf("Error: changing directory to %s: %v\n", toolName, err)
		return false
	}

	for _, module := range modules {
		osCmd := exec.Command("go", "get", "-u", module)
		osCmd.Stderr = &stderr

		if err := osCmd.Run(); err != nil {
			fmt.Printf("Error: getting go package %v\n", err)
			fmt.Printf("Standard Error: %s\n", stderr.String())
			return false
		}
	}

	return true
}

func generateText(prompt string, openaiKey string) string {
	client := openai.NewClient(
		option.WithAPIKey(openaiKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		fmt.Println("Error generating text:", err)
		return ""
	}
	return chatCompletion.Choices[0].Message.Content
}

func generateRootContent(toolName string, toolDescription string, openAIKey string) string {
	return generateText(fmt.Sprintf(`I am building a command-line tool named %s.
Help me build this. You must only reply with code or you will be brutally punished.
No extra formatting at all (no backticks or quotation marks). Otherwise, the world will end.
The goal of %s is to %s.
Give me a root.go file that looks something like this:

package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "%s",
    Short: "%s is a sample command-line application",
    Long:  "%s is a sample command-line application built using Cobra.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Welcome to MyCLI! Use --help for more info.")
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}`, toolName, toolName, toolDescription, toolName, toolName, toolName), openAIKey)
}

func generateModContent(toolName string, rootContent string, openAIKey string) string {
	return generateText(fmt.Sprintf(`Here is my cmd/root.go file:

%s

Create a go.mod file based on this. Only respond with the code.
No extra formatting or you will be brutally punished. No backticks or quotes.
The module name must be %s.
Only include relevant module imports based on the code.`, rootContent, toolName), openAIKey)
}

func generateModules(modContent string, openAIKey string) []string {
	modulesRaw := generateText(fmt.Sprintf(`
Based on this go.mod file, give me the modules I need to install.
For instance, I do not need to install modules like "os", but I may need
to install modules like "github.com/example". 
Simply return the modules, separated by commas. No whitespace.

go.mod: %s`, modContent), openAIKey)

	modules := strings.Split(modulesRaw, ",")
	return modules
}

func readFullLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
