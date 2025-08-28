import 'package:json_annotation/json_annotation.dart';
import 'user.dart';
import 'game.dart';

part 'move.g.dart';

@JsonSerializable()
class Move {
  final String id;
  final String gameId;
  final String moverId;
  final User mover;
  final String san; // Standard Algebraic Notation
  final String from;
  final String to;
  final String? promotion;
  final DateTime createdAt;

  Move({
    required this.id,
    required this.gameId,
    required this.moverId,
    required this.mover,
    required this.san,
    required this.from,
    required this.to,
    this.promotion,
    required this.createdAt,
  });

  factory Move.fromJson(Map<String, dynamic> json) => _$MoveFromJson(json);
  Map<String, dynamic> toJson() => _$MoveToJson(this);
}