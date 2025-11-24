package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	// Pour valider le format de l'URL

	// Pour valider le format de l'URL
	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// Driver SQLite pour GORM
)

// UPDATE : Faire une variable longURLFlag qui stockera la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {

		// UPDATE 1: Valider que le flag --url a été fourni.
		if longURLFlag == "" {
			log.Println("ERREUR: Vous devez fournir une URL avec --url")
			os.Exit(1)
		}

		// UDPATE Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)
		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			log.Printf("ERREUR: URL invalide : %v\n", err)
			os.Exit(1)
		}

		// UPDATED : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatal("Error : configuration is nil")
		}

		// UPDATE : Initialiser la connexion à la base de données SQLite.

		db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Impossible de se connecter à la base SQLite: %v", err)
		}

		// Obtenir l'objet sql.DB sous-jacent pour pouvoir le fermer proprement
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("sql.DB: %v", err)
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Printf(" fermeture de la base de données: %v", err)
			}
		}()

		// UPDATE : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)
		linkService := services.NewLinkService(linkRepo, clickRepo)

		// UPDATE : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		// os.Exit(1) si erreur
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Printf("ERREUR: Impossible de créer le lien : %v\n", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {

	// UPDATE : Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")

	// UPDATE :  Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")

	// UPDATE : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)

}
