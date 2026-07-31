package config

import (
	"testing"

	"github.com/one2nc/cloudlens/internal"
)

func TestEC2SnapshotAliases(t *testing.T) {
	aliases := NewAliases()
	aliases.loadDefaultAliases(internal.AWS)

	for _, alias := range []string{"ec2:s", "ec2:S", "Ec2:S"} {
		got, ok := aliases.Get(alias)
		if !ok {
			t.Errorf("alias %q not registered", alias)
			continue
		}
		if got != "ec2:s" {
			t.Errorf("alias %q resolved to %q, want ec2:s", alias, got)
		}
	}
}
