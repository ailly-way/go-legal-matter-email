package main

import "fmt"

type Matter struct {
	ID          string
	ClientEmail string
	DocumentURL string
	DueDays     int
}

func FollowUpSubject(m Matter) string {
	if m.DueDays <= 0 {
		return "Action required: signed document"
	}
	return fmt.Sprintf("Matter %s: signature due in %d days", m.ID, m.DueDays)
}

func DeliverSignedDocument(c *InfraiClient, m Matter) (sendResult, error) {
	body := fmt.Sprintf("<p>Your signed document for matter %s is ready.</p><p>Download: %s</p>", m.ID, m.DocumentURL)
	return c.SendEmail(m.ClientEmail, "Signed document delivery", body, "matter-"+m.ID+"-delivery")
}

func SendDeadlineFollowUp(c *InfraiClient, m Matter) (sendResult, error) {
	body := fmt.Sprintf("<p>Please review the deadline for matter %s.</p><p>%s</p>", m.ID, FollowUpSubject(m))
	return c.SendEmail(m.ClientEmail, FollowUpSubject(m), body, "matter-"+m.ID+"-follow-up")
}
