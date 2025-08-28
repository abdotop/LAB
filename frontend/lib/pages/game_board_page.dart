import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/game_provider.dart';
import '../providers/auth_provider.dart';
import '../widgets/chess_board.dart';

class GameBoardPage extends ConsumerWidget {
  final String gameId;

  const GameBoardPage({
    super.key,
    required this.gameId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final gameState = ref.watch(gameProvider(gameId));
    final authState = ref.watch(authProvider);
    final currentUserId = authState.user?.id;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Chess Game'),
        backgroundColor: Colors.brown,
        foregroundColor: Colors.white,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.go('/'),
        ),
        actions: [
          // Connection status indicator
          Icon(
            gameState.isConnected ? Icons.wifi : Icons.wifi_off,
            color: gameState.isConnected ? Colors.green : Colors.red,
          ),
          const SizedBox(width: 8),
          // Game menu
          PopupMenuButton(
            itemBuilder: (context) => [
              PopupMenuItem(
                child: const Text('Resign'),
                onTap: () => _showResignDialog(context, ref),
              ),
              PopupMenuItem(
                child: const Text('Offer Draw'),
                onTap: () => ref.read(gameProvider(gameId).notifier).offerDraw(),
              ),
            ],
          ),
        ],
      ),
      body: gameState.isLoading
          ? const Center(child: CircularProgressIndicator())
          : gameState.error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'Error: ${gameState.error}',
                        style: const TextStyle(color: Colors.red),
                        textAlign: TextAlign.center,
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: () => context.go('/'),
                        child: const Text('Back to Menu'),
                      ),
                    ],
                  ),
                )
              : gameState.gameSnapshot != null
                  ? _buildGameView(context, ref, gameState, currentUserId)
                  : const Center(child: Text('Loading game...')),
    );
  }

  Widget _buildGameView(
    BuildContext context,
    WidgetRef ref,
    GameState gameState,
    String? currentUserId,
  ) {
    final game = gameState.gameSnapshot!;
    final isWhite = game.players.white.id == currentUserId;
    final isYourTurn = (game.turn == 'w' && isWhite) || (game.turn == 'b' && !isWhite);

    return Column(
      children: [
        // Opponent info
        _buildPlayerInfo(
          user: isWhite ? game.players.black : game.players.white,
          timeMs: isWhite ? game.clocks.bMs : game.clocks.wMs,
          isActive: !isYourTurn && game.turn != '',
        ),
        
        const SizedBox(height: 8),

        // Chess board
        Expanded(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Center(
              child: ChessBoard(
                fen: game.fen,
                legalMoves: game.legalMoves ?? [],
                yourColor: gameState.yourColor,
                isYourTurn: isYourTurn,
                onMove: (from, to, {promotion}) {
                  ref.read(gameProvider(gameId).notifier).makeMove(
                        from,
                        to,
                        promotion: promotion,
                      );
                },
                // TODO: Track last move for highlighting
              ),
            ),
          ),
        ),

        const SizedBox(height: 8),

        // Current player info
        _buildPlayerInfo(
          user: isWhite ? game.players.white : game.players.black,
          timeMs: isWhite ? game.clocks.wMs : game.clocks.bMs,
          isActive: isYourTurn && game.turn != '',
        ),

        // Game status
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(16),
          color: isYourTurn ? Colors.green[100] : Colors.grey[100],
          child: Text(
            _getGameStatusText(game, isYourTurn, gameState.yourColor),
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: isYourTurn ? Colors.green[800] : Colors.grey[700],
            ),
            textAlign: TextAlign.center,
          ),
        ),
      ],
    );
  }

  Widget _buildPlayerInfo({
    required user,
    required int timeMs,
    required bool isActive,
  }) {
    final minutes = timeMs ~/ 60000;
    final seconds = (timeMs % 60000) ~/ 1000;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: isActive ? Colors.blue[50] : null,
      child: Row(
        children: [
          CircleAvatar(
            backgroundColor: Colors.brown,
            foregroundColor: Colors.white,
            child: Text(
              user.username[0].toUpperCase(),
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  user.username,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Text(
                  'Rating: ${user.rating}',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey[600],
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              color: isActive ? Colors.red : Colors.grey[300],
              borderRadius: BorderRadius.circular(6),
            ),
            child: Text(
              '${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: isActive ? Colors.white : Colors.black87,
              ),
            ),
          ),
        ],
      ),
    );
  }

  String _getGameStatusText(game, bool isYourTurn, String? yourColor) {
    // Check if game is over
    if (game.result != null) {
      switch (game.result) {
        case '1-0':
          return 'White wins! ${_getTerminationText(game.termination)}';
        case '0-1':
          return 'Black wins! ${_getTerminationText(game.termination)}';
        case '1/2-1/2':
          return 'Draw! ${_getTerminationText(game.termination)}';
      }
    }

    if (isYourTurn) {
      return 'Your turn';
    } else {
      final opponent = yourColor == 'WHITE' ? 'Black' : 'White';
      return '$opponent\'s turn';
    }
  }

  String _getTerminationText(String? termination) {
    switch (termination) {
      case 'CHECKMATE':
        return 'Checkmate';
      case 'STALEMATE':
        return 'Stalemate';
      case 'RESIGN':
        return 'Resignation';
      case 'AGREEMENT':
        return 'Draw by agreement';
      case 'THREEFOLD':
        return 'Threefold repetition';
      case 'FIFTY_MOVE':
        return 'Fifty-move rule';
      case 'INSUFFICIENT':
        return 'Insufficient material';
      case 'TIME':
        return 'Time out';
      default:
        return '';
    }
  }

  void _showResignDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Resign Game'),
        content: const Text('Are you sure you want to resign? This will end the game immediately.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(context).pop();
              ref.read(gameProvider(gameId).notifier).resign();
            },
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Resign', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );
  }
}