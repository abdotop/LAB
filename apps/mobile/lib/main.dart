import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import 'providers/auth_provider.dart';
import 'providers/sponsorship_provider.dart';
import 'screens/splash_screen.dart';
import 'screens/auth/login_screen.dart';
import 'screens/auth/register_screen.dart';
import 'screens/home/home_screen.dart';
import 'screens/sponsorship/candidates_screen.dart';
import 'screens/sponsorship/sponsorship_screen.dart';
import 'utils/theme.dart';

void main() {
  runApp(const EParrainagesApp());
}

class EParrainagesApp extends StatelessWidget {
  const EParrainagesApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProvider(create: (_) => SponsorshipProvider()),
      ],
      child: Consumer<AuthProvider>(
        builder: (context, authProvider, _) {
          return MaterialApp.router(
            title: 'E-Parrainages Sénégal',
            theme: AppTheme.lightTheme,
            debugShowCheckedModeBanner: false,
            routerConfig: _router,
            locale: const Locale('fr', 'SN'),
            supportedLocales: const [
              Locale('fr', 'SN'), // Français Sénégal
            ],
          );
        },
      ),
    );
  }
}

final GoRouter _router = GoRouter(
  initialLocation: '/splash',
  routes: [
    GoRoute(
      path: '/splash',
      builder: (context, state) => const SplashScreen(),
    ),
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginScreen(),
    ),
    GoRoute(
      path: '/register',
      builder: (context, state) => const RegisterScreen(),
    ),
    GoRoute(
      path: '/home',
      builder: (context, state) => const HomeScreen(),
    ),
    GoRoute(
      path: '/candidates',
      builder: (context, state) => const CandidatesScreen(),
    ),
    GoRoute(
      path: '/sponsorship/:candidateId',
      builder: (context, state) => SponsorshipScreen(
        candidateId: state.pathParameters['candidateId']!,
      ),
    ),
  ],
);