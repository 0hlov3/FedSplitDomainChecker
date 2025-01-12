package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCmdFlags(t *testing.T) {
	tests := []struct {
		flagName    string
		defaultVal  string
		expectedVal string
	}{
		{"accountDomain", "https://gotosocial.org", "https://example.org"},
		{"hostDomain", "https://gts.gotosocial.org", "https://example-host.org"},
		{"account", "admin@gotosocial.org", "test@example.org"},
	}

	for _, test := range tests {
		t.Run(test.flagName, func(t *testing.T) {
			viper.Set(test.flagName, test.expectedVal)

			flag := rootCmd.PersistentFlags().Lookup(test.flagName)
			assert.NotNil(t, flag, "Flag should be registered")
			assert.Equal(t, test.defaultVal, flag.DefValue, "Default value should match")
			assert.Equal(t, test.expectedVal, viper.GetString(test.flagName), "Expected value should match")
		})
	}
}

func TestRootCmdExecution(t *testing.T) {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)

	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	assert.NoError(t, err, "Execute should not return an error")
}

func TestPersistentPreRun(t *testing.T) {
	called := false

	originalPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		called = true
	}

	dummyCmd := &cobra.Command{
		Use:   "dummy",
		Short: "A dummy command for testing",
		Run:   func(cmd *cobra.Command, args []string) {}, // No-op
	}
	rootCmd.AddCommand(dummyCmd)

	rootCmd.SetArgs([]string{"dummy"})
	err := rootCmd.Execute()

	rootCmd.RemoveCommand(dummyCmd)
	rootCmd.PersistentPreRun = originalPreRun

	assert.NoError(t, err, "Execute should not return an error")
	assert.True(t, called, "PersistentPreRun should be called")
}
