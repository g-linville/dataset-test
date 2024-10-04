package printers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/g-linville/dataset-test/excel/pkg/graph"
)

func PrintWorkbookInfos(infos []graph.WorkbookInfo) {
	for _, info := range infos {
		fmt.Printf("Name: %s\n", info.Name)
		fmt.Printf("  ID: %s\n", info.ID)
	}
}

func PrintWorksheetInfos(infos []graph.WorksheetInfo) {
	for _, info := range infos {
		fmt.Printf("Name: %s\n", info.Name)
		fmt.Printf("  ID: %s\n", info.ID)
		fmt.Printf("  Workbook ID: %s\n", info.WorkbookID)
	}
}

type dataset struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type data struct {
	Data []any `json:"data"`
}

func PrintWorksheetData(d [][]any) error {
	data := data{
		Data: arrayToAnyArray(d),
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}

	dataHash := fmt.Sprintf("%x", sha256.Sum256(dataBytes))
	dataset := dataset{
		ID:   dataHash + ".json",
		Type: "gptscript_dataset",
	}

	if err := os.WriteFile(os.Getenv("GPTSCRIPT_WORKSPACE_DIR")+string(filepath.Separator)+dataset.ID, dataBytes, 0644); err != nil {
		return fmt.Errorf("error writing data file: %v", err)
	}

	datasetJSON, _ := json.Marshal(dataset)
	if len(datasetJSON) == 0 {
		fmt.Println("no data")
	} else {
		fmt.Println(string(datasetJSON))
	}
	return nil
}

func arrayToAnyArray[T any](arr []T) []any {
	var result []any
	for _, item := range arr {
		result = append(result, item)
	}
	return result
}
