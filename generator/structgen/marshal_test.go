package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	structgen "github.com/korylprince/go-adm/generator/structgen"
)

// write drops a schema into a temp file for GenerateFromFiles to read.
func write(t *testing.T, name, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func generate(t *testing.T, path string) string {
	t.Helper()

	buf := new(strings.Builder)
	if err := structgen.GenerateFromFiles([]string{path}, "test", nil, []string{"json"}, buf); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// A schema whose top level payload keys describe a free-form dictionary becomes
// a schema.Map rather than a schema.Struct. The encoder only switched on Enum
// and Struct, so such a schema generated nothing at all -- and reported no
// error, which is the part that makes it easy to miss.
func TestGenerateTopLevelMap(t *testing.T) {
	path := write(t, "anydict.yaml", `
title: Any Dictionary Schema
description: Arbitrary keys.
payloadkeys:
- key: ANY
  type: <any>
  presence: optional
`)

	src := generate(t, path)

	if want := "type AnyDictionarySchema map[string]any"; !strings.Contains(src, want) {
		t.Errorf("expected %q in output:\n%s", want, src)
	}
}

// The ordinary case, to show the map handling didn't displace it.
func TestGenerateTopLevelStruct(t *testing.T) {
	path := write(t, "normal.yaml", `
title: Normal Schema
payloadkeys:
- key: Name
  type: <string>
  presence: required
`)

	src := generate(t, path)

	if want := "type NormalSchema struct"; !strings.Contains(src, want) {
		t.Errorf("expected %q in output:\n%s", want, src)
	}
}
