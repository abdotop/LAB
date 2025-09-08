import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../providers/auth_provider.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _ninController = TextEditingController();
  final _pinController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.background,
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Logo et titre
                const Icon(
                  Icons.how_to_vote,
                  size: 80,
                  color: Color(0xFF00853F),
                ),
                const SizedBox(height: 24),
                
                Text(
                  'Connexion',
                  style: Theme.of(context).textTheme.displayMedium,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 8),
                
                Text(
                  'Utilisez votre NIN et votre code PIN',
                  style: Theme.of(context).textTheme.bodyMedium,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 48),
                
                // Champ NIN
                TextFormField(
                  controller: _ninController,
                  decoration: const InputDecoration(
                    labelText: 'Numéro d\'Identification Nationale (NIN)',
                    prefixIcon: Icon(Icons.person),
                  ),
                  keyboardType: TextInputType.number,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Veuillez saisir votre NIN';
                    }
                    if (value.length != 13) {
                      return 'Le NIN doit contenir 13 chiffres';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
                
                // Champ PIN
                TextFormField(
                  controller: _pinController,
                  decoration: const InputDecoration(
                    labelText: 'Code PIN (4 chiffres)',
                    prefixIcon: Icon(Icons.lock),
                  ),
                  keyboardType: TextInputType.number,
                  obscureText: true,
                  maxLength: 4,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Veuillez saisir votre code PIN';
                    }
                    if (value.length != 4) {
                      return 'Le code PIN doit contenir 4 chiffres';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 32),
                
                // Bouton de connexion
                Consumer<AuthProvider>(
                  builder: (context, authProvider, child) {
                    return ElevatedButton(
                      onPressed: authProvider.isLoading ? null : _handleLogin,
                      child: authProvider.isLoading
                          ? const CircularProgressIndicator(color: Colors.white)
                          : const Text('Se connecter'),
                    );
                  },
                ),
                const SizedBox(height: 16),
                
                // Lien vers l'inscription
                TextButton(
                  onPressed: () => context.go('/register'),
                  child: const Text('Première connexion ? Créer un compte'),
                ),
                
                const SizedBox(height: 24),
                
                // Informations légales
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.blue.shade50,
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.blue.shade200),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Icon(Icons.info_outline, color: Colors.blue.shade600, size: 20),
                          const SizedBox(width: 8),
                          Text(
                            'Important',
                            style: TextStyle(
                              fontWeight: FontWeight.w600,
                              color: Colors.blue.shade700,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Cette application est destinée aux citoyens sénégalais pour le parrainage électoral officiel. Votre NIN sera vérifié avec la base de données électorale.',
                        style: TextStyle(
                          fontSize: 12,
                          color: Colors.blue.shade700,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _handleLogin() async {
    if (_formKey.currentState!.validate()) {
      final authProvider = Provider.of<AuthProvider>(context, listen: false);
      
      final success = await authProvider.loginWithNinAndPin(
        _ninController.text,
        _pinController.text,
      );
      
      if (success && mounted) {
        context.go('/home');
      } else if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('NIN ou code PIN incorrect'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  @override
  void dispose() {
    _ninController.dispose();
    _pinController.dispose();
    super.dispose();
  }
}