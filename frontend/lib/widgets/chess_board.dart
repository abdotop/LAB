import 'package:flutter/material.dart';

// Chess piece unicode symbols
class ChessPieces {
  static const Map<String, String> pieces = {
    'wk': '♔', // White King
    'wq': '♕', // White Queen
    'wr': '♖', // White Rook
    'wb': '♗', // White Bishop
    'wn': '♘', // White Knight
    'wp': '♙', // White Pawn
    'bk': '♚', // Black King
    'bq': '♛', // Black Queen
    'br': '♜', // Black Rook
    'bb': '♝', // Black Bishop
    'bn': '♞', // Black Knight
    'bp': '♟', // Black Pawn
  };

  static String? getPiece(String? type, String? color) {
    if (type == null || color == null) return null;
    return pieces['$color$type'];
  }
}

class ChessBoard extends StatefulWidget {
  final String fen;
  final List<String> legalMoves;
  final String? yourColor;
  final bool isYourTurn;
  final Function(String from, String to, {String? promotion}) onMove;
  final String? lastMoveFrom;
  final String? lastMoveTo;

  const ChessBoard({
    super.key,
    required this.fen,
    required this.legalMoves,
    this.yourColor,
    required this.isYourTurn,
    required this.onMove,
    this.lastMoveFrom,
    this.lastMoveTo,
  });

  @override
  State<ChessBoard> createState() => _ChessBoardState();
}

class _ChessBoardState extends State<ChessBoard> {
  String? _selectedSquare;
  List<String> _highlightedSquares = [];

  @override
  void didUpdateWidget(ChessBoard oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.fen != oldWidget.fen) {
      // Clear selection when position changes
      _selectedSquare = null;
      _highlightedSquares = [];
    }
  }

  List<List<ChessPiece?>> _parseFEN(String fen) {
    final board = List.generate(8, (i) => List<ChessPiece?>.filled(8, null));
    final boardPart = fen.split(' ')[0];
    final rows = boardPart.split('/');

    for (int row = 0; row < 8; row++) {
      int col = 0;
      for (final char in rows[row].split('')) {
        if (char.contains(RegExp(r'[1-8]'))) {
          // Empty squares
          col += int.parse(char);
        } else {
          // Piece
          final color = char == char.toUpperCase() ? 'w' : 'b';
          final type = char.toLowerCase();
          board[row][col] = ChessPiece(type: type, color: color);
          col++;
        }
      }
    }

    return board;
  }

  String _squareToString(int row, int col) {
    final file = String.fromCharCode('a'.codeUnitAt(0) + col);
    final rank = (8 - row).toString();
    return '$file$rank';
  }

  (int, int) _stringToSquare(String square) {
    final file = square.codeUnitAt(0) - 'a'.codeUnitAt(0);
    final rank = int.parse(square[1]);
    return (8 - rank, file);
  }

  void _onSquareTap(int row, int col) {
    if (!widget.isYourTurn) return;

    final square = _squareToString(row, col);
    final board = _parseFEN(widget.fen);
    final piece = board[row][col];

    if (_selectedSquare == null) {
      // First tap - select piece
      if (piece != null && piece.color == widget.yourColor?[0]) {
        setState(() {
          _selectedSquare = square;
          _highlightedSquares = _getLegalMovesForSquare(square);
        });
      }
    } else {
      // Second tap - move or change selection
      if (_selectedSquare == square) {
        // Deselect
        setState(() {
          _selectedSquare = null;
          _highlightedSquares = [];
        });
      } else if (piece != null && piece.color == widget.yourColor?[0]) {
        // Select different piece
        setState(() {
          _selectedSquare = square;
          _highlightedSquares = _getLegalMovesForSquare(square);
        });
      } else if (_highlightedSquares.contains(square)) {
        // Make move
        _makeMove(_selectedSquare!, square);
        setState(() {
          _selectedSquare = null;
          _highlightedSquares = [];
        });
      }
    }
  }

  List<String> _getLegalMovesForSquare(String fromSquare) {
    return widget.legalMoves
        .where((move) => move.startsWith(fromSquare))
        .map((move) => move.substring(2, 4))
        .toList();
  }

  void _makeMove(String from, String to) {
    final board = _parseFEN(widget.fen);
    final (toRow, toCol) = _stringToSquare(to);
    final (fromRow, fromCol) = _stringToSquare(from);
    final piece = board[fromRow][fromCol];

    // Check for pawn promotion
    if (piece?.type == 'p' && (toRow == 0 || toRow == 7)) {
      _showPromotionDialog(from, to);
    } else {
      widget.onMove(from, to);
    }
  }

  void _showPromotionDialog(String from, String to) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        title: const Text('Promote Pawn'),
        content: Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _PromotionPiece(
              piece: ChessPieces.getPiece('q', widget.yourColor?[0])!,
              onTap: () {
                Navigator.of(context).pop();
                widget.onMove(from, to, promotion: 'q');
              },
            ),
            _PromotionPiece(
              piece: ChessPieces.getPiece('r', widget.yourColor?[0])!,
              onTap: () {
                Navigator.of(context).pop();
                widget.onMove(from, to, promotion: 'r');
              },
            ),
            _PromotionPiece(
              piece: ChessPieces.getPiece('b', widget.yourColor?[0])!,
              onTap: () {
                Navigator.of(context).pop();
                widget.onMove(from, to, promotion: 'b');
              },
            ),
            _PromotionPiece(
              piece: ChessPieces.getPiece('n', widget.yourColor?[0])!,
              onTap: () {
                Navigator.of(context).pop();
                widget.onMove(from, to, promotion: 'n');
              },
            ),
          ],
        ),
      ),
    );
  }

  Color _getSquareColor(int row, int col) {
    final square = _squareToString(row, col);
    
    // Last move highlighting
    if (square == widget.lastMoveFrom || square == widget.lastMoveTo) {
      return Colors.yellow.withOpacity(0.6);
    }
    
    // Selected square
    if (square == _selectedSquare) {
      return Colors.blue.withOpacity(0.6);
    }
    
    // Legal move highlighting
    if (_highlightedSquares.contains(square)) {
      return Colors.green.withOpacity(0.4);
    }
    
    // Default square colors
    final isLight = (row + col) % 2 == 0;
    return isLight ? Colors.grey[200]! : Colors.brown[400]!;
  }

  @override
  Widget build(BuildContext context) {
    final board = _parseFEN(widget.fen);
    final isFlipped = widget.yourColor == 'BLACK';

    return AspectRatio(
      aspectRatio: 1.0,
      child: Container(
        decoration: BoxDecoration(
          border: Border.all(color: Colors.brown, width: 2),
        ),
        child: Column(
          children: List.generate(8, (row) {
            final displayRow = isFlipped ? 7 - row : row;
            return Expanded(
              child: Row(
                children: List.generate(8, (col) {
                  final displayCol = isFlipped ? 7 - col : col;
                  final piece = board[displayRow][displayCol];
                  
                  return Expanded(
                    child: GestureDetector(
                      onTap: () => _onSquareTap(displayRow, displayCol),
                      child: Container(
                        color: _getSquareColor(displayRow, displayCol),
                        child: Center(
                          child: piece != null
                              ? Text(
                                  ChessPieces.getPiece(piece.type, piece.color) ?? '',
                                  style: const TextStyle(
                                    fontSize: 36,
                                    fontWeight: FontWeight.w400,
                                  ),
                                )
                              : null,
                        ),
                      ),
                    ),
                  );
                }),
              ),
            );
          }),
        ),
      ),
    );
  }
}

class ChessPiece {
  final String type; // 'k', 'q', 'r', 'b', 'n', 'p'
  final String color; // 'w', 'b'

  ChessPiece({required this.type, required this.color});
}

class _PromotionPiece extends StatelessWidget {
  final String piece;
  final VoidCallback onTap;

  const _PromotionPiece({
    required this.piece,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: 60,
        height: 60,
        decoration: BoxDecoration(
          color: Colors.grey[200],
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: Colors.brown),
        ),
        child: Center(
          child: Text(
            piece,
            style: const TextStyle(fontSize: 36),
          ),
        ),
      ),
    );
  }
}