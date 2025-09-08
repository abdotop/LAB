import { Router, Request, Response } from 'express';
import { asyncHandler } from '../middleware/errorHandler';

const router = Router();

// GET /api/users/profile - Profil utilisateur
router.get('/profile', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Profil utilisateur en cours de développement' });
}));

// PUT /api/users/profile - Mise à jour du profil
router.put('/profile', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Mise à jour profil en cours de développement' });
}));

export default router;