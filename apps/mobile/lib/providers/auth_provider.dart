import 'package:flutter/material.dart';

class AuthProvider with ChangeNotifier {
  bool _isAuthenticated = false;
  String? _userId;
  String? _nin;
  String? _phone;
  bool _isLoading = false;

  // Getters
  bool get isAuthenticated => _isAuthenticated;
  String? get userId => _userId;
  String? get nin => _nin;
  String? get phone => _phone;
  bool get isLoading => _isLoading;

  // Inscription avec NIN
  Future<bool> registerWithNin(String nin, String phone, String firstName, String lastName) async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter l'appel API pour l'inscription
      await Future.delayed(const Duration(seconds: 2)); // Simulation
      
      // Simulation d'envoi OTP
      _nin = nin;
      _phone = phone;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  // Vérification OTP
  Future<bool> verifyOtp(String otp) async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter la vérification OTP
      await Future.delayed(const Duration(seconds: 1)); // Simulation
      
      if (otp == '1234') { // Code de test
        _isAuthenticated = true;
        _userId = 'test_user_id';
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

  // Connexion avec NIN + PIN
  Future<bool> loginWithNinAndPin(String nin, String pin) async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter l'appel API pour la connexion
      await Future.delayed(const Duration(seconds: 2)); // Simulation
      
      if (nin.isNotEmpty && pin == '1234') { // Test
        _isAuthenticated = true;
        _userId = 'test_user_id';
        _nin = nin;
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

  // Déconnexion
  Future<void> logout() async {
    _isAuthenticated = false;
    _userId = null;
    _nin = null;
    _phone = null;
    notifyListeners();
  }

  // Définir un PIN
  Future<bool> setPin(String pin) async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implémenter la sauvegarde du PIN (hashé)
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