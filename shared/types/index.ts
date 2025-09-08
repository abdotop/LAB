// Types partagés pour la plateforme E-Parrainages

// Énumérations
export type UserRole = 
  | 'SUPER_ADMIN_DGE'
  | 'ADMIN_DGE'
  | 'GESTIONNAIRE_CANDIDATURE'
  | 'MANDATAIRE'
  | 'CITOYEN';

export type ElectoralCycleStatus = 
  | 'DRAFT'
  | 'ACTIVE'
  | 'PAUSED'
  | 'CLOSED'
  | 'ARCHIVED';

export type CandidateStatus = 
  | 'PENDING'
  | 'APPROVED'
  | 'REJECTED'
  | 'SUSPENDED';

export type SponsorshipStatus = 
  | 'ACTIVE'
  | 'CONTESTED'
  | 'CANCELLED'
  | 'FRAUDULENT';

// Interfaces utilisateur
export interface User {
  id: string;
  nin: string;
  phone: string;
  email?: string;
  firstName: string;
  lastName: string;
  role: UserRole;
  isActive: boolean;
  isVerified: boolean;
  createdAt: string;
  updatedAt: string;
  lastLogin?: string;
}

// Interface pour les cycles électoraux
export interface ElectoralCycle {
  id: string;
  name: string;
  description?: string;
  startDate: string;
  endDate: string;
  sponsorshipStart: string;
  sponsorshipEnd: string;
  status: ElectoralCycleStatus;
  minSponsorsRequired: number;
  createdAt: string;
  updatedAt: string;
}

// Interface pour les candidatures
export interface Candidacy {
  id: string;
  electoralCycleId: string;
  userId: string;
  candidateName: string;
  partyName?: string;
  slogan?: string;
  photoUrl?: string;
  status: CandidateStatus;
  cautionReceipt?: string;
  criminalRecord?: string;
  investitureDeclaration?: string;
  submittedAt: string;
  reviewedAt?: string;
  reviewedBy?: string;
  reviewNotes?: string;
  sponsorshipCount?: number;
}

// Interface pour les parrainages
export interface Sponsorship {
  id: string;
  electoralCycleId: string;
  candidacyId: string;
  citizenId: string;
  status: SponsorshipStatus;
  ipAddress: string;
  userAgent?: string;
  geoLocation?: string;
  contestedAt?: string;
  contestReason?: string;
  createdAt: string;
  updatedAt: string;
}

// Réponses API
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

export interface LoginRequest {
  nin: string;
  pin: string;
}

export interface LoginResponse {
  user: User;
  token: string;
  expiresIn: number;
}

export interface RegisterRequest {
  nin: string;
  phone: string;
  firstName: string;
  lastName: string;
}

export interface OtpVerificationRequest {
  userId: string;
  code: string;
  purpose: string;
}

// Types pour les statistiques
export interface DashboardStats {
  activeCycles: number;
  pendingCandidacies: number;
  totalSponsorships: number;
  fraudReports: number;
}

export interface CandidacyStats {
  candidacyId: string;
  candidateName: string;
  currentSponsorships: number;
  targetSponsorships: number;
  progressPercentage: number;
  dailyGrowth: number;
}

// Types pour les logs d'audit
export interface AuditLog {
  id: string;
  userId: string;
  action: string;
  resource: string;
  resourceId?: string;
  oldData?: any;
  newData?: any;
  ipAddress: string;
  userAgent?: string;
  createdAt: string;
}