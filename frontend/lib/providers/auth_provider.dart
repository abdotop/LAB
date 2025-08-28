import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/user.dart';
import '../services/api_service.dart';

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.read(apiServiceProvider));
});

class AuthState {
  final User? user;
  final String? accessToken;
  final bool isLoading;
  final String? error;

  AuthState({
    this.user,
    this.accessToken,
    this.isLoading = false,
    this.error,
  });

  AuthState copyWith({
    User? user,
    String? accessToken,
    bool? isLoading,
    String? error,
  }) {
    return AuthState(
      user: user ?? this.user,
      accessToken: accessToken ?? this.accessToken,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  final ApiService _apiService;

  AuthNotifier(this._apiService) : super(AuthState()) {
    _loadFromStorage();
  }

  Future<void> _loadFromStorage() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('access_token');
    final userJson = prefs.getString('user');

    if (token != null && userJson != null) {
      try {
        final user = User.fromJson(Map<String, dynamic>.from(
          userJson as Map,
        ));
        
        state = state.copyWith(
          user: user,
          accessToken: token,
        );
        
        _apiService.setAuthToken(token);
      } catch (e) {
        // Clear invalid data
        await _clearStorage();
      }
    }
  }

  Future<void> _saveToStorage(User user, String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('access_token', token);
    await prefs.setString('user', user.toJson().toString());
  }

  Future<void> _clearStorage() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('access_token');
    await prefs.remove('user');
  }

  Future<bool> login(String username, String password) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _apiService.login(username, password);
      
      await _saveToStorage(response['user'], response['accessToken']);
      _apiService.setAuthToken(response['accessToken']);
      
      state = state.copyWith(
        user: response['user'],
        accessToken: response['accessToken'],
        isLoading: false,
      );
      
      return true;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return false;
    }
  }

  Future<bool> register(String username, String password) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _apiService.register(username, password);
      
      await _saveToStorage(response['user'], response['accessToken']);
      _apiService.setAuthToken(response['accessToken']);
      
      state = state.copyWith(
        user: response['user'],
        accessToken: response['accessToken'],
        isLoading: false,
      );
      
      return true;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return false;
    }
  }

  Future<void> logout() async {
    await _clearStorage();
    _apiService.clearAuthToken();
    
    state = AuthState();
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}