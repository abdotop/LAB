import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'pages/main_menu_page.dart';
import 'pages/waiting_room_page.dart';
import 'pages/game_board_page.dart';
import 'pages/login_page.dart';
import 'providers/auth_provider.dart';

void main() {
  runApp(const ProviderScope(child: ChessApp()));
}

class ChessApp extends ConsumerWidget {
  const ChessApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = GoRouter(
      initialLocation: '/login',
      redirect: (context, state) {
        final authState = ref.read(authProvider);
        final isLoggedIn = authState.user != null;
        
        if (!isLoggedIn && state.matchedLocation != '/login') {
          return '/login';
        }
        
        if (isLoggedIn && state.matchedLocation == '/login') {
          return '/';
        }
        
        return null;
      },
      routes: [
        GoRoute(
          path: '/login',
          builder: (context, state) => const LoginPage(),
        ),
        GoRoute(
          path: '/',
          builder: (context, state) => const MainMenuPage(),
        ),
        GoRoute(
          path: '/waiting/:lobbyId',
          builder: (context, state) => WaitingRoomPage(
            lobbyId: state.pathParameters['lobbyId']!,
          ),
        ),
        GoRoute(
          path: '/game/:gameId',
          builder: (context, state) => GameBoardPage(
            gameId: state.pathParameters['gameId']!,
          ),
        ),
      ],
    );

    return MaterialApp.router(
      title: 'Chess App',
      theme: ThemeData(
        primarySwatch: Colors.brown,
        fontFamily: 'Inter',
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      routerConfig: router,
    );
  }
}