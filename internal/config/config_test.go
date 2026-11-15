package config_test

import (
	"bytes"
	"testing"

	"git.sr.ht/~tjex/dotf/internal/config"
	toml "github.com/pelletier/go-toml/v2"
)

func TestModuleAcceptsStringAndTable(t *testing.T) {
	src := `paths = ["~/a", { path = "~/b", repo = "bb" }, { path = "~/c" }]`
	var mc config.ModuleConfig
	dec := toml.NewDecoder(bytes.NewReader([]byte(src)))
	dec.EnableUnmarshalerInterface()
	if err := dec.Decode(&mc); err != nil {
		t.Fatal(err)
	}
	want := []config.Module{{Path: "~/a"}, {Path: "~/b", Repo: "bb"}, {Path: "~/c"}}
	if len(mc.Paths) != len(want) {
		t.Fatalf("got %v", mc.Paths)
	}
	for i := range want {
		if mc.Paths[i] != want[i] {
			t.Errorf("[%d] got %+v, want %+v", i, mc.Paths[i], want[i])
		}
	}
}
