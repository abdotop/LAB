import 'dart:convert';
import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

final websocketServiceProvider = Provider<WebSocketService>((ref) {
  return WebSocketService();
});

class WebSocketService {
  static const String wsUrl = 'ws://localhost:8080/ws';
  WebSocketChannel? _channel;
  StreamController<Map<String, dynamic>>? _messageController;

  Stream<Map<String, dynamic>> get messages {
    _messageController ??= StreamController<Map<String, dynamic>>.broadcast();
    return _messageController!.stream;
  }

  void connect(String authToken) {
    final uri = Uri.parse('$wsUrl?token=$authToken');
    _channel = WebSocketChannel.connect(uri);
    
    _messageController ??= StreamController<Map<String, dynamic>>.broadcast();
    
    _channel!.stream.listen(
      (data) {
        try {
          final message = json.decode(data);
          _messageController!.add(message);
        } catch (e) {
          print('Failed to parse WebSocket message: $e');
        }
      },
      onError: (error) {
        print('WebSocket error: $error');
        _messageController!.addError(error);
      },
      onDone: () {
        print('WebSocket connection closed');
        _channel = null;
      },
    );

    // Send heartbeat every 30 seconds
    Timer.periodic(const Duration(seconds: 30), (timer) {
      if (_channel == null) {
        timer.cancel();
        return;
      }
      _sendMessage('heartbeat', {});
    });
  }

  void disconnect() {
    _channel?.sink.close();
    _channel = null;
    _messageController?.close();
    _messageController = null;
  }

  void _sendMessage(String type, Map<String, dynamic> data) {
    if (_channel == null) return;
    
    final message = {
      'type': type,
      'data': data,
    };
    
    _channel!.sink.add(json.encode(message));
  }

  void joinGame(String gameId) {
    _sendMessage('join', {'gameId': gameId});
  }

  void leaveGame(String gameId) {
    _sendMessage('leave', {'gameId': gameId});
  }

  void makeMove(String gameId, String from, String to, {String? promotion}) {
    final data = <String, dynamic>{
      'gameId': gameId,
      'from': from,
      'to': to,
    };
    
    if (promotion != null) {
      data['promotion'] = promotion;
    }
    
    _sendMessage('move', data);
  }

  void offerDraw(String gameId) {
    _sendMessage('offerDraw', {'gameId': gameId});
  }

  void respondDraw(String gameId, bool accept) {
    _sendMessage('respondDraw', {
      'gameId': gameId,
      'accept': accept,
    });
  }

  void resign(String gameId) {
    _sendMessage('resign', {'gameId': gameId});
  }

  bool get isConnected => _channel != null;
}