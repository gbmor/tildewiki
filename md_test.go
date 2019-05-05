package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

var mdTestData1, _ = ioutil.ReadFile("pages/example.md")
var mdTestData2, _ = ioutil.ReadFile("pages/test1.md")
var markdownTests = []struct {
	name  string
	css   string
	title string
	data  []byte
}{
	{
		name:  "one",
		css:   "assets/wiki.css",
		title: "Example Page",
		data:  mdTestData1,
	},
	{
		name:  "two",
		css:   "assets/wiki.css",
		title: "No Description",
		data:  mdTestData2,
	},
}

// Make sure setupMarkdown is returning a valid
// blackfriday.HTMLRenderer type
func Test_setupMarkdown(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			var got interface{} = setupMarkdown(tt.css, tt.title)
			if _, ok := got.(*bf.HTMLRenderer); !ok || got == nil {
				t.Errorf("setupMarkdown() returned incorrect type or is nil: %v", reflect.TypeOf(got))
			}
		})
	}
}
func Benchmark_setupMarkdown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range markdownTests {
			setupMarkdown(c.css, c.title)
		}
	}
}

// Previously, I was using bytes.Equal(a, b) to test the
// output of render. However, I can't control for variations
// in blackfriday's output, so I'm just testing to make sure
// it's returning *something*
func Test_render(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			var got []byte
			if got = render(tt.data, tt.css, tt.title); got == nil {
				t.Errorf("render() outputting nil bytes\n")
			}
		})
	}
}
func Benchmark_render(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range markdownTests {
			render(c.data, c.css, c.title)
		}
	}
}
