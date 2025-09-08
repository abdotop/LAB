import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import rateLimit from 'express-rate-limit';
import dotenv from 'dotenv';
import { PrismaClient } from '@prisma/client';

// Import des routes
import authRoutes from './routes/auth';
import userRoutes from './routes/users';
import electoralCycleRoutes from './routes/electoral-cycles';
import candidacyRoutes from './routes/candidacies';
import sponsorshipRoutes from './routes/sponsorships';
import auditRoutes from './routes/audit';

// Import des middlewares
import { errorHandler } from './middleware/errorHandler';
import { logger } from './utils/logger';

// Configuration
dotenv.config();

const app = express();
const prisma = new PrismaClient();
const PORT = process.env.PORT || 3001;

// Middlewares de sécurité
app.use(helmet());

// CORS
const corsOptions = {
  origin: process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:3000'],
  credentials: true,
};
app.use(cors(corsOptions));

// Rate limiting
const limiter = rateLimit({
  windowMs: parseInt(process.env.RATE_LIMIT_WINDOW_MS || '900000'), // 15 minutes
  max: parseInt(process.env.RATE_LIMIT_MAX_REQUESTS || '100'),
  message: 'Trop de requêtes depuis cette IP, veuillez réessayer plus tard.',
});
app.use(limiter);

// Logging
app.use(morgan('combined', { stream: { write: message => logger.info(message.trim()) } }));

// Parsing
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Health check
app.get('/health', (req, res) => {
  res.json({ 
    status: 'OK', 
    timestamp: new Date().toISOString(),
    version: '1.0.0',
    service: 'e-parrainages-backend'
  });
});

// Routes API
app.use('/api/auth', authRoutes);
app.use('/api/users', userRoutes);
app.use('/api/electoral-cycles', electoralCycleRoutes);
app.use('/api/candidacies', candidacyRoutes);
app.use('/api/sponsorships', sponsorshipRoutes);
app.use('/api/audit', auditRoutes);

// Route 404
app.use('*', (req, res) => {
  res.status(404).json({ 
    error: 'Route non trouvée',
    path: req.originalUrl 
  });
});

// Gestionnaire d'erreurs global
app.use(errorHandler);

// Démarrage du serveur
const startServer = async () => {
  try {
    // Test de la connexion à la base de données
    await prisma.$connect();
    logger.info('Connexion à la base de données établie');

    app.listen(PORT, () => {
      logger.info(`🚀 Serveur démarré sur le port ${PORT}`);
      logger.info(`🌍 Environnement: ${process.env.NODE_ENV || 'development'}`);
      logger.info(`📖 API Documentation: http://localhost:${PORT}/health`);
    });
  } catch (error) {
    logger.error('Erreur lors du démarrage du serveur:', error);
    process.exit(1);
  }
};

// Gestion de l'arrêt propre
process.on('SIGINT', async () => {
  logger.info('Arrêt du serveur...');
  await prisma.$disconnect();
  process.exit(0);
});

process.on('SIGTERM', async () => {
  logger.info('Arrêt du serveur...');
  await prisma.$disconnect();
  process.exit(0);
});

startServer();

export default app;