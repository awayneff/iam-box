// cmd/revoke.go

package cmd

import (
	"context"
	"fmt"

	"iam-box/app/repository"
	"iam-box/app/service"

	"github.com/spf13/cobra"
)

var revokeCmd = &cobra.Command{
	Use:   "revoke [user_id] [action] [resource_type] [resource_id]",
	Short: "Revoke a permission from a user",
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

		err := permissionService.Revoke(context.Background(), userID, action, resourceType, &resourceID)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return
		}

		fmt.Printf("✅ Revoked %s on %s/%s from user %s\n", action, resourceType, resourceID, userID)
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)
}
