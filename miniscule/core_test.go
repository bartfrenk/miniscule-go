package miniscule

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestEnvResolver(t *testing.T) {
	var s interface{}
	err := Unmarshal([]byte("!env HELLO"), &s)
	if err != nil {
		t.Errorf("Failed to unmarshall '!env HELLO': %s", err)
	}
	if s != "" {
		t.Errorf("Expected nil, got %s", s)
	}
}

func TestEnvResolver2(t *testing.T) {
	var out yaml.Node
	var in = []byte("- !env HELLO")
	err := Unmarshal(in, &out)
	if err != nil {
		t.Errorf("Failed to unmarshall %s", in)
	}
	if len(out.Content) != 1 || out.Content[0].Value != "" {
		t.Error("Bla")
	}
}

func TestEnvResolver3(t *testing.T) {
	var out yaml.Node
	var in = []byte("x: !env HELLO")
	err := Unmarshal(in, &out)
	if err != nil {
		t.Errorf("Failed to unmarshall %s", in)
	}
	s, err := yaml.Marshal(out)
	t.Errorf("%s", s)
}
