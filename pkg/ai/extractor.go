package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

type AIExtractor struct {
	models map[string]Model
	config *AIConfig
}

type Model interface {
	Extract(ctx context.Context, input *ExtractionInput) (*ExtractionResult, error)
	Train(ctx context.Context, data *TrainingData) error
	Predict(ctx context.Context, features []float64) ([]float64, error)
}

type AIConfig struct {
	DefaultModel    string            `json:"default_model"`
	Models          map[string]ModelConfig `json:"models"`
	CacheEnabled    bool              `json:"cache_enabled"`
	CacheTTL        int               `json:"cache_ttl"`
	MaxTokens       int               `json:"max_tokens"`
	Temperature     float64           `json:"temperature"`
	Confidence      float64           `json:"confidence_threshold"`
}

type ModelConfig struct {
	Type       string                 `json:"type"`
	Endpoint   string                 `json:"endpoint"`
	APIKey     string                 `json:"api_key"`
	Parameters map[string]interface{} `json:"parameters"`
}

type ExtractionInput struct {
	HTML        string                 `json:"html"`
	URL         string                 `json:"url"`
	Schema      *ExtractionSchema      `json:"schema"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Options     *ExtractionOptions     `json:"options,omitempty"`
}

type ExtractionSchema struct {
	Fields      []FieldSchema          `json:"fields"`
	Validation  *ValidationRules       `json:"validation,omitempty"`
	PostProcess []PostProcessRule      `json:"post_process,omitempty"`
}

type FieldSchema struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Selector    string   `json:"selector,omitempty"`
	Attribute   string   `json:"attribute,omitempty"`
	Required    bool     `json:"required"`
	Multiple    bool     `json:"multiple"`
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

type ValidationRules struct {
	MinLength   int      `json:"min_length,omitempty"`
	MaxLength   int      `json:"max_length,omitempty"`
	Pattern     string   `json:"pattern,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
}

type PostProcessRule struct {
	Field     string `json:"field"`
	Operation string `json:"operation"`
	Value     string `json:"value,omitempty"`
}

type ExtractionOptions struct {
	UseAI           bool    `json:"use_ai"`
	FallbackToCSS   bool    `json:"fallback_to_css"`
	ConfidenceMin   float64 `json:"confidence_min"`
	MaxRetries      int     `json:"max_retries"`
	Timeout         int     `json:"timeout"`
}

type ExtractionResult struct {
	Data       map[string]interface{} `json:"data"`
	Confidence float64                `json:"confidence"`
	Method     string                 `json:"method"`
	Errors     []string               `json:"errors,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type TrainingData struct {
	Examples []TrainingExample `json:"examples"`
	Schema   *ExtractionSchema `json:"schema"`
}

type TrainingExample struct {
	HTML     string                 `json:"html"`
	Expected map[string]interface{} `json:"expected"`
	URL      string                 `json:"url,omitempty"`
}

func NewAIExtractor(config *AIConfig) *AIExtractor {
	extractor := &AIExtractor{
		models: make(map[string]Model),
		config: config,
	}

	for name, modelConfig := range config.Models {
		model := extractor.createModel(modelConfig)
		if model != nil {
			extractor.models[name] = model
		}
	}

	return extractor
}

func (a *AIExtractor) Extract(ctx context.Context, input *ExtractionInput) (*ExtractionResult, error) {
	cssResult := a.extractWithCSS(input)
	
	if input.Options != nil && input.Options.UseAI {
		aiResult, err := a.extractWithAI(ctx, input)
		if err == nil && aiResult.Confidence >= input.Options.ConfidenceMin {
			return aiResult, nil
		}
	}

	if input.Options != nil && input.Options.FallbackToCSS {
		return cssResult, nil
	}

	return nil, fmt.Errorf("extraction failed")
}

func (a *AIExtractor) extractWithCSS(input *ExtractionInput) *ExtractionResult {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(input.HTML))
	if err != nil {
		return &ExtractionResult{
			Errors: []string{fmt.Sprintf("HTML parse error: %v", err)},
		}
	}

	data := make(map[string]interface{})
	errors := []string{}

	for _, field := range input.Schema.Fields {
		if field.Selector == "" {
			continue
		}

		selection := doc.Find(field.Selector)
		if selection.Length() == 0 && field.Required {
			errors = append(errors, fmt.Sprintf("Required field '%s' not found", field.Name))
			continue
		}

		var value interface{}
		if field.Multiple {
			var values []string
			selection.Each(func(i int, s *goquery.Selection) {
				val := a.extractValue(s, field)
				if val != "" {
					values = append(values, val)
				}
			})
			value = values
		} else {
			value = a.extractValue(selection.First(), field)
		}

		if value != nil {
			data[field.Name] = value
		}
	}

	return &ExtractionResult{
		Data:       data,
		Confidence: 0.8, 
		Method:     "css",
		Errors:     errors,
	}
}

func (a *AIExtractor) extractValue(selection *goquery.Selection, field FieldSchema) string {
	if field.Attribute != "" {
		val, exists := selection.Attr(field.Attribute)
		if exists {
			return strings.TrimSpace(val)
		}
		return ""
	}
	return strings.TrimSpace(selection.Text())
}

func (a *AIExtractor) extractWithAI(ctx context.Context, input *ExtractionInput) (*ExtractionResult, error) {
	modelName := a.config.DefaultModel
	model, exists := a.models[modelName]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", modelName)
	}

	return model.Extract(ctx, input)
}

func (a *AIExtractor) createModel(config ModelConfig) Model {
	switch config.Type {
	case "openai":
		return &MockModel{modelType: "openai"}
	case "huggingface":
		return &MockModel{modelType: "huggingface"}
	case "local":
		return &MockModel{modelType: "local"}
	default:
		return nil
	}
}

type MockModel struct {
	modelType string
}

func (m *MockModel) Extract(ctx context.Context, input *ExtractionInput) (*ExtractionResult, error) {
	return &ExtractionResult{
		Data: map[string]interface{}{
			"title": "Mock Title",
			"price": 99.99,
		},
		Confidence: 0.9,
		Method:     m.modelType,
	}, nil
}

func (m *MockModel) Train(ctx context.Context, data *TrainingData) error {
	return nil
}

func (m *MockModel) Predict(ctx context.Context, features []float64) ([]float64, error) {
	return []float64{0.9}, nil
}

type SmartExtractor struct {
	aiExtractor *AIExtractor
	patterns    map[string]*ExtractionPattern
	cache       map[string]*ExtractionResult
}

type ExtractionPattern struct {
	Name        string            `json:"name"`
	URLPattern  string            `json:"url_pattern"`
	Schema      *ExtractionSchema `json:"schema"`
	Confidence  float64           `json:"confidence"`
	LastUpdated string            `json:"last_updated"`
}

func NewSmartExtractor(aiExtractor *AIExtractor) *SmartExtractor {
	return &SmartExtractor{
		aiExtractor: aiExtractor,
		patterns:    make(map[string]*ExtractionPattern),
		cache:       make(map[string]*ExtractionResult),
	}
}

func (s *SmartExtractor) LearnPattern(url string, result *ExtractionResult) {
	domain := extractDomain(url)
	
	pattern, exists := s.patterns[domain]
	if !exists {
		pattern = &ExtractionPattern{
			Name:       domain,
			URLPattern: fmt.Sprintf("*%s*", domain),
			Schema:     s.generateSchema(result.Data),
			Confidence: result.Confidence,
		}
		s.patterns[domain] = pattern
	} else {
		pattern.Confidence = (pattern.Confidence + result.Confidence) / 2
	}
}

func (s *SmartExtractor) generateSchema(data map[string]interface{}) *ExtractionSchema {
	var fields []FieldSchema
	
	for key, value := range data {
		field := FieldSchema{
			Name:     key,
			Type:     inferType(value),
			Required: true,
		}
		fields = append(fields, field)
	}
	
	return &ExtractionSchema{
		Fields: fields,
	}
}

func inferType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int64, float64:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "string"
	}
}

func extractDomain(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return url
}

type JSONExtractor struct {
	schema *JSONSchema
}

type JSONSchema struct {
	Paths map[string]string `json:"paths"`
}

func NewJSONExtractor(schema *JSONSchema) *JSONExtractor {
	return &JSONExtractor{schema: schema}
}

func (j *JSONExtractor) Extract(jsonData string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	for field, path := range j.schema.Paths {
		value := gjson.Get(jsonData, path)
		if value.Exists() {
			result[field] = value.Value()
		}
	}
	
	return result, nil
}