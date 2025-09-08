import 'package:flutter/material.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Accueil E-Parrainages')),
      body: const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.how_to_vote, size: 100, color: Color(0xFF00853F)),
            SizedBox(height: 20),
            Text('Bienvenue dans E-Parrainages', style: TextStyle(fontSize: 24)),
            SizedBox(height: 10),
            Text('Application en cours de développement', style: TextStyle(fontSize: 16)),
          ],
        ),
      ),
    );
  }
}