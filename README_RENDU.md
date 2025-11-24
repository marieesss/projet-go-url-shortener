README - Rendu du Projet URL Shortener
Travail Réalisé
1. Architecture et Structure
Le projet suit l'architecture modulaire proposée avec une séparation claire des responsabilités :

cmd/ : Commandes CLI avec Cobra (create, stats, migrate, run-server)
internal/ : Logique métier, repositories, services, workers et monitoring
configs/ : Configuration centralisée avec Viper

2. Fonctionnalités Implémentées
Core Features

✅ Raccourcissement d'URLs : Génération de codes courts uniques (6 caractères alphanumériques) avec gestion des collisions via retry logic
✅ Redirection instantanée : Redirection HTTP 302 sans latence
✅ Analytics asynchrones : Enregistrement des clics en arrière-plan via goroutines et channel bufferisé (non-bloquant)
✅ Surveillance d'URLs : Moniteur périodique vérifiant l'accessibilité des URLs avec notifications de changement d'état
✅ API REST : Endpoints complets (health, création, redirection, statistiques)
✅ Interface CLI : Commandes create, stats, migrate et run-server

APIs REST Implémentées

GET /health : Vérification de l'état du service
POST /api/v1/links : Création d'URL courte
GET /{shortCode} : Redirection et analytics asynchrone
GET /api/v1/links/{shortCode}/stats : Récupération des statistiques

3. Technologies Utilisées

Gin : Framework web pour les APIs REST
Cobra : Interface CLI
Viper : Gestion de configuration
GORM : ORM avec SQLite pour la persistance
Goroutines & Channels : Traitement asynchrone des clics

4. Patterns de Conception

Repository Pattern : Abstraction de la couche d'accès aux données
Service Layer : Logique métier isolée
Worker Pool : Pool de goroutines pour traiter les événements de clics
Channel Buffering : Communication non-bloquante entre composants

5. Points Techniques Notables
Concurrence

Channel bufferisé pour les événements de clics (taille configurable)
Pool de workers configurable pour l'enregistrement asynchrone
Multiplexage avec select/default pour éviter le blocage sur channel plein
Mutex pour protéger l'accès concurrent au moniteur d'URLs

Gestion d'Erreurs

Erreurs personnalisées avec wrapping (fmt.Errorf avec %w)
Gestion spécifique de gorm.ErrRecordNotFound
Logging approprié des erreurs critiques

Génération Sécurisée

Utilisation de crypto/rand pour la génération de codes courts imprévisibles
Logique de retry pour garantir l'unicité des codes

6. Configuration
Configuration flexible via configs/config.yaml :

Port du serveur et URL de base
Taille du buffer analytics et nombre de workers
Intervalle de monitoring des URLs
Chemin de la base de données

7. Limitations et Améliorations Possibles

Pas d'URLs personnalisées (feature bonus)
Pas d'expiration de liens (feature bonus)
Pas de rate limiting (feature bonus)
Monitoring basique sans retry pour les clics perdus

Compilation et Tests
bash# Construction
go build -o url-shortener

# Migration
./url-shortener migrate

# Lancement du serveur
./url-shortener run-server

# Création d'URL
./url-shortener create --url="https://example.com"

# Consultation des stats
./url-shortener stats --code="ABC123"


Boris PRINCE AGBODJAN
Marie ESPINOSA
Thomas STECINSKI
Matthéo NAEGELLEN