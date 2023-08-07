package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

func auditLoop(ctx context.Context) {
	signInActivityMap = make(map[string]SignInWrite)
	directoryAuditsMap = make(map[string]DirectoryAuditWrite)

	updateTicker := time.NewTicker(time.Duration(config.GetUpdateInterval()) * time.Second)
	defer updateTicker.Stop()

	trimTicker := time.NewTicker(time.Duration(config.GetTrimEntryAge()) * time.Second)
	defer trimTicker.Stop()

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// kick off initial audit (first run)
	audit(ctx, writer)
	writer.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-trimTicker.C:
			daTrimmed := trimDirectoryAudits()
			logTrimActivity(writer, "directory audit", daTrimmed.startTime, daTrimmed.thresholdTime, daTrimmed.count)
			siTrimmed := trimSignInActivity()
			logTrimActivity(writer, "sign-in activity", siTrimmed.startTime, siTrimmed.thresholdTime, siTrimmed.count)
			writer.Flush()
		case <-updateTicker.C:
			audit(ctx, writer)
			writer.Flush()
		}
	}
}

func audit(ctx context.Context, writer *tabwriter.Writer) {
	cred, err := azidentity.NewClientSecretCredential(config.GetTenantId(), config.GetClientId(), config.GetClientSecret(), nil)
	if err != nil {
		log.Fatal(err)
	}

	client, err := graph.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		fmt.Fprintf(writer, "%v\terror creating client: %v\n", time.Now().UTC().Format("2006-01-02T15:04:05Z"), err)
	} else {

		startTime, currentTime := getTimes()
		directoryAuditFilter := fmt.Sprintf("activityDateTime ge %s and activityDateTime lt %s", startTime, currentTime)
		DirectoryAudit(client, ctx, directoryAuditFilter)
		logAuditActivity(writer, "directory_audit", startTime, currentTime, len(directoryAuditsMap))

		signinActivityFilter := fmt.Sprintf("createdDateTime ge %s and createdDateTime lt %s", startTime, currentTime)
		SignInActivityAudit(client, ctx, signinActivityFilter)
		logAuditActivity(writer, "sign-in activity", startTime, currentTime, len(signInActivityMap))

		writer.Flush()

		daCount, err := WriteDirectoryAuditsToSQL()
		writeToDatabase(writer, daCount, err, "directory audit")

		siCount, err := WriteSignInActivityToSQL()
		writeToDatabase(writer, siCount, err, "signin activity")

		for k, v := range addSkips {
			fmt.Fprintf(writer, "%v\tskipped %d signin activity records with error code %d\n", time.Now().UTC().Format("2006-01-02T15:04:05Z"), v, k)
		}

		addSkips = make(map[int32]int)
	}

}
func getTimes() (string, string) {
	s := time.Now().UTC()
	startTime := s.Add(-time.Duration(config.TrimEntryAge) * time.Second).Format("2006-01-02T15:04:05Z")
	currentTime := s.Format("2006-01-02T15:04:05Z")
	return startTime, currentTime
}

func logAuditActivity(writer io.Writer, eventType string, startTime, currentTime string, eventCount int) {
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	fmt.Fprintf(writer, "%v\t%s events between\t[%s] and [%s]:\t%d\n", timeStamp, eventType, startTime, currentTime, eventCount)
}

func logTrimActivity(writer io.Writer, eventType string, startTime, currentTime time.Time, eventCount int) {
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	fmt.Fprintf(writer, "%v\ttrimmed %d %s events between\t[%s] and [%s]\n", timeStamp, eventCount, eventType, currentTime.Format("2006-01-02T15:04:05Z"), startTime.Format("2006-01-02T15:04:05Z"))
}
func writeToDatabase(writer io.Writer, count int, err error, eventType string) {
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if err != nil {
		fmt.Fprintf(writer, "%v\terror writing %s records to database:\t%v\n", timeStamp, eventType, err)
	} else {
		fmt.Fprintf(writer, "%v\t%s records written to database:\t%d\n", timeStamp, eventType, count)
	}
}
