import { Router, Request, Response } from 'express';
import { asyncHandler } from '../middleware/errorHandler';

const router = Router();

// GET /api/electoral-cycles - Liste des cycles électoraux
router.get('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Liste des cycles électoraux en cours de développement' });
}));

// POST /api/electoral-cycles - Créer un cycle électoral (Admin DGE)
router.post('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Création cycle électoral en cours de développement' });
}));

// GET /api/electoral-cycles/:id - Détails d'un cycle
router.get('/:id', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Détails cycle électoral en cours de développement' });
}));

export default router;