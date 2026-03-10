package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"ai_devs_4/s01e01/app"
)

func main() {
	przeslij()
	//content, err := app.Fetch()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//items, err := app.ParseCSV(content)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//items = app.FilterItems(items)
	//resp := app.FinalJsonResponse{
	//	Apikey: app.API_KEY,
	//	Task:   "people",
	//	Answer: nil,
	//}
	//answer := make([]app.Answer, 0)
	//jobs := make([]string, 0)
	//for _, item := range items {
	//	jobs = append(jobs, item.Job)
	//
	//	//tags, err := app.GetTags(item.Job)
	//	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	//	//
	//	//answer = append(answer, app.Answer{
	//	//	Name:    item.Name,
	//	//	Surname: item.Surname,
	//	//	Gender:  item.Gender,
	//	//	Born:    app.CalculateAge(item.BirthDate),
	//	//	City:    item.BirthPlace,
	//	//	Tags:    tags,
	//	//})
	//	//break
	//}
	//resp.Answer = answer
	//fmt.Println(resp)
	//
	//tagss, err := app.BatchTagJobs(jobs)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if len(jobs) != len(tagss) {
	//	log.Fatal("uneven tags for jobs")
	//}
	//for i, tags := range tagss {
	//	item := items[i]
	//	answer = append(answer, app.Answer{
	//		Name:    item.Name,
	//		Surname: item.Surname,
	//		Gender:  item.Gender,
	//		Born:    app.CalculateAge(item.BirthDate),
	//		City:    item.BirthPlace,
	//		Tags:    tags,
	//	})
	//}
	//resp.Answer = answer
	//data, err := json.MarshalIndent(resp, "", "  ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = os.WriteFile("./output.json", data, fs.ModePerm)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func przeslij() {
	content, err := os.ReadFile("./output.json")
	if err != nil {
		log.Fatal(err)
	}
	resp := app.FinalJsonResponse{}
	json.Unmarshal(content, &resp)
	items := app.GetFilteredItems()
	fmt.Println(len(resp.Answer), len(items))
	for i, answer := range resp.Answer {
		for _, item := range items {
			if answer.Name == item.Name && answer.Surname == item.Surname {
				resp.Answer[i].Born = item.BirthDate.Year()
			}
		}
	}
	content, _ = json.Marshal(resp)

	r, err := http.Post("https://hub.ag3nts.org/verify", "application/json", bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	fmt.Println(r.StatusCode)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
