package main

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	testMessageSet1 := MessageSet{
		SetName: "foo",
		Messages: []LoadingMessage{
			{
				Text:       "foo1",
				MinSeconds: 1,
				MaxSeconds: 1,
			},
			{
				Text:       "foo2",
				MinSeconds: 2,
				MaxSeconds: 2,
			},
		},
	}

	testMessageSet2 := MessageSet{
		SetName: "bar",
		Messages: []LoadingMessage{
			{
				Text:       "bar1",
				MinSeconds: 1,
				MaxSeconds: 1,
			},
			{
				Text:       "bar2",
				MinSeconds: 2,
				MaxSeconds: 2,
			},
		},
	}

	testFile := MessageFile{
		Sets: []MessageSet{
			testMessageSet1,
			testMessageSet2,
		},
	}

	marshalled, err := json.Marshal(testFile)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	parsed1, err := parseConfig(marshalled, testMessageSet1.SetName)
	if err != nil {
		t.Errorf("error parsing first test message set: %v", err)
	}

	if !reflect.DeepEqual(parsed1, &testMessageSet1) {
		diff := cmp.Diff(&testMessageSet1, parsed1)
		t.Errorf("first message set parsed incorrectly: %v\n", diff)
	}

	parsed2, err := parseConfig(marshalled, testMessageSet2.SetName)
	if err != nil {
		t.Errorf("error parsing first test message set: %v", err)
	}

	if !reflect.DeepEqual(parsed2, &testMessageSet2) {
		diff := cmp.Diff(&testMessageSet2, parsed2)
		t.Errorf("second message set parsed incorrectly: %v\n", diff)
	}

}
