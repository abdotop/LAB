# Documentation E-Parrainages

## Vue d'ensemble
Cette documentation couvre la plateforme complète de parrainage électoral pour le Sénégal.

## Architecture du Projet

```
LAB/
├── backend/                 # API Node.js + TypeScript + Prisma
│   ├── src/
│   │   ├── controllers/     # Contrôleurs API
│   │   ├── middleware/      # Middlewares Express
│   │   ├── models/          # Modèles métier
│   │   ├── routes/          # Définition des routes
│   │   ├── services/        # Services métier
│   │   └── utils/           # Utilitaires
│   ├── prisma/             # Schéma et migrations base de données
│   └── Dockerfile          # Image Docker backend
├── apps/
│   ├── dge-admin/          # Interface web DGE (Next.js)
│   ├── mandataires-web/    # Interface web Mandataires (Next.js)
│   └── mobile/             # Application mobile (Flutter)
├── shared/                 # Types et utilitaires partagés
├── docs/                   # Documentation
├── .github/workflows/      # CI/CD GitHub Actions
└── docker-compose.yml      # Configuration Docker
```

## Fonctionnalités Principales

### Authentification et Sécurité
- Inscription par NIN (Numéro d'Identification Nationale)
- Vérification OTP par SMS
- Code PIN à 4 chiffres pour les actions sensibles
- JWT pour l'authentification API
- RBAC (Role-Based Access Control)

### Rôles Utilisateurs
1. **Super Admin DGE**
   - Accès complet à toutes les fonctionnalités
   - Consultation de la piste d'audit immuable
   - Gestion des cycles électoraux
   - Validation des candidatures

2. **Admin DGE**
   - Gestion opérationnelle
   - Pas d'accès à la piste d'audit

3. **Gestionnaire de Candidature**
   - Dépôt de dossiers de candidature
   - Suivi du statut de validation

4. **Mandataire**
   - Assistance terrain
   - Suivi des parrainages

5. **Citoyen**
   - Parrainage électoral
   - Consultation des candidats

### Cycle de Vie d'un Parrainage

1. **Inscription Citoyen**
   - Saisie du NIN
   - Vérification OTP SMS
   - Définition du code PIN

2. **Consultation des Candidats**
   - Liste des candidats validés
   - Informations sur les candidatures

3. **Acte de Parrainage**
   - Sélection du candidat
   - Triple confirmation dont code PIN
   - Enregistrement sécurisé

4. **Fenêtre de Contestation**
   - 24h pour contester
   - Signalement de fraude limité (règle du "Joker Unique")

## Installation et Déploiement

### Développement Local

```bash
# Cloner le repository
git clone https://github.com/abdotop/LAB.git
cd LAB

# Backend
cd backend
npm install
cp .env.example .env
# Configurer la base de données dans .env
npx prisma migrate dev
npm run dev

# DGE Admin
cd ../apps/dge-admin
npm install
npm run dev

# Mandataires Web
cd ../apps/mandataires-web
npm install
npm run dev

# Mobile
cd ../apps/mobile
flutter pub get
flutter run
```

### Déploiement avec Docker

```bash
# Démarrer tous les services
docker-compose up -d

# Accès aux applications
# DGE Admin: http://localhost:3000
# Mandataires: http://localhost:3002
# API: http://localhost:3001
```

## Configuration de Base de Données

Le schéma Prisma définit toutes les tables nécessaires :
- `users` : Utilisateurs et leurs rôles
- `electoral_cycles` : Cycles électoraux
- `candidacies` : Candidatures
- `sponsorships` : Parrainages
- `audit_logs` : Piste d'audit immuable
- `otp_sessions` : Sessions OTP
- `fraud_reports` : Signalements de fraude

## API Endpoints

### Authentification
- `POST /api/auth/register` : Inscription
- `POST /api/auth/login` : Connexion
- `POST /api/auth/verify-otp` : Vérification OTP
- `POST /api/auth/logout` : Déconnexion

### Cycles Électoraux
- `GET /api/electoral-cycles` : Liste des cycles
- `POST /api/electoral-cycles` : Créer un cycle (Admin DGE)
- `GET /api/electoral-cycles/:id` : Détails d'un cycle

### Candidatures
- `GET /api/candidacies` : Liste des candidatures
- `POST /api/candidacies` : Déposer une candidature
- `PUT /api/candidacies/:id/approve` : Approuver (Admin DGE)

### Parrainages
- `GET /api/sponsorships` : Mes parrainages
- `POST /api/sponsorships` : Parrainer un candidat
- `POST /api/sponsorships/:id/contest` : Contester

### Audit
- `GET /api/audit` : Piste d'audit (Super Admin DGE uniquement)

## Contexte Sénégalais

### Langue
Toutes les interfaces sont en français, langue officielle du Sénégal.

### Institutions
- **DGE** : Direction Générale des Élections
- **DAF** : Direction Administrative et Financière

### Connectivité
L'application mobile est optimisée pour les zones à faible connectivité avec :
- Assistance mandataire pour scan CNI
- Interface simplifiée
- Cache local des données essentielles

## Sécurité et Anti-Fraude

### Mesures de Sécurité
- Chiffrement des données sensibles
- Audit trail immuable
- Rate limiting sur les API
- Validation stricte des données

### Anti-Fraude
- Un seul parrainage par citoyen et par cycle
- Règle du "Joker Unique" pour les signalements
- Gel automatique des comptes suspects
- Géolocalisation des actions

## Tests et Qualité

### Tests Backend
```bash
cd backend
npm test
```

### Tests Frontend
```bash
cd apps/dge-admin
npm run lint
npm run build

cd ../mandataires-web
npm run lint
npm run build
```

### Tests Mobile
```bash
cd apps/mobile
flutter analyze
flutter test
```

## CI/CD

Le pipeline GitHub Actions automatise :
- Tests de tous les composants
- Build des applications
- Déploiement Docker
- Contrôles de qualité du code

## Support et Contribution

Pour toute question ou contribution, veuillez consulter :
- Issues GitHub pour les bugs
- Pull Requests pour les contributions
- Documentation technique dans `/docs`