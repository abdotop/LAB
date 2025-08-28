import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../models/game.dart';
import '../services/api_service.dart';
import '../services/websocket_service.dart';
import 'auth_provider.dart';

final gameProvider = StateNotifierProvider.family<GameNotifier, GameState, String>((ref, gameId) {
  final auth = ref.watch(authProvider);
  final apiService = ref.read(apiServiceProvider);
  final wsService = ref.read(websocketServiceProvider);
  
  return GameNotifier(gameId, auth.accessToken, apiService, wsService);
});

class GameState {
  final GameSnapshot? gameSnapshot;
  final List<String> legalMoves;
  final bool isLoading;
  final String? error;
  final bool isConnected;
  final String? yourColor; // 'WHITE' or 'BLACK'

  GameState({
    this.gameSnapshot,
    this.legalMoves = const [],
    this.isLoading = false,
    this.error,
    this.isConnected = false,
    this.yourColor,
  });

  GameState copyWith({
    GameSnapshot? gameSnapshot,
    List<String>? legalMoves,
    bool? isLoading,
    String? error,
    bool? isConnected,
    String? yourColor,
  }) {
    return GameState(
      gameSnapshot: gameSnapshot ?? this.gameSnapshot,
      legalMoves: legalMoves ?? this.legalMoves,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      isConnected: isConnected ?? this.isConnected,
      yourColor: yourColor ?? this.yourColor,
    );
  }
}

class GameNotifier extends StateNotifier<GameState> {
  final String gameId;
  final String? authToken;
  final ApiService _apiService;
  final WebSocketService _wsService;

  GameNotifier(
    this.gameId,
    this.authToken,
    this._apiService,
    this._wsService,
  ) : super(GameState(isLoading: true)) {
    _initialize();
  }

  Future<void> _initialize() async {
    if (authToken == null) {
      state = state.copyWith(
        isLoading: false,
        error: 'Not authenticated',
      );
      return;
    }

    try {
      // Load game data
      final gameData = await _apiService.getGame(gameId);
      final gameSnapshot = GameSnapshot.fromJson(gameData);
      
      state = state.copyWith(
        gameSnapshot: gameSnapshot,
        legalMoves: gameSnapshot.legalMoves ?? [],
        isLoading: false,
      );

      // Connect to WebSocket
      await _connectWebSocket();
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  Future<void> _connectWebSocket() async {
    try {
      _wsService.connect(authToken!);
      _wsService.joinGame(gameId);
      
      // Listen for WebSocket messages
      _wsService.messages.listen((message) {
        _handleWebSocketMessage(message);
      });
      
      state = state.copyWith(isConnected: true);
    } catch (e) {
      state = state.copyWith(error: 'WebSocket connection failed: $e');
    }
  }

  void _handleWebSocketMessage(Map<String, dynamic> message) {
    final type = message['type'];
    final data = message['data'];

    switch (type) {
      case 'joined':
        state = state.copyWith(yourColor: data['youAre']);
        break;
      case 'state':
        final gameSnapshot = GameSnapshot.fromJson(data['game']);
        state = state.copyWith(
          gameSnapshot: gameSnapshot,
          legalMoves: gameSnapshot.legalMoves ?? [],
        );
        break;
      case 'moveAccepted':
        // Update game state after move
        _refreshGameState();
        break;
      case 'illegalMove':
        state = state.copyWith(error: 'Illegal move: ${data['reason']}');
        break;
      case 'gameOver':
        // Handle game over
        _refreshGameState();
        break;
      case 'error':
        state = state.copyWith(error: data['message']);
        break;
    }
  }

  Future<void> _refreshGameState() async {
    try {
      final gameData = await _apiService.getGame(gameId);
      final gameSnapshot = GameSnapshot.fromJson(gameData);
      
      state = state.copyWith(
        gameSnapshot: gameSnapshot,
        legalMoves: gameSnapshot.legalMoves ?? [],
      );
    } catch (e) {
      // Handle error silently for refresh
    }
  }

  void makeMove(String from, String to, {String? promotion}) {
    _wsService.makeMove(gameId, from, to, promotion: promotion);
  }

  void resign() {
    _wsService.resign(gameId);
  }

  void offerDraw() {
    _wsService.offerDraw(gameId);
  }

  void clearError() {
    state = state.copyWith(error: null);
  }

  @override
  void dispose() {
    _wsService.disconnect();
    super.dispose();
  }
}