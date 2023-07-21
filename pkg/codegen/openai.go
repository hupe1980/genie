package codegen

import (
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
)

const (
	// DefaultOpenAIModelName is the default model name used for code generation.
	DefaultOpenAIModelName = "gpt-3.5-turbo"

	// DefaultOpenAITemperature is the default temperature used for code generation.
	DefaultOpenAITemperature = 0.4

	// DefaultOpenAIMaxTokens is the default maximum number of tokens to predict per generation.
	DefaultOpenAIMaxTokens = -1
)

// OpenAIOptions contains options for configuring the OpenAI.
type OpenAIOptions struct {
	// ModelName is the name of the model to use for code generation.
	ModelName string

	// Temperatur is the temperature setting for text generation. Higher values produce more random output.
	Temperature float32

	// MaxTokens denotes the number of tokens to predict per generation.
	MaxTokens int
}

// OpenAI represents a code generator using OpenAI models.
type OpenAI struct {
	CodeGen
	info *callback.OpenAIHandler
}

// NewOpenAI creates a new instance of OpenAI code generator with the provided API key.
func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := OpenAIOptions{
		ModelName:   DefaultOpenAIModelName,
		Temperature: DefaultOpenAITemperature,
		MaxTokens:   DefaultOpenAIMaxTokens,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	info := callback.NewOpenAIHandler()

	openAI, err := chatmodel.NewOpenAI(apiKey, func(o *chatmodel.OpenAIOptions) {
		o.ModelName = opts.ModelName
		o.Temperature = opts.Temperature
		o.MaxTokens = opts.MaxTokens
		o.Callbacks = []schema.Callback{info}
	})
	if err != nil {
		return nil, err
	}

	return &OpenAI{
		CodeGen: New(openAI),
		info:    info,
	}, nil
}

// Info returns the information about the code generator.
func (cg *OpenAI) Info() string {
	return cg.info.String()
}
