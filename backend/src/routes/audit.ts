import { Router, Request, Response } from 'express';
import { asyncHandler } from '../middleware/errorHandler';

const router = Router();

// GET /api/audit - Piste d'audit (Super Admin DGE uniquement)
router.get('/', asyncHandler(async (req: Request, res: Response) => {
  res.json({ message: 'Piste d\'audit en cours de développement' });
}));

export default router;