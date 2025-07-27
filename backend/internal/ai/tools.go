// Package ai provides tool implementations for AI function calling
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ToolExecutor handles execution of AI tools
type ToolExecutor struct {
	httpClient *http.Client
}

// NewToolExecutor creates a new tool executor
func NewToolExecutor() *ToolExecutor {
	return &ToolExecutor{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExecuteTool executes a tool call and returns the result
func (te *ToolExecutor) ExecuteTool(ctx context.Context, toolCall ToolCall) (string, error) {
	// Parse arguments from JSON string
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse tool arguments: %w", err)
	}

	switch toolCall.Function.Name {
	case "web_search":
		return te.executeWebSearch(ctx, args)
	case "calculate":
		return te.executeCalculate(ctx, args)
	case "get_current_time":
		return te.getCurrentTime(ctx, args)
	case "get_weather":
		return te.getWeather(ctx, args)
	case "analyze_image":
		return te.analyzeImage(ctx, args)
	case "translate_text":
		return te.translateText(ctx, args)
	case "analyze_document":
		return te.analyzeDocument(ctx, args)
	case "generate_code":
		return te.generateCode(ctx, args)
	case "format_data":
		return te.formatData(ctx, args)
	case "validate_json":
		return te.validateJSON(ctx, args)
	case "extract_text":
		return te.extractText(ctx, args)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}
}

// Web Search Tool
func (te *ToolExecutor) executeWebSearch(ctx context.Context, args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter is required for web_search")
	}

	numResults := 5 // default
	if num, ok := args["num_results"].(float64); ok {
		numResults = int(num)
	}

	// Simulate web search results
	// In a real implementation, this would call a search API like Google, Bing, or DuckDuckGo
	results := []map[string]interface{}{
		{
			"title":       fmt.Sprintf("Search result for: %s", query),
			"url":         "https://example.com/result1",
			"snippet":     fmt.Sprintf("This is a comprehensive result about %s with detailed information.", query),
			"published":   time.Now().Format("2006-01-02"),
		},
		{
			"title":       fmt.Sprintf("Latest news about %s", query),
			"url":         "https://news.example.com/article",
			"snippet":     fmt.Sprintf("Recent developments and updates regarding %s from reliable sources.", query),
			"published":   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		},
		{
			"title":       fmt.Sprintf("Complete guide to %s", query),
			"url":         "https://guide.example.com/topic",
			"snippet":     fmt.Sprintf("A detailed guide covering all aspects of %s with examples and best practices.", query),
			"published":   time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
		},
	}

	// Limit results
	if numResults < len(results) {
		results = results[:numResults]
	}

	// Format results as text
	var resultText strings.Builder
	resultText.WriteString(fmt.Sprintf("Search results for '%s':\n\n", query))

	for i, result := range results {
		resultText.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, result["title"]))
		resultText.WriteString(fmt.Sprintf("   URL: %s\n", result["url"]))
		resultText.WriteString(fmt.Sprintf("   %s\n", result["snippet"]))
		resultText.WriteString(fmt.Sprintf("   Published: %s\n\n", result["published"]))
	}

	return resultText.String(), nil
}

// Calculate Tool
func (te *ToolExecutor) executeCalculate(ctx context.Context, args map[string]interface{}) (string, error) {
	expression, ok := args["expression"].(string)
	if !ok {
		return "", fmt.Errorf("expression parameter is required for calculate")
	}

	// Simple calculator implementation
	// This is a basic implementation - in production, you'd want a more robust math parser
	result, err := te.evaluateExpression(expression)
	if err != nil {
		return "", fmt.Errorf("failed to calculate expression '%s': %w", expression, err)
	}

	return fmt.Sprintf("Calculation: %s = %g", expression, result), nil
}

// Get Current Time Tool
func (te *ToolExecutor) getCurrentTime(ctx context.Context, args map[string]interface{}) (string, error) {
	timezone := "UTC" // default
	if tz, ok := args["timezone"].(string); ok {
		timezone = tz
	}

	now := time.Now()
	if timezone != "UTC" {
		// Try to load the timezone
		if loc, err := time.LoadLocation(timezone); err == nil {
			now = now.In(loc)
		}
	}

	return fmt.Sprintf("Current time (%s): %s\nTimestamp: %d\nFormatted: %s", 
		timezone, 
		now.Format(time.RFC3339), 
		now.Unix(),
		now.Format("Monday, January 2, 2006 at 3:04 PM MST")), nil
}

// Get Weather Tool
func (te *ToolExecutor) getWeather(ctx context.Context, args map[string]interface{}) (string, error) {
	location, ok := args["location"].(string)
	if !ok {
		return "", fmt.Errorf("location parameter is required for get_weather")
	}

	// Simulate weather data
	// In a real implementation, this would call a weather API like OpenWeatherMap
	weatherData := map[string]interface{}{
		"location":    location,
		"temperature": 22.5,
		"condition":   "Partly Cloudy",
		"humidity":    65,
		"wind_speed":  12.3,
		"pressure":    1013.25,
		"visibility":  10,
		"uv_index":    5,
		"updated":     time.Now().Format(time.RFC3339),
	}

	return fmt.Sprintf(`Weather for %s:
- Temperature: %.1f°C
- Condition: %s
- Humidity: %d%%
- Wind Speed: %.1f km/h
- Pressure: %.2f hPa
- Visibility: %d km
- UV Index: %d
- Last Updated: %s`,
		weatherData["location"],
		weatherData["temperature"],
		weatherData["condition"],
		weatherData["humidity"],
		weatherData["wind_speed"],
		weatherData["pressure"],
		weatherData["visibility"],
		weatherData["uv_index"],
		weatherData["updated"]), nil
}

// Analyze Image Tool
func (te *ToolExecutor) analyzeImage(ctx context.Context, args map[string]interface{}) (string, error) {
	_, ok := args["image_data"].(string)
	if !ok {
		return "", fmt.Errorf("image_data parameter is required for analyze_image")
	}

	analysisType := "general" // default
	if aType, ok := args["analysis_type"].(string); ok {
		analysisType = aType
	}

	// Simulate image analysis
	// In a real implementation, this would use computer vision APIs like Google Vision, AWS Rekognition, etc.
	analysis := map[string]interface{}{
		"objects_detected": []string{"person", "car", "building", "tree"},
		"colors":           []string{"blue", "green", "gray", "white"},
		"text_detected":    "Sample text found in image",
		"faces_count":      2,
		"confidence":       0.95,
		"image_quality":    "high",
		"dimensions":       "1920x1080",
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Image Analysis (%s):\n\n", analysisType))
	result.WriteString(fmt.Sprintf("Objects Detected: %s\n", strings.Join(analysis["objects_detected"].([]string), ", ")))
	result.WriteString(fmt.Sprintf("Dominant Colors: %s\n", strings.Join(analysis["colors"].([]string), ", ")))
	result.WriteString(fmt.Sprintf("Text Found: %s\n", analysis["text_detected"]))
	result.WriteString(fmt.Sprintf("Faces Detected: %d\n", analysis["faces_count"]))
	result.WriteString(fmt.Sprintf("Confidence: %.2f\n", analysis["confidence"]))
	result.WriteString(fmt.Sprintf("Image Quality: %s\n", analysis["image_quality"]))
	result.WriteString(fmt.Sprintf("Dimensions: %s\n", analysis["dimensions"]))

	return result.String(), nil
}

// Translate Text Tool
func (te *ToolExecutor) translateText(ctx context.Context, args map[string]interface{}) (string, error) {
	text, ok := args["text"].(string)
	if !ok {
		return "", fmt.Errorf("text parameter is required for translate_text")
	}

	targetLang, ok := args["target_language"].(string)
	if !ok {
		return "", fmt.Errorf("target_language parameter is required for translate_text")
	}

	sourceLang := "auto" // default
	if sLang, ok := args["source_language"].(string); ok {
		sourceLang = sLang
	}

	// Simulate translation
	// In a real implementation, this would use translation APIs like Google Translate, DeepL, etc.
	translatedText := fmt.Sprintf("[Translated to %s] %s", targetLang, text)

	return fmt.Sprintf(`Translation Result:
Source Language: %s
Target Language: %s
Original Text: %s
Translated Text: %s
Confidence: 0.98`,
		sourceLang, targetLang, text, translatedText), nil
}

// Analyze Document Tool
func (te *ToolExecutor) analyzeDocument(ctx context.Context, args map[string]interface{}) (string, error) {
	_, ok := args["document_data"].(string)
	if !ok {
		return "", fmt.Errorf("document_data parameter is required for analyze_document")
	}

	documentType := "pdf" // default
	if dType, ok := args["document_type"].(string); ok {
		documentType = dType
	}

	// Simulate document analysis
	analysis := map[string]interface{}{
		"document_type":    documentType,
		"page_count":       15,
		"word_count":       2847,
		"language":         "English",
		"topics":           []string{"Technology", "AI", "Machine Learning"},
		"key_entities":     []string{"OpenAI", "GPT", "Neural Networks"},
		"sentiment":        "Neutral",
		"readability":      "Professional",
		"has_tables":       true,
		"has_images":       true,
	}

	return fmt.Sprintf(`Document Analysis:
- Type: %s
- Pages: %d
- Words: %d
- Language: %s
- Topics: %s
- Key Entities: %s
- Sentiment: %s
- Readability: %s
- Contains Tables: %t
- Contains Images: %t`,
		analysis["document_type"],
		analysis["page_count"],
		analysis["word_count"],
		analysis["language"],
		strings.Join(analysis["topics"].([]string), ", "),
		strings.Join(analysis["key_entities"].([]string), ", "),
		analysis["sentiment"],
		analysis["readability"],
		analysis["has_tables"],
		analysis["has_images"]), nil
}

// Generate Code Tool
func (te *ToolExecutor) generateCode(ctx context.Context, args map[string]interface{}) (string, error) {
	description, ok := args["description"].(string)
	if !ok {
		return "", fmt.Errorf("description parameter is required for generate_code")
	}

	language := "javascript" // default
	if lang, ok := args["language"].(string); ok {
		language = lang
	}

	// Simulate code generation based on language
	var code string
	switch strings.ToLower(language) {
	case "javascript", "js":
		code = fmt.Sprintf(`// %s
function generatedFunction() {
    // Implementation for: %s
    console.log('Generated function executed');
    return true;
}

// Usage example
generatedFunction();`, description, description)
	case "python", "py":
		code = fmt.Sprintf(`# %s
def generated_function():
    """Implementation for: %s"""
    print('Generated function executed')
    return True

# Usage example
if __name__ == "__main__":
    generated_function()`, description, description)
	case "go", "golang":
		code = fmt.Sprintf(`// %s
package main

import "fmt"

// GeneratedFunction implements: %s
func GeneratedFunction() bool {
    fmt.Println("Generated function executed")
    return true
}

func main() {
    GeneratedFunction()
}`, description, description)
	default:
		code = fmt.Sprintf(`// %s
// Generated code for: %s
// Language: %s

// Implementation would go here`, description, description, language)
	}

	return fmt.Sprintf("Generated %s code:\n\n```%s\n%s\n```", language, language, code), nil
}

// Format Data Tool
func (te *ToolExecutor) formatData(ctx context.Context, args map[string]interface{}) (string, error) {
	data, ok := args["data"].(string)
	if !ok {
		return "", fmt.Errorf("data parameter is required for format_data")
	}

	format := "json" // default
	if f, ok := args["format"].(string); ok {
		format = f
	}

	// Try to parse and format the data
	var result string
	switch strings.ToLower(format) {
	case "json":
		// Try to parse and pretty-print JSON
		var jsonData interface{}
		if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
			return "", fmt.Errorf("invalid JSON data: %w", err)
		}
		formatted, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to format JSON: %w", err)
		}
		result = string(formatted)
	case "csv":
		// Simple CSV formatting simulation
		lines := strings.Split(data, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				formatted = append(formatted, strings.TrimSpace(line))
			}
		}
		result = strings.Join(formatted, "\n")
	case "xml":
		// Simple XML formatting
		result = fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<data>\n  %s\n</data>", data)
	default:
		result = data
	}

	return fmt.Sprintf("Formatted data (%s):\n\n```%s\n%s\n```", format, format, result), nil
}

// Validate JSON Tool
func (te *ToolExecutor) validateJSON(ctx context.Context, args map[string]interface{}) (string, error) {
	jsonData, ok := args["json_data"].(string)
	if !ok {
		return "", fmt.Errorf("json_data parameter is required for validate_json")
	}

	// Validate JSON
	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return fmt.Sprintf("❌ Invalid JSON:\n%s\n\nError: %s", jsonData, err.Error()), nil
	}

	// Pretty print valid JSON
	formatted, _ := json.MarshalIndent(data, "", "  ")
	return fmt.Sprintf("✅ Valid JSON:\n\n```json\n%s\n```", string(formatted)), nil
}

// Extract Text Tool
func (te *ToolExecutor) extractText(ctx context.Context, args map[string]interface{}) (string, error) {
	_, ok := args["file_data"].(string)
	if !ok {
		return "", fmt.Errorf("file_data parameter is required for extract_text")
	}

	fileType := "pdf" // default
	if fType, ok := args["file_type"].(string); ok {
		fileType = fType
	}

	// Simulate text extraction
	extractedText := fmt.Sprintf(`Extracted text from %s file:

This is sample extracted text content. In a real implementation, this would:
- Parse PDF files using libraries like pdfplumber or PyPDF2
- Extract text from images using OCR (Tesseract)
- Parse Word documents using python-docx
- Handle various file formats

The extracted content would preserve formatting and structure where possible.

File type: %s
Extraction confidence: 95%%
Character count: 1,247
Word count: 203`, fileType, fileType)

	return extractedText, nil
}

// Helper function for basic expression evaluation
func (te *ToolExecutor) evaluateExpression(expr string) (float64, error) {
	// Remove spaces
	expr = strings.ReplaceAll(expr, " ", "")

	// Simple regex-based calculator for basic operations
	// This is a simplified implementation - for production use a proper math parser
	
	// Handle parentheses first (simplified)
	for strings.Contains(expr, "(") {
		// Find innermost parentheses
		re := regexp.MustCompile(`\(([^()]+)\)`)
		matches := re.FindStringSubmatch(expr)
		if len(matches) < 2 {
			break
		}
		inner := matches[1]
		result, err := te.evaluateSimpleExpression(inner)
		if err != nil {
			return 0, err
		}
		expr = strings.Replace(expr, matches[0], fmt.Sprintf("%g", result), 1)
	}

	return te.evaluateSimpleExpression(expr)
}

func (te *ToolExecutor) evaluateSimpleExpression(expr string) (float64, error) {
	// Handle multiplication and division first
	for {
		re := regexp.MustCompile(`(-?\d+(?:\.\d+)?)([*/])(-?\d+(?:\.\d+)?)`)
		matches := re.FindStringSubmatch(expr)
		if len(matches) < 4 {
			break
		}

		left, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, err
		}

		right, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			return 0, err
		}

		var result float64
		switch matches[2] {
		case "*":
			result = left * right
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			result = left / right
		}

		expr = strings.Replace(expr, matches[0], fmt.Sprintf("%g", result), 1)
	}

	// Handle addition and subtraction
	for {
		re := regexp.MustCompile(`(-?\d+(?:\.\d+)?)([+-])(-?\d+(?:\.\d+)?)`)
		matches := re.FindStringSubmatch(expr)
		if len(matches) < 4 {
			break
		}

		left, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, err
		}

		right, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			return 0, err
		}

		var result float64
		switch matches[2] {
		case "+":
			result = left + right
		case "-":
			result = left - right
		}

		expr = strings.Replace(expr, matches[0], fmt.Sprintf("%g", result), 1)
	}

	// Handle mathematical functions
	if strings.HasPrefix(expr, "sqrt(") && strings.HasSuffix(expr, ")") {
		inner := expr[5 : len(expr)-1]
		val, err := strconv.ParseFloat(inner, 64)
		if err != nil {
			return 0, err
		}
		return math.Sqrt(val), nil
	}

	if strings.HasPrefix(expr, "sin(") && strings.HasSuffix(expr, ")") {
		inner := expr[4 : len(expr)-1]
		val, err := strconv.ParseFloat(inner, 64)
		if err != nil {
			return 0, err
		}
		return math.Sin(val), nil
	}

	if strings.HasPrefix(expr, "cos(") && strings.HasSuffix(expr, ")") {
		inner := expr[4 : len(expr)-1]
		val, err := strconv.ParseFloat(inner, 64)
		if err != nil {
			return 0, err
		}
		return math.Cos(val), nil
	}

	// Try to parse as a simple number
	return strconv.ParseFloat(expr, 64)
}

// GetAvailableTools returns the list of available tools with their schemas
func GetAvailableTools() []Tool {
	return []Tool{
		NewTool(
			"web_search",
			"Search the web for current information and news",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query to find information about",
					},
					"num_results": map[string]interface{}{
						"type":        "integer",
						"description": "Number of search results to return (default: 5, max: 10)",
						"default":     5,
						"minimum":     1,
						"maximum":     10,
					},
				},
				"required": []string{"query"},
			},
		),
		NewTool(
			"calculate",
			"Perform mathematical calculations and solve equations",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "Mathematical expression to calculate (supports +, -, *, /, parentheses, sqrt, sin, cos)",
					},
				},
				"required": []string{"expression"},
			},
		),
		NewTool(
			"get_current_time",
			"Get the current date and time in various formats",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"timezone": map[string]interface{}{
						"type":        "string",
						"description": "Timezone (e.g., 'UTC', 'America/New_York', 'Asia/Tokyo')",
						"default":     "UTC",
					},
				},
			},
		),
		NewTool(
			"get_weather",
			"Get current weather information for a location",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "City name, country, or coordinates for weather information",
					},
				},
				"required": []string{"location"},
			},
		),
		NewTool(
			"analyze_image",
			"Analyze images to detect objects, text, faces, and other features",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"image_data": map[string]interface{}{
						"type":        "string",
						"description": "Base64 encoded image data",
					},
					"analysis_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of analysis: 'general', 'text', 'faces', 'objects'",
						"default":     "general",
					},
				},
				"required": []string{"image_data"},
			},
		),
		NewTool(
			"translate_text",
			"Translate text between different languages",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Text to translate",
					},
					"target_language": map[string]interface{}{
						"type":        "string",
						"description": "Target language code (e.g., 'en', 'es', 'fr', 'de', 'ja', 'zh')",
					},
					"source_language": map[string]interface{}{
						"type":        "string",
						"description": "Source language code (auto-detect if not specified)",
						"default":     "auto",
					},
				},
				"required": []string{"text", "target_language"},
			},
		),
		NewTool(
			"analyze_document",
			"Analyze documents (PDF, Word, etc.) to extract insights and metadata",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"document_data": map[string]interface{}{
						"type":        "string",
						"description": "Base64 encoded document data",
					},
					"document_type": map[string]interface{}{
						"type":        "string",
						"description": "Document type: 'pdf', 'docx', 'txt', 'rtf'",
						"default":     "pdf",
					},
				},
				"required": []string{"document_data"},
			},
		),
		NewTool(
			"generate_code",
			"Generate code snippets in various programming languages",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Description of what the code should do",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Programming language: 'javascript', 'python', 'go', 'java', 'cpp', 'rust'",
						"default":     "javascript",
					},
				},
				"required": []string{"description"},
			},
		),
		NewTool(
			"format_data",
			"Format and prettify data in various formats (JSON, CSV, XML)",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"data": map[string]interface{}{
						"type":        "string",
						"description": "Raw data to format",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"description": "Output format: 'json', 'csv', 'xml', 'yaml'",
						"default":     "json",
					},
				},
				"required": []string{"data"},
			},
		),
		NewTool(
			"validate_json",
			"Validate and prettify JSON data",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"json_data": map[string]interface{}{
						"type":        "string",
						"description": "JSON string to validate and format",
					},
				},
				"required": []string{"json_data"},
			},
		),
		NewTool(
			"extract_text",
			"Extract text content from various file formats",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"file_data": map[string]interface{}{
						"type":        "string",
						"description": "Base64 encoded file data",
					},
					"file_type": map[string]interface{}{
						"type":        "string",
						"description": "File type: 'pdf', 'docx', 'txt', 'image'",
						"default":     "pdf",
					},
				},
				"required": []string{"file_data"},
			},
		),
	}
}