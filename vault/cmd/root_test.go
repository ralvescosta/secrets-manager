package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_Execute_should_execute_correctly(t *testing.T) {
	cmd := &cobra.Command{}
	Execute(cmd)
}
