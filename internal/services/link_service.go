package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// UPDATED Créer la struct
// LinkService est une structure qui g fournit des méthodes pour la logique métier des liens.
// Elle détient linkRepo qui est une référence vers une interface LinkRepository.
// IMPORTANT : Le champ doit être du type de l'interface (non-pointeur).

type LinkService struct {
	linkRepo  repository.LinkRepository
	clickRepo repository.ClickRepository
}

// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// UPDATED Créer la méthode GenerateShortCode
// GenerateShortCode est une méthode rattachée à LinkService
// Elle génère un code court aléatoire d'une longueur spécifiée. Elle prend une longueur en paramètre et retourne une string et une erreur
// Il utilise le package 'crypto/rand' pour éviter la prévisibilité.
// Je vous laisse chercher un peu :) C'est faisable en une petite dizaine de ligne

// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {

	var shortCode string
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		// Générer un shortCode
		code, err := s.GenerateShortCode(6)
		if err != nil {
			return nil, err
		}

		existing, err := s.linkRepo.GetLinkByShortCode(code)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code
				break
			}
			return nil, fmt.Errorf("error checking shortCode: %w", err)
		}

		if existing != nil {
			log.Printf("Erreur sur '%s', tentative %d/%d...", code, i+1, maxRetries)
			continue
		}

		shortCode = code
		break
	}

	if shortCode == "" {
		return nil, errors.New("failed to generate a unique shortCode after retries")
	}

	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
		IsUp:      true,
	}

	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("failed to save link: %w", err)
	}

	return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get link for shortcode %s: %w", shortCode, err)
	}
	return link, nil
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	link, err := s.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.clickRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks: %w", err)
	}

	return link, count, nil
}
