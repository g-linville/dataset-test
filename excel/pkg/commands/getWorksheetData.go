package commands

import (
	"context"

	"github.com/g-linville/dataset-test/excel/pkg/client"
	"github.com/g-linville/dataset-test/excel/pkg/global"
	"github.com/g-linville/dataset-test/excel/pkg/graph"
	"github.com/g-linville/dataset-test/excel/pkg/printers"
)

func GetWorksheetData(ctx context.Context, workbookID, worksheetID string) error {
	c, err := client.NewClient(global.ReadOnlyScopes)
	if err != nil {
		return err
	}

	data, _, err := graph.GetWorksheetData(ctx, c, workbookID, worksheetID)

	return printers.PrintWorksheetData(data)
}
