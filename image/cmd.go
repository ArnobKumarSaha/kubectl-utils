package image

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	name      string
	namespace string
	resource  string
	hash      bool
)

func NewCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Example: `kubectl utils image list -r ds,dep -n kube-system`,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	addCommonFlags(cmd)
	cmd.AddCommand(NewListCmd())
	return cmd
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&name, "name", name, "")
	// TODO: Need to take the `name` into account.
	// Possibly, If name is given , show the containerName & imageId.
	// Otherwise, show the images only.    NEAT & CLEAN
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", namespace, "")
	cmd.PersistentFlags().StringVarP(&resource, "resource", "r", resource, "")

	cmd.PersistentFlags().BoolVarP(&hash, "hash", "", hash, "shows image hash too")
	cmd.PersistentFlags().Lookup("hash").NoOptDefVal = "true"
	_ = cmd.MarkFlagRequired("resource")
}

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			err := listAll()
			if err != nil {
				_ = fmt.Errorf("error on listing : %v", err)
			}
		},
	}
	cmd.AddCommand()
	return cmd
}
