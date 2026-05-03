package cmd

import (
	"fmt"
	"iam-box/app/entities"
	"log"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost     string
	dbUser     string
	dbPassword string
	dbName     string
	dbPort     string
	dbSSLMode  string
)

var rootCmd = &cobra.Command{
	Use:   "iam-box",
	Short: "IAM authorization service",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbHost, "db-host", "localhost", "database host")
	rootCmd.PersistentFlags().StringVar(&dbUser, "db-user", "iam_user", "database user")
	rootCmd.PersistentFlags().StringVar(&dbPassword, "db-password", "iam_password", "database password")
	rootCmd.PersistentFlags().StringVar(&dbName, "db-name", "iam_db", "database name")
	rootCmd.PersistentFlags().StringVar(&dbPort, "db-port", "5432", "database port")
	rootCmd.PersistentFlags().StringVar(&dbSSLMode, "db-sslmode", "disable", "database ssl mode")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func initDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&entities.Permission{}, &entities.Decision{})
	return db
}
