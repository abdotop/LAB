import 'package:flutter/material.dart';

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