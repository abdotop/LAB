import { Router, Request, Response } from 'express';
import { asyncHandler } from '../middleware/errorHandler';

const router = Router();

// GET /api/sponsorships - Mes parrainages
router.get('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Mes parrainages en cours de développement' });
}));

// POST /api/sponsorships - Parrainer un candidat
router.post('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Parrainage en cours de développement' });
}));

// POST /api/sponsorships/:id/contest - Contester un parrainage
router.post('/:id/contest', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Contestation en cours de développement' });
}));

export default router;