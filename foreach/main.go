package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gptscript-ai/go-gptscript"
)

// These values are hardcoded for the proof of concept.
// In a real-world scenario, they would be configured by the user in the UI.
const (
	instructions = "Find the star count of the repo. Return just the star number."
	criterion    = "The org name of this GitHub repo is gptscript-ai."
	tools        = "github.com/gptscript-ai/tools/apis/github/read,github.com/gptscript-ai/tools/apis/excel/write"
)

type data struct {
	Data []any `json:"data"`
}

func main() {
	datasetID := os.Getenv("DATASET")
	if datasetID == "" {
		fmt.Println("Error: no dataset ID provided")
		os.Exit(1)
	}

	datasetType := os.Getenv("DATASET_TYPE")
	if datasetType == "" {
		fmt.Println("Error: no dataset type provided")
		os.Exit(1)
	}

	if datasetType != "gptscript_dataset" {
		fmt.Println("Error: not a GPTScript dataset")
		os.Exit(1)
	} else if !strings.HasSuffix(datasetID, ".json") {
		fmt.Println("Error: dataset ID must end in .json")
		os.Exit(1)
	}

	fileName := os.Getenv("GPTSCRIPT_WORKSPACE_DIR") + string(filepath.Separator) + datasetID
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening dataset file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = file.Close()
	}()

	var data data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Printf("Error decoding dataset file: %v\n", err)
		os.Exit(1)
	}

	client, err := gptscript.NewGPTScript()
	if err != nil {
		fmt.Printf("Error creating GPTScript client: %v\n", err)
		os.Exit(1)
	}

	output := map[string]string{}
	for _, item := range data.Data {
		itemStr, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("Error marshalling item: %v\n", err)
			os.Exit(1)
		}

		// First, evaluate the criterion to see if we should process this item.
		criterionDef := gptscript.ToolDef{
			Instructions: fmt.Sprintf(`Determine whether the following data meets the criterion. Return exactly TRUE if it does, and exactly FALSE if it does not.

Data: %s
Criterion: %s`, string(itemStr), criterion),
		}

		criterionRun, err := client.Evaluate(context.Background(), gptscript.Options{}, criterionDef)
		if err != nil {
			fmt.Printf("Error evaluating criterion: %v\n", err)
			os.Exit(1)
		}

		criterionOut, err := criterionRun.Text()
		if err != nil {
			fmt.Printf("Error getting criterion output text: %v\n", err)
			os.Exit(1)
		}

		if strings.Contains(criterionOut, "FALSE") || !strings.Contains(criterionOut, "TRUE") {
			continue
		}

		def := gptscript.ToolDef{
			Tools:        strings.Split(tools, ","),
			Instructions: instructions,
		}

		run, err := client.Evaluate(context.Background(), gptscript.Options{
			Input: string(itemStr),
		}, def)
		if err != nil {
			fmt.Printf("Error evaluating GPTScript: %v\n", err)
			os.Exit(1)
		}

		out, err := run.Text()
		if err != nil {
			fmt.Printf("Error running GPTScript: %v\n", err)
			os.Exit(1)
		}

		output[string(itemStr)] = out
	}

	outputBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Printf("Error marshalling output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(outputBytes))
}
