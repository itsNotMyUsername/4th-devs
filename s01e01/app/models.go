package app

import "time"

type Item struct {
	Name         string
	Surname      string
	Gender       string
	BirthDate    time.Time
	BirthPlace   string
	BirthCountry string
	Job          string
}

type PostRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Text  struct {
		Format ResponseFormat `json:"format"`
	} `json:"text"`
}

type Response struct {
	Id                string        `json:"id"`
	Object            string        `json:"object"`
	CreatedAt         int           `json:"created_at"`
	Model             string        `json:"model"`
	Status            string        `json:"status"`
	CompletedAt       int           `json:"completed_at"`
	Output            []Output      `json:"output"`
	Error             interface{}   `json:"error"`
	IncompleteDetails interface{}   `json:"incomplete_details"`
	Tools             []interface{} `json:"tools"`
	ToolChoice        string        `json:"tool_choice"`
	ParallelToolCalls bool          `json:"parallel_tool_calls"`
	MaxOutputTokens   interface{}   `json:"max_output_tokens"`
	Temperature       int           `json:"temperature"`
	TopP              int           `json:"top_p"`
	PresencePenalty   int           `json:"presence_penalty"`
	FrequencyPenalty  int           `json:"frequency_penalty"`
	TopLogprobs       int           `json:"top_logprobs"`
	MaxToolCalls      interface{}   `json:"max_tool_calls"`
	Metadata          struct {
	} `json:"metadata"`
	Background         bool        `json:"background"`
	PreviousResponseId interface{} `json:"previous_response_id"`
	ServiceTier        string      `json:"service_tier"`
	Truncation         string      `json:"truncation"`
	Store              bool        `json:"store"`
	Instructions       interface{} `json:"instructions"`
	Text               struct {
		Format struct {
			Type string `json:"type"`
		} `json:"format"`
	} `json:"text"`
	Reasoning        interface{} `json:"reasoning"`
	SafetyIdentifier interface{} `json:"safety_identifier"`
	PromptCacheKey   interface{} `json:"prompt_cache_key"`
	User             interface{} `json:"user"`
	Usage            struct {
		InputTokens        int `json:"input_tokens"`
		InputTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"input_tokens_details"`
		OutputTokens        int `json:"output_tokens"`
		OutputTokensDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"output_tokens_details"`
		TotalTokens int  `json:"total_tokens"`
		Cost        int  `json:"cost"`
		IsByok      bool `json:"is_byok"`
		CostDetails struct {
			UpstreamInferenceCost       int `json:"upstream_inference_cost"`
			UpstreamInferenceInputCost  int `json:"upstream_inference_input_cost"`
			UpstreamInferenceOutputCost int `json:"upstream_inference_output_cost"`
		} `json:"cost_details"`
	} `json:"usage"`
}

type Output struct {
	Id      string    `json:"id"`
	Type    string    `json:"type"`
	Role    string    `json:"role"`
	Status  string    `json:"status"`
	Content []Content `json:"content"`
}

type Content struct {
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Annotations []interface{} `json:"annotations"`
	Logprobs    []interface{} `json:"logprobs"`
}

type ResponseFormat struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
	Strict bool                   `json:"strict"`
}

var tagsSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"tags": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "string",
				"enum": []string{
					"IT",
					"transport",
					"edukacja",
					"medycyna",
					"praca z ludźmi",
					"praca z pojazdami",
					"praca fizyczna",
				},
			},
		},
	},
	"required":             []string{"tags"},
	"additionalProperties": false,
}

var batchTagsSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"results": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"tags": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
							"enum": []string{
								"IT",
								"transport",
								"edukacja",
								"medycyna",
								"praca z ludźmi",
								"praca z pojazdami",
								"praca fizyczna",
							},
						},
					},
				},
				"required": []string{"tags"},
			},
		},
	},
	"required": []string{"results"},
}

type FinalJsonResponse struct {
	Apikey string   `json:"apikey"`
	Task   string   `json:"task"`
	Answer []Answer `json:"answer"`
}

type Answer struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Gender  string   `json:"gender"`
	Born    int      `json:"born"`
	City    string   `json:"city"`
	Tags    []string `json:"tags"`
}
