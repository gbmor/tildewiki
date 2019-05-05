package main

import (
	"testing"
)

var cssLocalTests = []struct {
	name []byte
	want bool
}{
	{
		name: []byte("https://google.com/test.css"),
		want: false,
	},
	{
		name: []byte("style.css"),
		want: true,
	},
}

// Make sure it's parsing the CSS location correctly
// and returning the correct bool
func Test_cssLocal(t *testing.T) {
	for _, tt := range cssLocalTests {
		t.Run(string(tt.name), func(t *testing.T) {
			if got := cssLocal(tt.name); got != tt.want {
				t.Errorf("cssLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Benchmark_cssLocal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range cssLocalTests {
			cssLocal(c.name)
		}
	}
}
