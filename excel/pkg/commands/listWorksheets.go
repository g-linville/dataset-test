package commands

import (
	"context"

	"github.com/g-linville/dataset-test/excel/pkg/client"
	"github.com/g-linville/dataset-test/excel/pkg/global"
	"github.com/g-linville/dataset-test/excel/pkg/graph"
	"github.com/g-linville/dataset-test/excel/pkg/printers"
)

func ListWorksheets(ctx context.Context, workbookID string) error {
	c, err := client.NewClient(global.ReadOnlyScopes)
	if err != nil {
		return err
	}

	infos, err := graph.ListWorksheetsInWorkbook(ctx, c, workbookID)
	if err != nil {
		return err
	}

	printers.PrintWorksheetInfos(infos)
	return nil
}
