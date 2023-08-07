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

func SignInActivityAudit(client *graph.GraphServiceClient, ctx context.Context, filter string) {

	signinRequestQuery := auditlogs.SignInsRequestBuilderGetQueryParameters{
		Filter: &filter,
	}

	signinRequestBuilder := auditlogs.SignInsRequestBuilderGetRequestConfiguration{
		QueryParameters: &signinRequestQuery,
	}

	var signinResult models.SignInCollectionResponseable

	signinResult, err := client.AuditLogs().SignIns().Get(ctx, &signinRequestBuilder)
	if err != nil {
		printOdataError(err)
		log.Fatal(err)
	}

	//spew.Dump(signinResult)
	pageIterator, err := graphcore.NewPageIterator[*models.SignIn](signinResult, client.GetAdapter(), models.CreateSignInActivityFromDiscriminatorValue)
	if err != nil {
		printOdataError(err)
		log.Fatal(err)
	}

	err = pageIterator.Iterate(context.Background(), ProcessSignInActivityAudit)
	if err != nil {
		printOdataError(err)
	}
}

func ProcessSignInActivityAudit(event *models.SignIn) bool {
	if event != nil {
		AddSignInActivity(event)
		return true
	} else {
		return false
	}
}

// function to add a event *models.SignInActivity to a global slice of *models.SignInActivity
type SignInWrite struct {
	SignIn              *models.SignIn
	hasModelBeenWritten bool
	SignInNL            string
}

var signInActivityMap = make(map[string]SignInWrite)
var addSkips = make(map[int32]int)

func AddSignInActivity(event *models.SignIn) {
	siLock.Lock()
	defer siLock.Unlock()

	if event != nil {
		if _, exists := signInActivityMap[*event.GetId()]; !exists {
			var s = SignInWrite{
				SignIn:              event,
				hasModelBeenWritten: false,
			}
			if event.GetStatus().GetErrorCode() != nil {
				for _, v := range config.GetSkipErrorCodes() {
					if *event.GetStatus().GetErrorCode() == v {
						addSkips[v]++
						return
					}
				}
				signInActivityMap[*event.GetId()] = s
			}
		}
	}
}

func GetLatestSignInTimestamp() string {
	siLock.Lock()
	defer siLock.Unlock()

	var latest string

	for _, v := range signInActivityMap {
		timestamp := v.SignIn.GetCreatedDateTime().Format("2006-01-02T15:04:05Z")
		if timestamp > latest {
			latest = timestamp
		}
	}
	return latest
}

func trimSignInActivity() trimRecord {
	siLock.Lock()
	defer siLock.Unlock()

	currentTime := time.Now().UTC()
	var trimResult = trimRecord{
		startTime:     currentTime,
		thresholdTime: currentTime.Add(-time.Duration(config.TrimEntryAge) * time.Second),
		count:         0,
	}

	beforeTrim := len(signInActivityMap)
	for k, v := range signInActivityMap {
		if v.SignIn.GetCreatedDateTime().Before(trimResult.thresholdTime) {
			delete(signInActivityMap, k)
		}
	}
	trimResult.count = beforeTrim - len(signInActivityMap)
	return trimResult
}
