package resolver_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/melvinsh/cnay/pkg/resolver"
)

func TestGetInputReader(t *testing.T) {
	testCases := []struct {
		listFile string
		wantErr  bool
	}{
		{listFile: "", wantErr: true},
		{listFile: "../../test/hosts.txt", wantErr: false},
		{listFile: "nonexistentfile.txt", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.listFile, func(t *testing.T) {
			reader, err := resolver.GetInputReader(tc.listFile, false)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetInputReader(%q): expected an error but got none", tc.listFile)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetInputReader(%q): unexpected error: %v", tc.listFile, err)
			}
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(reader)
			got := strings.TrimSpace(buf.String())
			want := "setting.bitfinex.com\nh.bitfinex.com\neos-bp.bitfinex.com"
			if got != want {
				t.Errorf("GetInputReader(%q): got %q, want %q", tc.listFile, got, want)
			}
		})
	}
}

func TestReadHostnames(t *testing.T) {
	input := `example.com
foo.bar
test.example.org`
	want := []string{"example.com", "foo.bar", "test.example.org"}

	reader := strings.NewReader(input)
	got := resolver.ReadHostnames(reader, false)

	if len(got) != len(want) {
		t.Fatalf("ReadHostnames(): got %d hostnames, want %d", len(got), len(want))
	}

	for i, hostname := range want {
		if got[i] != hostname {
			t.Errorf("ReadHostnames(): got %q at position %d, want %q", got[i], i, hostname)
		}
	}
}

func TestResolveHostnames(t *testing.T) {
	testCases := []struct {
		name           string
		hostnames      []string
		debug          bool
		showHostname   bool
		useProgressBar bool
		want           []string
	}{
		{
			name:      "single hostname, don't show hostnames",
			hostnames: []string{"example.com"},
			want:      []string{"93.184.216.34"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolver.ResolveHostnames(tc.hostnames, tc.debug, tc.showHostname, tc.useProgressBar)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ResolveHostnames(%v, %v, %v, %v) = %v, want %v", tc.hostnames, tc.debug, tc.showHostname, tc.useProgressBar, got, tc.want)
			}
		})
	}
}
