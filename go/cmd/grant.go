package cmd

import (
	"context"
	"fmt"

	"iam-box/app/repository"
	"iam-box/app/service"

	"github.com/spf13/cobra"
)

var grantCmd = &cobra.Command{
	Use:   "grant [user_id] [action] [resource_type] [resource_id]",
	Short: "Grant a permission",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		db := initDB()
		permissionRepo := repository.NewPermissionRepository(db)
		decisionRepo := repository.NewDecisionRepository(db)
		permissionService := service.NewPermissionService(*permissionRepo, *decisionRepo)

		userID := args[0]
		action := args[1]
		resourceType := args[2]
		resourceID := args[3]

		err := permissionService.Create(context.Background(), userID, action, resourceType, &resourceID)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return
		}

		fmt.Printf("✅ Granted %s %s on %s/%s\n", userID, action, resourceType, resourceID)
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)
}
