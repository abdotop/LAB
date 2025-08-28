import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http/http.dart' as http;
import '../models/user.dart';
import '../models/lobby.dart';
import '../models/game.dart';
import '../models/move.dart';

final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService();
});

class ApiService {
  static const String baseUrl = 'http://localhost:8080/api';
  String? _authToken;

  void setAuthToken(String token) {
    _authToken = token;
  }

  void clearAuthToken() {
    _authToken = null;
  }

  Map<String, String> get _headers {
    final headers = <String, String>{
      'Content-Type': 'application/json',
    };
    
    if (_authToken != null) {
      headers['Authorization'] = 'Bearer $_authToken';
    }
    
    return headers;
  }

  Future<Map<String, dynamic>> _handleResponse(http.Response response) async {
    final body = json.decode(response.body);
    
    if (response.statusCode >= 200 && response.statusCode < 300) {
      return body;
    } else {
      throw Exception(body['error'] ?? 'Request failed');
    }
  }

  // Auth endpoints
  Future<Map<String, dynamic>> login(String username, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/login'),
      headers: _headers,
      body: json.encode({
        'username': username,
        'password': password,
      }),
    );

    final data = await _handleResponse(response);
    return {
      'user': User.fromJson(data['user']),
      'accessToken': data['accessToken'],
    };
  }

  Future<Map<String, dynamic>> register(String username, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/register'),
      headers: _headers,
      body: json.encode({
        'username': username,
        'password': password,
      }),
    );

    final data = await _handleResponse(response);
    return {
      'user': User.fromJson(data['user']),
      'accessToken': data['accessToken'],
    };
  }

  Future<User> getMe() async {
    final response = await http.get(
      Uri.parse('$baseUrl/users/me'),
      headers: _headers,
    );

    final data = await _handleResponse(response);
    return User.fromJson(data['user']);
  }

  Future<User> updateMe({String? avatarUrl, String? fcmToken}) async {
    final body = <String, dynamic>{};
    if (avatarUrl != null) body['avatarUrl'] = avatarUrl;
    if (fcmToken != null) body['fcmToken'] = fcmToken;

    final response = await http.patch(
      Uri.parse('$baseUrl/users/me'),
      headers: _headers,
      body: json.encode(body),
    );

    final data = await _handleResponse(response);
    return User.fromJson(data['user']);
  }

  // Lobby endpoints
  Future<Lobby> createLobby({
    required String type,
    required TimeControl timeControl,
    bool inviteOnly = false,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/lobbies'),
      headers: _headers,
      body: json.encode({
        'type': type,
        'timeControl': timeControl.toJson(),
        'inviteOnly': inviteOnly,
      }),
    );

    final data = await _handleResponse(response);
    return Lobby.fromJson(data['lobby']);
  }

  Future<List<Lobby>> getPublicLobbies() async {
    final response = await http.get(
      Uri.parse('$baseUrl/lobbies/public'),
      headers: _headers,
    );

    final data = await _handleResponse(response);
    return (data['lobbies'] as List)
        .map((lobby) => Lobby.fromJson(lobby))
        .toList();
  }

  Future<Map<String, dynamic>> joinLobby(String lobbyId) async {
    final response = await http.post(
      Uri.parse('$baseUrl/lobbies/$lobbyId/join'),
      headers: _headers,
    );

    return await _handleResponse(response);
  }

  Future<void> deleteLobby(String lobbyId) async {
    final response = await http.delete(
      Uri.parse('$baseUrl/lobbies/$lobbyId'),
      headers: _headers,
    );

    await _handleResponse(response);
  }

  // Game endpoints
  Future<Game> createGame({
    required String whiteId,
    required String blackId,
    required TimeControl timeControl,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/games'),
      headers: _headers,
      body: json.encode({
        'whiteId': whiteId,
        'blackId': blackId,
        'timeControl': timeControl.toJson(),
      }),
    );

    final data = await _handleResponse(response);
    return Game.fromJson(data['game']);
  }

  Future<Map<String, dynamic>> getGame(String gameId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/games/$gameId'),
      headers: _headers,
    );

    return await _handleResponse(response);
  }

  Future<List<Move>> getGameMoves(String gameId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/games/$gameId/moves'),
      headers: _headers,
    );

    final data = await _handleResponse(response);
    return (data as List)
        .map((move) => Move.fromJson(move))
        .toList();
  }

  Future<void> offerDraw(String gameId) async {
    final response = await http.post(
      Uri.parse('$baseUrl/games/$gameId/draw-offer'),
      headers: _headers,
    );

    await _handleResponse(response);
  }

  Future<void> respondToDraw(String gameId, bool accept) async {
    final response = await http.post(
      Uri.parse('$baseUrl/games/$gameId/draw-respond'),
      headers: _headers,
      body: json.encode({'accept': accept}),
    );

    await _handleResponse(response);
  }

  Future<void> resign(String gameId) async {
    final response = await http.post(
      Uri.parse('$baseUrl/games/$gameId/resign'),
      headers: _headers,
    );

    await _handleResponse(response);
  }
}