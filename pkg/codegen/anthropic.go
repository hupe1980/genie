package codegen

import (
	"github.com/hupe1980/golc/model/chatmodel"
)

const (
	// DefaultAnthropicModelName is the default model name used for code generation.
	DefaultAnthropicModelName = "claude-v1"

	// DefaultAnthropicTemperature is the default temperature used for code generation.
	DefaultAnthropicTemperature = 0.4

	// DefaultAnthropicMaxTokens is the default maximum number of tokens to predict per generation.
	DefaultAnthropicMaxTokens = 4000
)

// AnthropicOptions contains options for configuring the Anthropic model.
type AnthropicOptions struct {
	// ModelName is the name of the model to use for code generation.
	ModelName string

	// Temperatur is the temperature setting for text generation. Higher values produce more random output.
	Temperature float64

	// MaxTokens denotes the number of tokens to predict per generation.
	MaxTokens int
}

// Anthropic represents a code generator using Anthropic models.
type Anthropic struct {
	CodeGen
}

// NewAnthropic creates a new instance of Anthropic code generator with the provided API key.
func NewAnthropic(apiKey string, optFns ...func(o *AnthropicOptions)) (*Anthropic, error) {
	opts := AnthropicOptions{
		ModelName:   DefaultAnthropicModelName,
		Temperature: DefaultAnthropicTemperature,
		MaxTokens:   DefaultAnthropicMaxTokens,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	anthropic, err := chatmodel.NewAnthropic(apiKey, func(o *chatmodel.AnthropicOptions) {
		o.ModelName = opts.ModelName
		o.Temperature = opts.Temperature
		o.MaxTokens = opts.MaxTokens
	})
	if err != nil {
		return nil, err
	}

	return &Anthropic{
		CodeGen: New(anthropic),
	}, nil
}
