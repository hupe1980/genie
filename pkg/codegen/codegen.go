// Package codegen provides code generation functionality using different models.
package codegen

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"gopkg.in/yaml.v2"
)

//go:embed prompts/file_paths_system_prompt.tpl
var filePathsSystemPrompt string

//go:embed prompts/file_paths_human_prompt.tpl
var filePathsHumanPrompt string

//go:embed prompts/shared_dependencies_system_prompt.tpl
var sharedDependenciesSystemPrompt string

//go:embed prompts/code_generation_system_prompt.tpl
var codeGenerationSystemPrompt string

//go:embed prompts/code_generation_human_prompt.tpl
var codeGenerationHumanPrompt string

// CodeGenOptions contains options for configuring the CodeGen.
type CodeGenOptions struct{}

// CodeGen represents the interface for the code generation.
type CodeGen interface {
	// FilePaths generates file paths based on the input prompt.
	FilePaths(ctx context.Context, input *FilePathsInput) (*FilePathsOutput, error)

	// SharedDependencies generates shared dependencies based on the input prompt and file paths.
	SharedDependencies(ctx context.Context, input *SharedDependenciesInput) (*SharedDependenciesOutput, error)

	// GenerateSourceCode generates source code based on the input prompt, filename, file paths,
	// and shared dependencies.
	GenerateSourceCode(ctx context.Context, input *GenerateSourceCodeInput) (*GenerateSourceCodeOutput, error)
}

// codeGen represents a code generator based on a chat model.
type codeGen struct {
	model schema.ChatModel
	opts  CodeGenOptions
}

// New creates a new instance of the code generator.
func New(model schema.ChatModel, optFns ...func(o *CodeGenOptions)) CodeGen {
	opts := CodeGenOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &codeGen{
		model: model,
		opts:  opts,
	}
}

// FilePathsInput represents the input for generating file paths.
type FilePathsInput struct {
	Prompt string
}

// FilePathsOutput represents the output for generating file paths.
type FilePathsOutput struct {
	Reasoning []string `json:"reasoning"`
	FilePaths []string `json:"file_paths"`
}

// FilePaths generates a list of file paths based on the input prompt.
// It takes a FilePathsInput containing the prompt string as input and returns a FilePathsOutput
// containing the generated file paths and reasoning for the generated output.
// The function returns an error if there is any issue during the generation process.
func (cg *codeGen) FilePaths(ctx context.Context, input *FilePathsInput) (*FilePathsOutput, error) {
	ct := prompt.NewChatTemplate([]prompt.MessageTemplate{
		prompt.NewSystemMessageTemplate(filePathsSystemPrompt),
		prompt.NewHumanMessageTemplate(filePathsHumanPrompt),
	})

	pv, err := ct.FormatPrompt(map[string]any{
		"prompt": input.Prompt,
	})
	if err != nil {
		return nil, err
	}

	result, err := model.GeneratePrompt(ctx, cg.model, pv)
	if err != nil {
		return nil, err
	}

	output := FilePathsOutput{}
	if err = json.Unmarshal(toJSON(result.Generations[0].Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w\nRaw data: %v", err, result.Generations[0].Text)
	}

	return &output, nil
}

// SharedDependenciesInput represents the input for generating shared dependencies.
type SharedDependenciesInput struct {
	Prompt    string
	FilePaths []string
}

// SharedDependency represents a shared dependency with name, description, and symbols.
type SharedDependency struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Symbols     []string `json:"symbols"`
}

// SharedDependenciesOutput represents the output for generating shared dependencies.
type SharedDependenciesOutput struct {
	Reasoning          []string           `json:"reasoning"`
	SharedDependencies []SharedDependency `json:"shared_dependencies"`
}

// SharedDependencies generates shared dependencies based on the input prompt and file paths.
// It takes a SharedDependenciesInput containing the prompt string and file paths as input
// and returns a SharedDependenciesOutput containing the generated shared dependencies,
// along with reasoning for the generated output.
// The function returns an error if there is any issue during the generation process.
func (cg *codeGen) SharedDependencies(ctx context.Context, input *SharedDependenciesInput) (*SharedDependenciesOutput, error) {
	t := prompt.NewSystemMessageTemplate(sharedDependenciesSystemPrompt)

	pv, err := t.FormatPrompt(map[string]any{
		"prompt":    input.Prompt,
		"filePaths": input.FilePaths,
	})
	if err != nil {
		return nil, err
	}

	result, err := model.GeneratePrompt(ctx, cg.model, pv)
	if err != nil {
		return nil, err
	}

	output := SharedDependenciesOutput{}
	if err = json.Unmarshal(toJSON(result.Generations[0].Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w\nRaw data: %v", err, result.Generations[0].Text)
	}

	return &output, nil
}

// GenerateSourceCodeInput represents the input for generating source code.
type GenerateSourceCodeInput struct {
	Prompt             string
	Filename           string
	FilePaths          []string
	SharedDependencies []SharedDependency
}

// GenerateSourceCodeOutput represents the output for generating source code.
type GenerateSourceCodeOutput struct {
	Filename string `json:"filename"`
	Source   string `json:"source"`
}

// GenerateSourceCode generates source code based on the input prompt, filename, file paths,
// and shared dependencies.
// It takes a GenerateSourceCodeInput containing the prompt string, filename, file paths,
// and shared dependencies as input.
// The function returns a GenerateSourceCodeOutput containing the generated source code
// and the specified filename.
// The function returns an error if there is any issue during the generation process.
func (cg *codeGen) GenerateSourceCode(ctx context.Context, input *GenerateSourceCodeInput) (*GenerateSourceCodeOutput, error) {
	sharedDepsYaml, err := yaml.Marshal(input.SharedDependencies)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shared dependencies: %w", err)
	}

	ct := prompt.NewChatTemplate([]prompt.MessageTemplate{
		prompt.NewSystemMessageTemplate(codeGenerationSystemPrompt),
		prompt.NewHumanMessageTemplate(codeGenerationHumanPrompt),
	})

	pv, err := ct.FormatPrompt(map[string]any{
		"prompt":              input.Prompt,
		"filePaths":           input.FilePaths,
		"shared_dependencies": string(sharedDepsYaml),
		"filename":            input.Filename,
	})
	if err != nil {
		return nil, err
	}

	result, err := model.GeneratePrompt(ctx, cg.model, pv)
	if err != nil {
		return nil, err
	}

	return &GenerateSourceCodeOutput{
		Filename: input.Filename,
		Source:   result.Generations[0].Text,
	}, nil
}

// toJSON extracts a JSON string from a given string.
func toJSON(s string) []byte {
	re := regexp.MustCompile(`(?s)\{.*\}`)
	return re.Find([]byte(s))
}
