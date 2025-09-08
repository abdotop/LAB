import { Router, Request, Response } from 'express';
import { asyncHandler } from '../middleware/errorHandler';

const router = Router();

// GET /api/candidacies - Liste des candidatures
router.get('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Liste des candidatures en cours de développement' });
}));

// POST /api/candidacies - Déposer une candidature
router.post('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Dépôt de candidature en cours de développement' });
}));

// PUT /api/candidacies/:id/approve - Approuver candidature (Admin DGE)
router.put('/:id/approve', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Approbation candidature en cours de développement' });
}));

export default router;