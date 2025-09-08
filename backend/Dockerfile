# Dockerfile pour le backend Node.js
FROM node:18-alpine

WORKDIR /app

# Copier les fichiers package
COPY package*.json ./

# Installer les dépendances
RUN npm ci --only=production

# Copier le code source
COPY . .

# Générer le client Prisma
RUN npx prisma generate

# Build de l'application
RUN npm run build

# Créer un utilisateur non-root
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nodejs -u 1001

# Créer les dossiers nécessaires
RUN mkdir -p /app/uploads /app/logs
RUN chown -R nodejs:nodejs /app

USER nodejs

EXPOSE 3001

# Démarrer l'application
CMD ["npm", "start"]