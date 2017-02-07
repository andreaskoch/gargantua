package main

import (
	"strings"
	"testing"
)

func Test_getLatestLogMessages_NoMessagesGiven_EmptyArrayIsReturned(t *testing.T) {

	messages := []string{}

	result, _ := getLatestLogMessages(messages, 3)

	if len(result) != 3 {
		t.Fail()
		t.Logf("getLatestLogMessages([...], 3) should have returned and array of size 3")
	}

	if result[0] != "" {
		t.Fail()
		t.Logf("getLatestLogMessages([\"%s\"], 3) returned %q instead of %q", strings.Join(messages, "\",\""), result[0], "")
	}
}

func Test_getLatestLogMessages_RequestedCountEqualsNumberOfMessages(t *testing.T) {

	messages := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}

	result, _ := getLatestLogMessages(messages, 3)

	if len(result) != 3 {
		t.Fail()
		t.Logf("getLatestLogMessages([...], 3) should have returned and array of size 3")
	}

	if result[2] != "Line 3" {
		t.Fail()
		t.Logf("getLatestLogMessages([\"%s\"], 3) returned %q instead of %q", strings.Join(messages, "\",\""), result[0], "Line 3")
	}
}

func Test_getLatestLogMessages_RequestedCountIsSmallerThanNumberOfMessages(t *testing.T) {

	messages := []string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
	}

	result, _ := getLatestLogMessages(messages, 3)

	if len(result) != 3 {
		t.Fail()
		t.Logf("getLatestLogMessages([...], 3) should have returned and array of size 3")
	}

	if result[0] != "Line 3" {
		t.Fail()
		t.Logf("getLatestLogMessages([\"%s\"], 3) returned %q instead of %q", strings.Join(messages, "\",\""), result[0], "Line 3")
	}
}
