// Package cmd provides a command-line interface for the Genie code generation tool.
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hupe1980/genie/pkg/codegen"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// Execute runs the root command and handles any errors that may occur.
func Execute(version string) {
	printLogo()

	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// globalOptions holds the command-line flags for Genie.
type globalOptions struct {
	apiKey      string
	prompt      string
	model       string
	temperature float32
	outdir      string
}

// newRootCmd creates the root Cobra command for Genie.
func newRootCmd(version string) *cobra.Command {
	globalOpts := &globalOptions{}

	cmd := &cobra.Command{
		Use:     "genie",
		Version: version,
		Short:   "Genie is a Proof of Concept (POC) source code generator that showcases the potential of utilizing Large Language Models (LLMs) for code generation.",
		Example: `genie -p "Create a python hello world"
genie -p prompt.txt`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := globalOpts.apiKey
			if apiKey == "" {
				apiKey = os.Getenv("OPENAI_API_KEY")
			}

			ctx := context.Background()

			cg, err := codegen.New(apiKey, func(o *codegen.CodeGenOptions) {
				o.ModelName = globalOpts.model
			})
			if err != nil {
				return err
			}

			prompt, err := readFileOrString(globalOpts.prompt)
			if err != nil {
				return err
			}

			fmt.Println("Create list of files:")

			fpOutput, err := cg.FilePaths(ctx, &codegen.FilePathsInput{
				Prompt: prompt,
			})
			if err != nil {
				return err
			}

			printListWithBullets(fpOutput.FilePaths)
			fmt.Println()

			fmt.Println("Reasoning:")
			printListWithBullets(fpOutput.Reasoning)
			fmt.Println()

			fmt.Println("Create list of shared Dependecies:")

			sdOutput, err := cg.SharedDependencies(ctx, &codegen.SharedDependenciesInput{
				Prompt:    prompt,
				FilePaths: fpOutput.FilePaths,
			})
			if err != nil {
				return err
			}

			printListWithBullets(sdOutput.SharedDependencies)
			fmt.Println()

			fmt.Println("Reasoning:")
			printListWithBullets(sdOutput.Reasoning)
			fmt.Println()

			g := new(errgroup.Group)
			g.SetLimit(3)

			for _, fp := range fpOutput.FilePaths {
				fp := fp

				filePath := filepath.Join(globalOpts.outdir, fp)

				// Check if file already exists:
				if _, statErr := os.Stat(filePath); statErr == nil {
					fmt.Printf("File %v already exists, skipping\n", filePath)
					continue
				}

				g.Go(func() error {
					cgOutput, genErr := cg.GenerateSourceCode(ctx, &codegen.GenerateSourceCodeInput{
						Prompt:             prompt,
						Filename:           fp,
						FilePaths:          fpOutput.FilePaths,
						SharedDependencies: sdOutput.SharedDependencies,
					})
					if genErr != nil {
						return genErr
					}

					// Create the folder if it doesn't exist
					folder := filepath.Dir(filePath)

					if _, statErr := os.Stat(folder); os.IsNotExist(statErr) {
						if mkErr := os.MkdirAll(folder, 0755); mkErr != nil {
							return mkErr
						}
					}

					// Write the data to the file.
					if wErr := os.WriteFile(filePath, []byte(cgOutput.Source), 0600); wErr != nil {
						return wErr
					}

					fmt.Printf("File %s created\n", fp)

					return nil
				})
			}

			err = g.Wait()
			if err != nil {
				return err
			}

			fmt.Println()
			fmt.Println(cg.Info())

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&globalOpts.apiKey, "api-key", "", "", "openAI api key")
	cmd.PersistentFlags().StringVarP(&globalOpts.prompt, "prompt", "p", "", "prompt to use (required)")
	cmd.PersistentFlags().StringVarP(&globalOpts.model, "model", "m", codegen.DefaultModelName, "model to use")
	cmd.PersistentFlags().StringVarP(&globalOpts.outdir, "outdir", "o", "dist", "outdir to use")
	cmd.PersistentFlags().Float32VarP(&globalOpts.temperature, "temperature", "t", codegen.DefaultTemperature, "temperature to use")

	_ = cmd.MarkPersistentFlagRequired("prompt")

	return cmd
}

// readFileOrString reads the contents of a file if the input is a valid filename,
// otherwise, it returns the input string as it is.
func readFileOrString(input string) (string, error) {
	if fileInfo, err := os.Stat(input); err == nil && !fileInfo.IsDir() {
		// Assuming it's a valid file path
		content, err := os.ReadFile(input)
		if err != nil {
			return "", err
		}

		return string(content), nil
	}

	return input, nil
}

// printListWithBullets prints a slice with bullet points.
func printListWithBullets[T any](slice []T) {
	for _, item := range slice {
		fmt.Printf("• %v\n", item)
	}
}

// printLogo prints the Genie logo.
func printLogo() {
	logo := ` ██████  ███████ ███    ██ ██ ███████ 
██       ██      ████   ██ ██ ██      
██   ███ █████   ██ ██  ██ ██ █████   
██    ██ ██      ██  ██ ██ ██ ██      
 ██████  ███████ ██   ████ ██ ███████`

	fmt.Println()
	fmt.Println(logo)
	fmt.Println()
}
