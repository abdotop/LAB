import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/auth_provider.dart';
import '../services/api_service.dart';
import '../models/lobby.dart';

class MainMenuPage extends ConsumerWidget {
  const MainMenuPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final user = authState.user!;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Chess App'),
        backgroundColor: Colors.brown,
        foregroundColor: Colors.white,
        actions: [
          PopupMenuButton(
            icon: CircleAvatar(
              backgroundColor: Colors.white,
              child: Text(
                user.username[0].toUpperCase(),
                style: const TextStyle(
                  color: Colors.brown,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            itemBuilder: (context) => [
              PopupMenuItem(
                child: Text('${user.username} (${user.rating})'),
                onTap: () {},
              ),
              const PopupMenuDivider(),
              PopupMenuItem(
                child: const Text('Logout'),
                onTap: () {
                  ref.read(authProvider.notifier).logout();
                  context.go('/login');
                },
              ),
            ],
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // Welcome message
            Text(
              'Welcome, ${user.username}!',
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 8),
            Text(
              'Rating: ${user.rating}',
              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                color: Colors.grey[600],
              ),
            ),
            const SizedBox(height: 48),

            // Menu buttons
            _MenuButton(
              icon: Icons.public,
              title: 'Play Online',
              subtitle: 'Find an opponent',
              onTap: () => _showPublicLobbies(context, ref),
            ),
            const SizedBox(height: 16),
            
            _MenuButton(
              icon: Icons.add,
              title: 'Create Private Game',
              subtitle: 'Invite friends',
              onTap: () => _createPrivateLobby(context, ref),
            ),
            const SizedBox(height: 16),
            
            _MenuButton(
              icon: Icons.login,
              title: 'Join with Code',
              subtitle: 'Enter invite code',
              onTap: () => _showJoinCodeDialog(context, ref),
            ),
            const SizedBox(height: 16),
            
            _MenuButton(
              icon: Icons.person,
              title: 'Profile',
              subtitle: 'View your stats',
              onTap: () {
                // TODO: Navigate to profile page
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showPublicLobbies(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => const PublicLobbiesSheet(),
    );
  }

  void _createPrivateLobby(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => const CreateLobbyDialog(),
    );
  }

  void _showJoinCodeDialog(BuildContext context, WidgetRef ref) {
    final codeController = TextEditingController();
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Join Game'),
        content: TextField(
          controller: codeController,
          decoration: const InputDecoration(
            labelText: 'Enter invite code',
            border: OutlineInputBorder(),
          ),
          textCapitalization: TextCapitalization.characters,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              final code = codeController.text.trim();
              if (code.isNotEmpty) {
                try {
                  final apiService = ref.read(apiServiceProvider);
                  final response = await apiService.joinLobby(code);
                  
                  if (context.mounted) {
                    Navigator.of(context).pop();
                    if (response['game'] != null) {
                      context.go('/game/${response['game']['id']}');
                    } else {
                      context.go('/waiting/${response['lobby']['id']}');
                    }
                  }
                } catch (e) {
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('Error: $e')),
                    );
                  }
                }
              }
            },
            child: const Text('Join'),
          ),
        ],
      ),
    );
  }
}

class _MenuButton extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  const _MenuButton({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: Icon(icon, size: 32, color: Colors.brown),
        title: Text(
          title,
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        subtitle: Text(subtitle),
        trailing: const Icon(Icons.chevron_right),
        onTap: onTap,
      ),
    );
  }
}

class PublicLobbiesSheet extends ConsumerStatefulWidget {
  const PublicLobbiesSheet({super.key});

  @override
  ConsumerState<PublicLobbiesSheet> createState() => _PublicLobbiesSheetState();
}

class _PublicLobbiesSheetState extends ConsumerState<PublicLobbiesSheet> {
  late Future<List<Lobby>> _lobbiesFuture;

  @override
  void initState() {
    super.initState();
    _lobbiesFuture = ref.read(apiServiceProvider).getPublicLobbies();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: MediaQuery.of(context).size.height * 0.7,
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Public Games',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              IconButton(
                onPressed: () => Navigator.of(context).pop(),
                icon: const Icon(Icons.close),
              ),
            ],
          ),
          const Divider(),
          Expanded(
            child: FutureBuilder<List<Lobby>>(
              future: _lobbiesFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                }

                if (snapshot.hasError) {
                  return Center(
                    child: Text('Error: ${snapshot.error}'),
                  );
                }

                final lobbies = snapshot.data ?? [];
                
                if (lobbies.isEmpty) {
                  return const Center(
                    child: Text('No public games available'),
                  );
                }

                return ListView.builder(
                  itemCount: lobbies.length,
                  itemBuilder: (context, index) {
                    final lobby = lobbies[index];
                    return Card(
                      child: ListTile(
                        title: Text(lobby.creator.username),
                        subtitle: Text('Time: ${lobby.timeControl.displayText}'),
                        trailing: ElevatedButton(
                          onPressed: () async {
                            try {
                              final apiService = ref.read(apiServiceProvider);
                              final response = await apiService.joinLobby(lobby.id);
                              
                              if (context.mounted) {
                                Navigator.of(context).pop();
                                context.go('/game/${response['game']['id']}');
                              }
                            } catch (e) {
                              if (context.mounted) {
                                ScaffoldMessenger.of(context).showSnackBar(
                                  SnackBar(content: Text('Error: $e')),
                                );
                              }
                            }
                          },
                          child: const Text('Join'),
                        ),
                      ),
                    );
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

class CreateLobbyDialog extends ConsumerStatefulWidget {
  const CreateLobbyDialog({super.key});

  @override
  ConsumerState<CreateLobbyDialog> createState() => _CreateLobbyDialogState();
}

class _CreateLobbyDialogState extends ConsumerState<CreateLobbyDialog> {
  bool _isPrivate = true;
  int _initialTime = 600; // 10 minutes
  int _increment = 0;

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Create Game'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SwitchListTile(
            title: const Text('Private Game'),
            value: _isPrivate,
            onChanged: (value) => setState(() => _isPrivate = value),
          ),
          const SizedBox(height: 16),
          Text('Initial Time: ${_initialTime ~/ 60} minutes'),
          Slider(
            value: _initialTime.toDouble(),
            min: 60,
            max: 1800,
            divisions: 29,
            onChanged: (value) => setState(() => _initialTime = value.toInt()),
          ),
          Text('Increment: $_increment seconds'),
          Slider(
            value: _increment.toDouble(),
            min: 0,
            max: 30,
            divisions: 30,
            onChanged: (value) => setState(() => _increment = value.toInt()),
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () async {
            try {
              final apiService = ref.read(apiServiceProvider);
              final lobby = await apiService.createLobby(
                type: _isPrivate ? 'PRIVATE' : 'PUBLIC',
                timeControl: TimeControl(
                  initial: _initialTime,
                  increment: _increment,
                ),
                inviteOnly: _isPrivate,
              );
              
              if (context.mounted) {
                Navigator.of(context).pop();
                context.go('/waiting/${lobby.id}');
              }
            } catch (e) {
              if (context.mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Error: $e')),
                );
              }
            }
          },
          child: const Text('Create'),
        ),
      ],
    );
  }
}