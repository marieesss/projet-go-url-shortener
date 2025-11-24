package cli

import (
	"fmt"
	"log"

	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// Driver SQLite pour GORM
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmdCobra *cobra.Command, args []string) {
		// UPDATED : Charger la configuration chargée globalement via cmd.cfg

		cfg := cmd.Cfg
		if cfg == nil {
			log.Fatal("Error config")
		}

		// UPDATED 2: Initialiser la connexion à la BDD

		db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erreur connexion SQLite: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Erreur  SQL: %v", err)
		}
		// UPDATED Assurez-vous que la connexion est fermée après la migration grâce à defer

		defer sqlDB.Close()

		// TODO 3: Exécuter les migrations automatiques de GORM.
		// Utilisez db.AutoMigrate() et passez-lui les pointeurs vers tous vos modèles.

		err = db.AutoMigrate(
			&models.Link{},
			&models.Click{},
		)
		if err != nil {
			log.Fatalf("Erreur lors des migrations GORM: %v", err)
		}

		// Pas touche au log
		fmt.Println("Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	// UPDZTED : Ajouter la commande à RootCmd
	cmd.RootCmd.AddCommand(MigrateCmd)
}
