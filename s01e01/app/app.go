package app

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const model = "openai/gpt-5-nano"

var API_KEY = os.Getenv("API_KEY")

var csvURL = fmt.Sprintf("https://hub.ag3nts.org/data/%s/people.csv", API_KEY)

func Fetch() ([]byte, error) {
	resp, err := http.Get(csvURL)
	if err != nil {
		slog.Error("failed to get url response")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to get url response")
		return nil, err
	}
	return body, nil
}

func ParseCSV(input []byte) ([]Item, error) {
	items := make([]Item, 0)

	r := csv.NewReader(bytes.NewReader(input))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] == "name" {
			continue
		}

		birthDate, err := time.Parse(time.DateOnly, record[3])
		if err != nil {
			log.Fatal(err)
		}

		items = append(items, Item{
			Name:         record[0],
			Surname:      record[1],
			Gender:       record[2],
			BirthDate:    birthDate,
			BirthPlace:   record[4],
			BirthCountry: record[5],
			Job:          record[6],
		})

	}
	return items, nil
}

func FilterItems(items []Item) []Item {
	output := make([]Item, 0)

	for _, item := range items {
		age := CalculateAge(item.BirthDate)

		if item.Gender != "M" {
			continue
		}

		if item.BirthPlace != "Grudziądz" {
			continue
		}

		if age <= 20 || age >= 40 {
			continue
		}

		output = append(output, item)
	}

	return output
}

func GetFilteredItems() []Item {
	csvData, _ := Fetch()
	items, _ := ParseCSV(csvData)
	return FilterItems(items)
}

func CalculateAge(input time.Time) int {
	now := time.Now()
	age := now.Year() - input.Year()
	if now.YearDay() < input.YearDay() {
		age--
	}
	return age
}

//func GetTags(input string) ([]string, error) {
//	prompt := fmt.Sprintf(
//		`Sklasyfikuj ten opis stanowiska do jednej lub więcej kategorii:
//"IT","transport","edukacja","medycyna","praca z ludźmi","praca z pojazdami","praca fizyczna".
//Opis: %s`, input)
//
//	response, err := requestToLLM(prompt, tagsSchema)
//	if err != nil {
//		return nil, err
//	}
//	tags := make([]string, 0)
//	for _, output := range response.Output {
//		for _, content := range output.Content {
//			if content.Type != "output_text" {
//				continue
//			}
//
//			var result struct {
//				Tags []string `json:"tags"`
//			}
//
//			err = json.Unmarshal([]byte(content.Text), &result)
//			if err != nil {
//				return nil, err
//			}
//
//			tags = append(tags, result.Tags...)
//		}
//	}
//	return tags, nil
//}

func BatchTagJobs(jobs []string) ([][]string, error) {
	prompt := fmt.Sprintf(`
Masz listę stanowisk pracy: %v
Otaguj każde stanowisko tylko tagami z listy:
["IT", "transport", "edukacja", "medycyna", "praca z ludźmi", "praca z pojazdami", "praca fizyczna"]
Odpowiedz wyłącznie tablicą obiektów {tags:[...]} w dokładnie tej samej kolejności co input.
`, jobs)

	// JSON Schema dla pojedynczego elementu
	itemSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"tags": map[string]any{
				"type": "array",
				"items": map[string]any{
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

	// JSON Schema dla całej tablicy
	jsonSchema := map[string]any{
		"name":   "job_tags",
		"strict": true,
		"schema": map[string]any{
			"type":  "array",
			"items": itemSchema,
		},
	}

	response, err := requestToLLM(prompt, jsonSchema)
	if err != nil {
		slog.Error("failed to requestToLLM", "err", err)
		return nil, err
	}

	tagss := make([][]string, 0)
	for _, output := range response.Output {
		for _, content := range output.Content {
			if content.Type != "output_text" {
				slog.Error("skipped", "content", content)
				continue
			}

			var result struct {
				Results []struct {
					Tags []string `json:"tags"`
				} `json:"results"`
			}

			err = json.Unmarshal([]byte(content.Text), &result)
			if err != nil {
				slog.Error("failed to unmarshal response result", "err", err)
				return nil, err
			}

			for _, res := range result.Results {
				tagss = append(tagss, res.Tags)
			}
		}
	}

	return tagss, nil
}

func requestToLLM(prompt string, schema map[string]interface{}) (*Response, error) {
	requestBody := map[string]any{
		"model": model,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]string{
					{"type": "text", "text": prompt},
				},
			},
		},
		"response_format": map[string]any{
			"type":        "json_schema",
			"json_schema": schema,
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		slog.Error("failed to marshal postRequest", "err", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://openrouter.ai/api/v1/responses", bytes.NewReader(body))
	if err != nil {
		slog.Error("failed to make NewRequest", "err", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENROUTER_API_KEY")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to Do", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to ReadAll", "err", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(time.Second)
		return requestToLLM(prompt, schema)
	}

	response := Response{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		slog.Error("failed to unmarshal response", "err", err)
		return nil, err
	}
	return &response, nil
}
