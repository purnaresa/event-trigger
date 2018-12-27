package main

import "testing"

func Test_trigger(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trigger(tt.args.event)
		})
	}
}
