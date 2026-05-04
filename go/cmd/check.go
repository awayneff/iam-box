// cmd/check.go

package cmd

import (
	"context"
	"fmt"

	"iam-box/app/repository"
	"iam-box/app/service"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [user_id] [action] [resource_type] [resource_id]",
	Short: "Check a permission",
	Args:  cobra.RangeArgs(3, 4), // 3 or 4 args
	Run: func(cmd *cobra.Command, args []string) {
		db := initDB()
		permissionRepo := repository.NewPermissionRepository(db)
		decisionRepo := repository.NewDecisionRepository(db)
		permissionService := service.NewPermissionService(*permissionRepo, *decisionRepo)

		userID := args[0]
		action := args[1]
		resourceType := args[2]

		// accept wildcard
		var resourceID *string
		if len(args) == 4 && args[3] != "" && args[3] != "null" {
			resourceID = &args[3]
		}

		allowed, err := permissionService.Check(context.Background(), userID, action, resourceType, resourceID)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		if allowed {
			fmt.Println("✅ ALLOWED")
		} else {
			fmt.Println("❌ DENIED")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
