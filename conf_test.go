package main

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/spf13/viper"
)

var initConfigRegexp = regexp.MustCompile(viper.GetString("ValidPath"))
var initConfigTests = []struct {
	name string
	want *regexp.Regexp
}{
	{
		name: "Config Test",
		want: initConfigRegexp,
	},
}

// Test to make sure it's returning a valid *regexp.Regexp.
// The function does a lot of stuff internally I can't really
// test for right here.
func Test_initConfigParams(t *testing.T) {
	for _, tt := range initConfigTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initConfigParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initConfigParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func Benchmark_initConfigParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range initConfigTests {
			initConfigParams()
		}
	}
}

func Test_error500(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			error500(tt.args.w, tt.args.r)
		})
	}
}

func Test_error404(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			error404(tt.args.w, tt.args.r)
		})
	}
}*/
