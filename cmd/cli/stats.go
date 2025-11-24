package cli

import (
	"errors" // ajout de l'import
	"fmt"
	"log"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// UPDATED : variable shortCodeFlag qui stockera la valeur du flag --code
var shortCodeFlag string

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {

		// UPDATED : Valider que le flag --code a été fourni.
		// os.Exit(1) si erreur
		if shortCodeFlag == "" {
			fmt.Fprintln(os.Stderr, "ERREUR: le flag --code est obligatoire. Exemple : url-shortener stats --code=\"xyz123\"")
			os.Exit(1)
		}

		// UPDATED : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg 

		// UPDATED 3: Initialiser la connexion à la BDD.
		// log.Fatalf si erreur
		db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: échec de la connexion à la base de données SQLite: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}


		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}


		// UPDATED S'assurer que la connexion est fermée à la fin de l'exécution de la commande grâce à defer
    	defer sqlDB.Close()

		// UPDATED : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		// linkRepo :=
		// linkService :=
		linkRepo := repository.NewLinkRepository(db)
    	linkService := services.NewLinkService(linkRepo)

		// UPDATED 5: Appeler GetLinkStats pour récupérer le lien et ses statistiques.
		// Attention, la fonction retourne 3 valeurs
		// Pour l'erreur, utilisez gorm.ErrRecordNotFound
		// Si erreur, os.Exit(1)
		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			// On utilise gorm.ErrRecordNotFound pour distinguer le "pas trouvé"
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Fprintf(os.Stderr, "Aucun lien trouvé pour le code court %q\n", shortCodeFlag)
			} else {
				fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des statistiques pour %q : %v\n", shortCodeFlag, err)
			}
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// UPDATED : Définir le flag --code pour la commande stats.
    link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
    if err != nil {
        // On utilise gorm.ErrRecordNotFound pour distinguer le "pas trouvé"
        if errors.Is(err, gorm.ErrRecordNotFound) {
            fmt.Fprintf(os.Stderr, "Aucun lien trouvé pour le code court %q\n", shortCodeFlag)
        } else {
            fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des statistiques pour %q : %v\n", shortCodeFlag, err)
        }
        os.Exit(1)
    }

	// UPDATED Marquer le flag comme requis
    if err := StatsCmd.MarkFlagRequired("code"); err != nil {
        log.Fatalf("FATAL: impossible de marquer le flag --code comme requis: %v", err)
    }

	// UPDATED : Ajouter la commande à RootCmd
    cmd2.RootCmd.AddCommand(StatsCmd)
}