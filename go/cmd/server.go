package cmd

import (
    "fmt"
    "iam-box/app/controllers"
    "iam-box/app/repository"
    "iam-box/app/service"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/spf13/cobra"
)

var (
    port string
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Start IAM API server",
    Run: func(cmd *cobra.Command, args []string) {
        db := initDB()
        
        // Repositories
        permissionRepo := repository.NewPermissionRepository(db)
        decisionRepo := repository.NewDecisionRepository(db)
        
        // Services
        permissionService := service.NewPermissionService(*permissionRepo, *decisionRepo)
        decisionService := service.NewDecisionService(*decisionRepo)
        
        // Controllers
        permissionsController := controllers.NewPermissionController(*permissionService)
        decisionController := controllers.NewDecisionController(*decisionService)
        
        // Router
        r := chi.NewRouter()
        r.Route("/api/v1", func(r chi.Router) {
            r.Route("/permissions", func(r chi.Router) {
                r.Post("/grant", permissionsController.Create)
                r.Get("/{user_id}", permissionsController.GetByUser)
                r.Post("/can", permissionsController.Check)
                r.Delete("/", permissionsController.Delete)
            })
            r.Route("/decisions", func(r chi.Router) {
                r.Get("/", decisionController.List)
            })
        })
        
        r.Get("/_core/health", func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"status": "ok"}`))
        })
        
        addr := ":" + port
        fmt.Printf("IAM server listening on %s\n", addr)
        log.Fatal(http.ListenAndServe(addr, r))
    },
}

func init() {
    serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to listen on")
    rootCmd.AddCommand(serverCmd)
}