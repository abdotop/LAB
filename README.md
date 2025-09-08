# E-Parrainages - Plateforme de Parrainage Électoral

## Description
Système numérique sécurisé et transparent permettant aux citoyens sénégalais de parrainer des candidats aux élections, avec des mesures anti-fraude et une supervision claire par les autorités électorales.

## Architecture du Projet

Le projet est structuré en 4 applications principales :

### 1. Backend API (`/backend`)
- **Technologie** : Node.js + TypeScript + Prisma + Express
- **Base de données** : PostgreSQL
- **Fonctionnalités** : API centralisée, authentification, RBAC, audit

### 2. DGE Admin Web (`/apps/dge-admin`)
- **Technologie** : Next.js + TypeScript + TailwindCSS
- **Utilisateurs** : Super Admin DGE, Admin DGE standard
- **Fonctionnalités** : Gestion des cycles électoraux, validation des candidatures, supervision

### 3. Mandataires Web (`/apps/mandataires-web`)
- **Technologie** : Next.js + TypeScript + TailwindCSS
- **Utilisateurs** : Gestionnaires de candidature, Mandataires
- **Fonctionnalités** : Dépôt de candidatures, suivi des parrainages

### 4. Application Mobile (`/apps/mobile`)
- **Technologie** : Flutter (Dart)
- **Utilisateurs** : Citoyens sénégalais
- **Fonctionnalités** : Inscription NIN, parrainage sécurisé, consultation des candidats

## Fonctionnalités Principales

### Sécurité
- Authentification par NIN (Numéro d'Identification Nationale)
- Vérification OTP par SMS
- Code PIN à 4 chiffres
- Audit trail immuable
- Chiffrement des données sensibles

### Rôles et Permissions
- **Super Admin DGE** : Accès complet + piste d'audit
- **Admin DGE** : Gestion opérationnelle
- **Gestionnaire de Candidature** : Dépôt de dossiers
- **Mandataire** : Assistance terrain
- **Citoyen** : Parrainage électoral

### Anti-Fraude
- Un parrainage unique par cycle électoral
- Fenêtre de contestation 24h
- Règle du "Joker Unique" pour signalements
- Gel des comptes frauduleux

## Installation et Développement

### Prérequis
- Node.js 18+
- PostgreSQL 14+
- Flutter 3.16+
- Docker (optionnel)

### Installation
```bash
# Cloner le repository
git clone https://github.com/abdotop/LAB.git
cd LAB

# Installer les dépendances backend
cd backend
npm install

# Installer les dépendances DGE Admin
cd ../apps/dge-admin
npm install

# Installer les dépendances Mandataires
cd ../apps/mandataires-web
npm install

# Configuration de la base de données
cd ../../backend
cp .env.example .env
# Éditer .env avec vos paramètres
npx prisma migrate dev
```

### Démarrage en mode développement
```bash
# Backend API
cd backend
npm run dev

# DGE Admin (nouveau terminal)
cd apps/dge-admin
npm run dev

# Mandataires Web (nouveau terminal)
cd apps/mandataires-web
npm run dev

# Mobile App (nouveau terminal)
cd apps/mobile
flutter run
```

## Structure du Projet
```
LAB/
├── backend/                 # API Node.js + TypeScript + Prisma
├── apps/
│   ├── dge-admin/          # Interface web DGE (Next.js)
│   ├── mandataires-web/    # Interface web Mandataires (Next.js)
│   └── mobile/             # Application mobile (Flutter)
├── shared/                  # Types et utilitaires partagés
├── docs/                   # Documentation
└── docker-compose.yml      # Configuration Docker
```

## Contexte Sénégalais
- Interface entièrement en français
- Adaptation aux institutions sénégalaises (DGE, DAF)
- Support des zones à faible connectivité
- Intégration avec les opérateurs SMS locaux

## Licence
MIT - Voir [LICENSE](LICENSE) pour plus de détails.