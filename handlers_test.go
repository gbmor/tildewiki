package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/viper"
)

// This is a pretty strict test. Make sure the
// output of pageHandler is byte-for-byte what
// I'm expecting it to be.
func Test_pageHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "example",
		},
		{
			name: "test1",
		},
	}

	hush, _ := os.Open("/dev/null")
	log.SetOutput(hush)
	genPageCache()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "localhost:8080/w/"+tt.name, nil)
			pageHandler(w, req, tt.name)
			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != 200 {
				t.Errorf("pageHandler(): %v\n", resp.StatusCode)
			}
			if !bytes.Equal(body, cachedPages[tt.name+".md"].Body) {
				t.Errorf("pageHandler(): Byte mismatch\n")
			}
		})
	}
}

// This is the same test type as pageHandler
func Test_indexHandler(t *testing.T) {
	name := "Index Handler Test"
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080", nil)
	t.Run(name, func(t *testing.T) {
		indexHandler(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("indexHandler(): %v\n", resp.StatusCode)
		}
		if !bytes.Equal(body, indexCache.Body) {
			t.Errorf("indexHandler(): Byte mismatch\n")
		}
	})
}

// This is the same test type as pageHandler
func Test_iconHandler(t *testing.T) {
	name := "Icon Handler Test"
	icon, _ := ioutil.ReadFile(viper.GetString("AssetsDir") + "/" + viper.GetString("Icon"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080/icon", nil)
	t.Run(name, func(t *testing.T) {
		iconHandler(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("iconHandler(): %v\n", resp.StatusCode)
		}
		if !bytes.Equal(body, icon) {
			t.Errorf("iconHandler(): Byte mismatch\n")
		}
	})
}

// This is the same test type as pageHandler
func Test_cssHandler(t *testing.T) {
	name := "CSS Handler Test"
	if !cssLocal([]byte(viper.GetString("CSS"))) {
		t.Skipf("cssHandler(): Set to use remote CSS in config, skipping test ...\n")
	}
	css, _ := ioutil.ReadFile(viper.GetString("CSS"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080/css", nil)
	t.Run(name, func(t *testing.T) {
		cssHandler(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("cssHandler(): %v\n", resp.StatusCode)
		}
		if !bytes.Equal(body, css) {
			t.Errorf("cssHandler(): Byte mismatch\n")
		}
	})
}

// this test is a bit hairy, will finish soon
/*
func Test_validatePath(t *testing.T) {
	w := httptest.NewRecorder()
	type args struct {
		fn func(*httptest.ResponseRecorder, *http.Request, string)
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		{
			name: "valid",
			args: func(w, httptest.NewRequest("GET", "localhost:8080/w/example", nil), "example"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePath(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validatePath() = %v, want %v", got, tt.want)
			}
		})
	}

}*/
