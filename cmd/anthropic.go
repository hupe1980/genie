package cmd

import (
	"context"
	"os"

	"github.com/hupe1980/genie/pkg/codegen"
	"github.com/spf13/cobra"
)

// anthropicOptions holds the command-line flags for the anthropic command.
type anthropicOptions struct {
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
}

// newAnthropicCmd creates a Cobra command for running codegen provided by Anthropc.
func newAnthropicCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &anthropicOptions{}

	cmd := &cobra.Command{
		Use:           "anthropic",
		Short:         "Run codegen provided by anthropic",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `genie anthropic -p "Create a python hello world"
genie anthropic -p prompt.txt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := opts.apiKey
			if apiKey == "" {
				apiKey = os.Getenv("ANTHROPIC_API_KEY")
			}

			cg, err := codegen.NewAnthropic(apiKey, func(o *codegen.AnthropicOptions) {
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

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.apiKey, "api-key", "", "", "anthropic api key")
	cmd.Flags().StringVarP(&opts.model, "model", "m", codegen.DefaultAnthropicModelName, "model to use")
	cmd.Flags().IntVarP(&opts.maxTokens, "max-tokens", "", codegen.DefaultOpenAIMaxTokens, "max tokens to use")
	cmd.Flags().Float64VarP(&opts.temperature, "temperature", "t", codegen.DefaultAnthropicTemperature, "temperature to use")

	return cmd
}
