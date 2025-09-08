import { Router, Request, Response } from 'express';
import { asyncHandler } from '../middleware/errorHandler';

const router = Router();

// POST /api/auth/register - Inscription citoyen avec NIN
router.post('/register', asyncHandler(async (req: Request, res: Response) => {
  // TODO: Implémenter l'inscription avec NIN + OTP
  res.json({ message: 'Inscription en cours de développement' });
}));

// POST /api/auth/login - Connexion avec NIN + PIN
router.post('/login', asyncHandler(async (req: Request, res: Response) => {
  // TODO: Implémenter la connexion
  res.json({ message: 'Connexion en cours de développement' });
}));

// POST /api/auth/verify-otp - Vérification OTP
router.post('/verify-otp', asyncHandler(async (req: Request, res: Response) => {
  // TODO: Implémenter la vérification OTP
  res.json({ message: 'Vérification OTP en cours de développement' });
}));

// POST /api/auth/logout - Déconnexion
router.post('/logout', asyncHandler(async (req: Request, res: Response) => {
  // TODO: Implémenter la déconnexion
  res.json({ message: 'Déconnexion réussie' });
}));

export default router;