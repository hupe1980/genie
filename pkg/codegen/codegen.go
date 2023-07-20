package codegen

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultModelName is the default model name used for code generation.
	DefaultModelName = "gpt-3.5-turbo"

	// DefaultTemperature is the default temperature used for code generation.
	DefaultTemperature = 0.4
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
type CodeGenOptions struct {
	// ModelName is the name of the model to use for code generation.
	ModelName string

	// Temperatur is the temperature setting for text generation. Higher values produce more random output.
	Temperature float32
}

// CodeGen represents a code generator based on a chat model.
type CodeGen struct {
	model schema.ChatModel
	info  *callback.OpenAIHandler
	opts  CodeGenOptions
}

// New creates a new CodeGen instance with the given API key and optional configuration options.
func New(apiKey string, optFns ...func(o *CodeGenOptions)) (*CodeGen, error) {
	opts := CodeGenOptions{
		ModelName:   DefaultModelName,
		Temperature: DefaultTemperature,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	info := callback.NewOpenAIHandler()

	openAI, err := chatmodel.NewOpenAI(apiKey, func(o *chatmodel.OpenAIOptions) {
		o.ModelName = opts.ModelName
		o.Temperature = opts.Temperature
		o.Callbacks = []schema.Callback{info}
	})
	if err != nil {
		return nil, err
	}

	return &CodeGen{
		model: openAI,
		info:  info,
		opts:  opts,
	}, nil
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
func (cg *CodeGen) FilePaths(input *FilePathsInput) (*FilePathsOutput, error) {
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

	result, err := model.GeneratePrompt(context.Background(), cg.model, pv)
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
func (cg *CodeGen) SharedDependencies(input *SharedDependenciesInput) (*SharedDependenciesOutput, error) {
	t := prompt.NewSystemMessageTemplate(sharedDependenciesSystemPrompt)

	pv, err := t.FormatPrompt(map[string]any{
		"prompt":    input.Prompt,
		"filePaths": input.FilePaths,
	})
	if err != nil {
		return nil, err
	}

	result, err := model.GeneratePrompt(context.Background(), cg.model, pv)
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

// GenerateSourceCode generates source code based on the input prompt, filename, file paths, and shared dependencies.
func (cg *CodeGen) GenerateSourceCode(input *GenerateSourceCodeInput) (*GenerateSourceCodeOutput, error) {
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

	result, err := model.GeneratePrompt(context.Background(), cg.model, pv)
	if err != nil {
		return nil, err
	}

	return &GenerateSourceCodeOutput{
		Filename: input.Filename,
		Source:   result.Generations[0].Text,
	}, nil
}

// Info returns the information about the code generator.
func (cg *CodeGen) Info() string {
	return cg.info.String()
}

// toJSON extracts a JSON string from a given string.
func toJSON(s string) []byte {
	re := regexp.MustCompile(`(?s)\{.*\}`)
	return re.Find([]byte(s))
}
