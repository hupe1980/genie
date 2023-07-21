package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hupe1980/genie/pkg/codegen"
	"github.com/spf13/cobra"
)

// openAIOptions holds the command-line flags for the openAI command.
type openAIOptions struct {
	apiKey      string
	model       string
	temperature float32
	maxTokens   int
}

// newOpenAICmd creates a Cobra command for running codegen provided by OpenAI.
func newOpenAICmd(globalOpts *globalOptions) *cobra.Command {
	opts := &openAIOptions{}

	cmd := &cobra.Command{
		Use:           "openai",
		Short:         "Run codegen provided by openAI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `genie openai -p "Create a python hello world"
genie openai -p prompt.txt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := opts.apiKey
			if apiKey == "" {
				apiKey = os.Getenv("OPENAI_API_KEY")
			}

			cg, err := codegen.NewOpenAI(apiKey, func(o *codegen.OpenAIOptions) {
				o.ModelName = opts.model
				o.Temperature = opts.temperature
				o.MaxTokens = opts.maxTokens
			})
			if err != nil {
				return err
			}

			if err := run(context.Background(), globalOpts, cg); err != nil {
				return err
			}

			fmt.Println()
			fmt.Println(cg.Info())

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.apiKey, "api-key", "", "", "openAI api key")
	cmd.Flags().StringVarP(&opts.model, "model", "m", codegen.DefaultOpenAIModelName, "model to use")
	cmd.Flags().IntVarP(&opts.maxTokens, "max-tokens", "", codegen.DefaultOpenAIMaxTokens, "max tokens to use")
	cmd.Flags().Float32VarP(&opts.temperature, "temperature", "t", codegen.DefaultOpenAITemperature, "temperature to use")

	return cmd
}
