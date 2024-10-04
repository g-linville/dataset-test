package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/g-linville/dataset-test/excel/pkg/util"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type WorkbookInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListWorkbooks(ctx context.Context, c *msgraphsdkgo.GraphServiceClient) ([]WorkbookInfo, error) {
	drive, err := c.Me().Drive().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var infos []WorkbookInfo
	workbooks, err := c.Drives().ByDriveId(util.Deref(drive.GetId())).SearchWithQ(util.Ptr("xlsx")).GetAsSearchWithQGetResponse(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, workbook := range workbooks.GetValue() {
		infos = append(infos, WorkbookInfo{
			ID:   util.Deref(workbook.GetId()),
			Name: util.Deref(workbook.GetName()),
		})
	}
	return infos, nil
}

type WorksheetInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	WorkbookID string `json:"workbook_id"`
}

func ListWorksheetsInWorkbook(ctx context.Context, c *msgraphsdkgo.GraphServiceClient, workbookID string) ([]WorksheetInfo, error) {
	drive, err := c.Me().Drive().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	sheets, err := c.Drives().ByDriveId(util.Deref(drive.GetId())).Items().ByDriveItemId(workbookID).Workbook().Worksheets().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var infos []WorksheetInfo
	for _, sheet := range sheets.GetValue() {
		infos = append(infos, WorksheetInfo{
			ID:         util.Deref(sheet.GetId()),
			Name:       util.Deref(sheet.GetName()),
			WorkbookID: workbookID,
		})
	}
	return infos, nil
}

func GetWorksheetData(ctx context.Context, c *msgraphsdkgo.GraphServiceClient, workbookID, worksheetID string) ([][]any, models.WorkbookRangeable, error) {
	drive, err := c.Me().Drive().Get(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	usedRange, err := c.Drives().ByDriveId(util.Deref(drive.GetId())).Items().ByDriveItemId(workbookID).Workbook().Worksheets().ByWorkbookWorksheetId(worksheetID).UsedRange().Get(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	result, err := serialization.SerializeToJson(usedRange.GetValues())
	if err != nil {
		return nil, nil, err
	}

	var data [][]any
	if err = json.Unmarshal(result, &data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return data, usedRange, nil
}
