import 'package:json_annotation/json_annotation.dart';
import 'user.dart';

part 'game.g.dart';

@JsonSerializable()
class Game {
  final String id;
  final String whiteId;
  final String blackId;
  final User white;
  final User black;
  final String fen;
  final String pgn;
  final String turn; // 'w' or 'b'
  final GameClocks clocks;
  final String? result;
  final String? termination;
  final DateTime createdAt;
  final DateTime updatedAt;

  Game({
    required this.id,
    required this.whiteId,
    required this.blackId,
    required this.white,
    required this.black,
    required this.fen,
    required this.pgn,
    required this.turn,
    required this.clocks,
    this.result,
    this.termination,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Game.fromJson(Map<String, dynamic> json) => _$GameFromJson(json);
  Map<String, dynamic> toJson() => _$GameToJson(this);
}

@JsonSerializable()
class GameClocks {
  final int wMs; // White time in milliseconds
  final int bMs; // Black time in milliseconds

  GameClocks({
    required this.wMs,
    required this.bMs,
  });

  factory GameClocks.fromJson(Map<String, dynamic> json) => _$GameClocksFromJson(json);
  Map<String, dynamic> toJson() => _$GameClocksToJson(this);
}

@JsonSerializable()
class GameSnapshot {
  final String id;
  final String fen;
  final String pgn;
  final String turn;
  final GameClocks clocks;
  final GamePlayers players;
  final List<String>? legalMoves;

  GameSnapshot({
    required this.id,
    required this.fen,
    required this.pgn,
    required this.turn,
    required this.clocks,
    required this.players,
    this.legalMoves,
  });

  factory GameSnapshot.fromJson(Map<String, dynamic> json) => _$GameSnapshotFromJson(json);
  Map<String, dynamic> toJson() => _$GameSnapshotToJson(this);
}

@JsonSerializable()
class GamePlayers {
  final User white;
  final User black;

  GamePlayers({
    required this.white,
    required this.black,
  });

  factory GamePlayers.fromJson(Map<String, dynamic> json) => _$GamePlayersFromJson(json);
  Map<String, dynamic> toJson() => _$GamePlayersToJson(this);
}