import 'package:json_annotation/json_annotation.dart';
import 'user.dart';

part 'lobby.g.dart';

@JsonSerializable()
class Lobby {
  final String id;
  final String type; // 'PUBLIC' or 'PRIVATE'
  final String creatorId;
  final User creator;
  final String? inviteCode;
  final TimeControl timeControl;
  final String status; // 'OPEN', 'MATCHED', 'CLOSED'
  final DateTime createdAt;

  Lobby({
    required this.id,
    required this.type,
    required this.creatorId,
    required this.creator,
    this.inviteCode,
    required this.timeControl,
    required this.status,
    required this.createdAt,
  });

  factory Lobby.fromJson(Map<String, dynamic> json) => _$LobbyFromJson(json);
  Map<String, dynamic> toJson() => _$LobbyToJson(this);
}

@JsonSerializable()
class TimeControl {
  final int initial; // seconds
  final int increment; // seconds

  TimeControl({
    required this.initial,
    required this.increment,
  });

  factory TimeControl.fromJson(Map<String, dynamic> json) => _$TimeControlFromJson(json);
  Map<String, dynamic> toJson() => _$TimeControlToJson(this);

  String get displayText {
    if (increment == 0) {
      return '${initial ~/ 60}min';
    } else {
      return '${initial ~/ 60}+${increment}s';
    }
  }
}