import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../services/api_service.dart';

class WaitingRoomPage extends ConsumerStatefulWidget {
  final String lobbyId;
  
  const WaitingRoomPage({
    super.key,
    required this.lobbyId,
  });

  @override
  ConsumerState<WaitingRoomPage> createState() => _WaitingRoomPageState();
}

class _WaitingRoomPageState extends ConsumerState<WaitingRoomPage> {
  bool _isLoading = true;
  String? _error;
  Map<String, dynamic>? _lobby;

  @override
  void initState() {
    super.initState();
    // TODO: Implement lobby polling or WebSocket updates
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Waiting Room'),
        backgroundColor: Colors.brown,
        foregroundColor: Colors.white,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.go('/'),
        ),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text('Error: $_error'),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: () => context.go('/'),
                        child: const Text('Back to Menu'),
                      ),
                    ],
                  ),
                )
              : Padding(
                  padding: const EdgeInsets.all(24.0),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(
                        Icons.hourglass_empty,
                        size: 80,
                        color: Colors.brown,
                      ),
                      const SizedBox(height: 24),
                      const Text(
                        'Waiting for opponent...',
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                        textAlign: TextAlign.center,
                      ),
                      const SizedBox(height: 16),
                      if (_lobby?['inviteCode'] != null) ...[
                        const Text(
                          'Share this code with your friend:',
                          style: TextStyle(fontSize: 16),
                        ),
                        const SizedBox(height: 8),
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 24,
                            vertical: 12,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.grey[200],
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            _lobby!['inviteCode'],
                            style: const TextStyle(
                              fontSize: 32,
                              fontWeight: FontWeight.bold,
                              letterSpacing: 4,
                            ),
                          ),
                        ),
                        const SizedBox(height: 16),
                        ElevatedButton.icon(
                          onPressed: () {
                            // TODO: Share invite code
                          },
                          icon: const Icon(Icons.share),
                          label: const Text('Share Code'),
                        ),
                      ],
                      const SizedBox(height: 48),
                      const CircularProgressIndicator(
                        color: Colors.brown,
                      ),
                      const SizedBox(height: 32),
                      ElevatedButton(
                        onPressed: () async {
                          try {
                            await ref
                                .read(apiServiceProvider)
                                .deleteLobby(widget.lobbyId);
                            if (context.mounted) {
                              context.go('/');
                            }
                          } catch (e) {
                            if (context.mounted) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(content: Text('Error: $e')),
                              );
                            }
                          }
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.red,
                          foregroundColor: Colors.white,
                        ),
                        child: const Text('Cancel'),
                      ),
                    ],
                  ),
                ),
    );
  }
}