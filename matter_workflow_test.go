package main

import "testing"

func TestFollowUpSubjectMakesUrgencyVisible(t *testing.T) {
	matter := Matter{ID: "M-1042", DueDays: 0}
	got := FollowUpSubject(matter)
	want := "Action required: signed document"
	if got != want {
		t.Fatalf("subject for a due matter = %q, want %q", got, want)
	}
}
