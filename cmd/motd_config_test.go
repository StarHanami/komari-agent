package cmd

import (
	"encoding/json"
	"testing"
)

func TestDisableMOTDConfiguration(t *testing.T) {
	saved := *flags
	defer func() { *flags = saved }()
	flag := RootCmd.PersistentFlags().Lookup("disable-motd")
	if flag == nil || flag.DefValue != "false" {
		t.Fatal("--disable-motd must exist and default to false")
	}
	wasChanged := flag.Changed
	defer func() { flag.Changed = wasChanged }()

	t.Run("default", func(t *testing.T) {
		if flags.DisableMOTD {
			t.Fatal("MOTD modifications must remain enabled by default")
		}
	})
	t.Run("command line", func(t *testing.T) {
		if err := RootCmd.PersistentFlags().Set("disable-motd", "true"); err != nil {
			t.Fatal(err)
		}
		if !flags.DisableMOTD {
			t.Fatal("CLI option did not disable MOTD modifications")
		}
	})
	for _, value := range []string{"true", "1"} {
		t.Run("environment "+value, func(t *testing.T) {
			flags.DisableMOTD = false
			t.Setenv("AGENT_DISABLE_MOTD", value)
			loadFromEnv()
			if !flags.DisableMOTD {
				t.Fatal("environment option did not disable MOTD modifications")
			}
		})
	}
	t.Run("JSON overrides earlier sources", func(t *testing.T) {
		for _, value := range []bool{false, true} {
			flags.DisableMOTD = !value
			raw, err := json.Marshal(map[string]bool{"disable_motd": value})
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(raw, flags); err != nil {
				t.Fatal(err)
			}
			if flags.DisableMOTD != value {
				t.Fatalf("JSON disable_motd = %v, want %v", flags.DisableMOTD, value)
			}
		}
	})
}
