package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"iam-box/app/repository"
	"iam-box/app/service"

	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listOffset int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all permissions",
	Run: func(cmd *cobra.Command, args []string) {
		db := initDB()
		permissionRepo := repository.NewPermissionRepository(db)
		decisionRepo := repository.NewDecisionRepository(db)
		permissionService := service.NewPermissionService(*permissionRepo, *decisionRepo)

		limit := listLimit
		offset := listOffset

		if limit < 1 || limit > 1000 {
			fmt.Println("limit must be between 1 and 1000")
			return
		}
		if offset < 0 {
			offset = 0
		}

		perms, err := permissionService.List(context.Background(), limit, offset)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return
		}

		if len(*perms) == 0 {
			fmt.Println("No permissions found")
			return
		}

		// Use tabwriter from stdlib
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tUser\tAction\tType\tResource")

		for _, p := range *perms {
			resourceID := ""
			if p.ResourceID != nil {
				resourceID = *p.ResourceID
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				p.ID,
				p.UserID,
				p.Action,
				p.ResourceType,
				resourceID)
		}
		w.Flush()
	},
}

func init() {
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "max results (1-1000)")
	listCmd.Flags().IntVar(&listOffset, "offset", 0, "number to skip")
	rootCmd.AddCommand(listCmd)
}
