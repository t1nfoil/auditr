package main

import (
	"context"
	"log"
	"time"

	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/auditlogs"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func DirectoryAudit(client *graph.GraphServiceClient, ctx context.Context, filter string) {

	auditRequestQuery := auditlogs.DirectoryAuditsRequestBuilderGetQueryParameters{
		Filter: &filter,
	}
	auditRequestBuilder := auditlogs.DirectoryAuditsRequestBuilderGetRequestConfiguration{
		QueryParameters: &auditRequestQuery,
	}

	var auditResult models.DirectoryAuditCollectionResponseable
	auditResult, err := client.AuditLogs().DirectoryAudits().Get(ctx, &auditRequestBuilder)
	if err != nil {
		printOdataError(err)
		log.Fatal(err)
	}

	pageIterator, err := graphcore.NewPageIterator[*models.DirectoryAudit](auditResult, client.GetAdapter(), models.CreateDirectoryAuditCollectionResponseFromDiscriminatorValue)
	if err != nil {
		printOdataError(err)
		log.Fatal(err)
	}

	err = pageIterator.Iterate(context.Background(), ProcessDirectoryAudit)
	if err != nil {
		printOdataError(err)
	}
}

func ProcessDirectoryAudit(event *models.DirectoryAudit) bool {
	if event != nil {
		AddDirectoryAudit(event)
		return true
	} else {
		return false
	}
}

type DirectoryAuditWrite struct {
	DirectoryAudit      *models.DirectoryAudit
	DirectoryNL         string
	hasModelBeenWritten bool
}

var directoryAuditsMap = make(map[string]DirectoryAuditWrite)

func AddDirectoryAudit(event *models.DirectoryAudit) {
	daLock.Lock()
	defer daLock.Unlock()
	if event != nil {
		if _, exists := directoryAuditsMap[*event.GetId()]; !exists {
			var d = DirectoryAuditWrite{
				DirectoryAudit:      event,
				hasModelBeenWritten: false,
			}
			directoryAuditsMap[*event.GetId()] = d
		}
	}
}

func GetMostRecentDirectoryAuditTimestamp() string {
	daLock.Lock()
	defer daLock.Unlock()
	var mostRecentTimeStamp string
	for _, v := range directoryAuditsMap {
		timestamp := v.DirectoryAudit.GetActivityDateTime().Format("2006-01-02T15:04:05Z")
		if timestamp > mostRecentTimeStamp {
			mostRecentTimeStamp = timestamp
		}
	}
	return mostRecentTimeStamp
}

type trimRecord struct {
	startTime     time.Time
	thresholdTime time.Time
	count         int
}

func trimDirectoryAudits() trimRecord {
	daLock.Lock()
	defer daLock.Unlock()

	currentTime := time.Now().UTC()
	var trimResult = trimRecord{
		startTime:     currentTime,
		thresholdTime: currentTime.Add(-time.Duration(config.TrimEntryAge) * time.Second),
		count:         0,
	}

	beforeTrim := len(directoryAuditsMap)
	for k, v := range directoryAuditsMap {
		if v.DirectoryAudit.GetActivityDateTime().Before(trimResult.thresholdTime) && v.hasModelBeenWritten {
			delete(directoryAuditsMap, k)
		}
	}
	trimResult.count = beforeTrim - len(directoryAuditsMap)
	return trimResult
}
