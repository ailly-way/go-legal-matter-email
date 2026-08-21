package main

import (
	"fmt"
	"os"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	matter := Matter{ID: "M-1042", ClientEmail: os.Getenv("MATTER_CLIENT_EMAIL"), DocumentURL: os.Getenv("SIGNED_DOCUMENT_URL"), DueDays: 3}
	if matter.ClientEmail == "" || matter.DocumentURL == "" {
		fmt.Fprintln(os.Stderr, "MATTER_CLIENT_EMAIL and SIGNED_DOCUMENT_URL are required")
		os.Exit(1)
	}
	result, err := DeliverSignedDocument(client, matter)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	followUp, err := SendDeadlineFollowUp(client, matter)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("matter %s: delivery=%s follow_up=%s\n", matter.ID, result.MessageID, followUp.MessageID)
}
