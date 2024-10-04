package main

import (
	"context"
	"fmt"
	"os"

	"github.com/g-linville/dataset-test/excel/pkg/commands"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gptscript-go-tool <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	var err error
	switch command {
	case "listWorkbooks":
		err = commands.ListWorkbooks(context.Background())
	case "listWorksheets":
		err = commands.ListWorksheets(context.Background(), os.Getenv("WORKBOOK_ID"))
	case "getWorksheetData":
		err = commands.GetWorksheetData(context.Background(), os.Getenv("WORKBOOK_ID"), os.Getenv("WORKSHEET_ID"))
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
