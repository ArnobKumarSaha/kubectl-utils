package rbac

import (
	"github.com/spf13/cobra"
)

var (
	name      string
	namespace string
	oyaml     bool
)

func NewCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "rbac",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	addCommonFlags(cmd)
	cmd.AddCommand(ServiceAccountCMD())
	cmd.AddCommand(RoleCMD())
	cmd.AddCommand(ClusterROleCMD())
	return cmd
}
