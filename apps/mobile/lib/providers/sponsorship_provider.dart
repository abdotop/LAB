import 'package:flutter/material.dart';

class Candidate {
  final String id;
  final String name;
  final String? party;
  final String? slogan;
  final String? photoUrl;
  final int sponsorshipCount;
  final int minSponsors;

  Candidate({
    required this.id,
    required this.name,
    this.party,
    this.slogan,
    this.photoUrl,
    required this.sponsorshipCount,
    required this.minSponsors,
  });

  double get progressPercentage => (sponsorshipCount / minSponsors * 100).clamp(0, 100);
}

class SponsorshipProvider with ChangeNotifier {
  List<Candidate> _candidates = [];
  String? _mySponsorshipId;
  bool _isLoading = false;
  bool _hasSponsored = false;

  // Getters
  List<Candidate> get candidates => _candidates;
  String? get mySponsorshipId => _mySponsorshipId;
  bool get isLoading => _isLoading;
  bool get hasSponsored => _hasSponsored;

  // Charger la liste des candidats
  Future<void> loadCandidates() async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter l'appel API pour récupérer les candidats
      await Future.delayed(const Duration(seconds: 1)); // Simulation
      
      // Données de test
      _candidates = [
        Candidate(
          id: '1',
          name: 'Macky Sall',
          party: 'Benno Bokk Yakaar',
          slogan: 'L\'Emergence en Marche',
          sponsorshipCount: 12847,
          minSponsors: 15000,
        ),
        Candidate(
          id: '2',
          name: 'Ousmane Sonko',
          party: 'PASTEF',
          slogan: 'Sénégal 2050',
          sponsorshipCount: 8956,
          minSponsors: 15000,
        ),
        Candidate(
          id: '3',
          name: 'Khalifa Sall',
          party: 'Taxawu Senegaal',
          slogan: 'Pour un Sénégal Uni',
          sponsorshipCount: 6234,
          minSponsors: 15000,
        ),
      ];
      
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Parrainer un candidat
  Future<bool> sponsorCandidate(String candidateId, String pin) async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter l'appel API pour le parrainage
      await Future.delayed(const Duration(seconds: 2)); // Simulation
      
      if (pin == '1234') { // Code de test
        _mySponsorshipId = 'sponsorship_id_${DateTime.now().millisecondsSinceEpoch}';
        _hasSponsored = true;
        
        // Mettre à jour le nombre de parrainages du candidat
        final candidateIndex = _candidates.indexWhere((c) => c.id == candidateId);
        if (candidateIndex != -1) {
          _candidates[candidateIndex] = Candidate(
            id: _candidates[candidateIndex].id,
            name: _candidates[candidateIndex].name,
            party: _candidates[candidateIndex].party,
            slogan: _candidates[candidateIndex].slogan,
            photoUrl: _candidates[candidateIndex].photoUrl,
            sponsorshipCount: _candidates[candidateIndex].sponsorshipCount + 1,
            minSponsors: _candidates[candidateIndex].minSponsors,
          );
        }
        
        _isLoading = false;
        notifyListeners();
        return true;
      }
      
      _isLoading = false;
      notifyListeners();
      return false;
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  // Vérifier si l'utilisateur a déjà parrainé
  Future<void> checkExistingSponsorship() async {
    try {
      // TODO: Implémenter la vérification du parrainage existant
      await Future.delayed(const Duration(milliseconds: 500)); // Simulation
      
      // Simulation: l'utilisateur n'a pas encore parrainé
      _hasSponsored = false;
      _mySponsorshipId = null;
      notifyListeners();
    } catch (e) {
      // Gérer l'erreur
    }
  }

  // Contester un parrainage (dans les 24h)
  Future<bool> contestSponsorship(String reason) async {
    if (_mySponsorshipId == null) return false;
    
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter la contestation
      await Future.delayed(const Duration(seconds: 1)); // Simulation
      
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }
}