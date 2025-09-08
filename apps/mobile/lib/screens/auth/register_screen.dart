import 'package:flutter/material.dart';

class RegisterScreen extends StatelessWidget {
  const RegisterScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Inscription')),
      body: const Center(
        child: Text('Écran d\'inscription - En cours de développement'),
      ),
    );
  }
}

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Accueil')),
      body: const Center(
        child: Text('Écran d\'accueil - En cours de développement'),
      ),
    );
  }
}

class CandidatesScreen extends StatelessWidget {
  const CandidatesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Candidats')),
      body: const Center(
        child: Text('Liste des candidats - En cours de développement'),
      ),
    );
  }
}

class SponsorshipScreen extends StatelessWidget {
  final String candidateId;
  
  const SponsorshipScreen({super.key, required this.candidateId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Parrainage')),
      body: Center(
        child: Text('Parrainage pour candidat $candidateId - En cours de développement'),
      ),
    );
  }
}