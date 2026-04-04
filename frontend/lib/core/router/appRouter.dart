import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/alerts/alertsPage.dart';
import '../../features/auth/authPages.dart';
import '../../features/cows/cowPages.dart';
import '../../features/dashboard/dashboardPage.dart';
import '../../features/reports/reportsPage.dart';
import '../../features/users/userPages.dart';
import '../../shared/appShell.dart';
import '../providers/auth_provider.dart';

Page<void> _fadePage(Widget child) => CustomTransitionPage<void>(
  child: child,
  transitionDuration: const Duration(milliseconds: 150),
  reverseTransitionDuration: const Duration(milliseconds: 150),
  transitionsBuilder: (context, animation, _, child) =>
      FadeTransition(opacity: animation, child: child),
);

/// tracks whether the initial token refresh has completed
final authInitProvider = FutureProvider<bool>((ref) async {
  return ref.read(authProvider.notifier).tryRefresh();
});

/// preserves the deep link path across the splash redirect
final _pendingPathProvider = StateProvider<String?>((ref) => null);

final appRouterProvider = Provider<GoRouter>((ref) {
  final router = GoRouter(
    initialLocation: '/splash',
    redirect: (context, state) {
      final initDone = ref.read(authInitProvider) is! AsyncLoading;
      if (!initDone) {
        if (state.uri.path != '/splash') {
          // save deep link so we can restore it after auth init
          ref.read(_pendingPathProvider.notifier).state = state.uri.toString();
          return '/splash';
        }
        return null;
      }

      final auth = ref.read(authProvider);
      final onAuthPage = state.uri.path == '/login' ||
          state.uri.path.startsWith('/forgot-password') ||
          state.uri.path.startsWith('/verify-code') ||
          state.uri.path.startsWith('/reset-password') ||
          state.uri.path.startsWith('/reset-success') ||
          state.uri.path == '/splash';

      if (!auth.isLoggedIn && !onAuthPage) return '/login';
      if (auth.isLoggedIn && (state.uri.path == '/login' || state.uri.path == '/splash')) {
        final pending = ref.read(_pendingPathProvider);
        if (pending != null) {
          ref.read(_pendingPathProvider.notifier).state = null;
          return pending;
        }
        return '/';
      }
      // init done, not logged in, on splash -> go to login
      if (!auth.isLoggedIn && state.uri.path == '/splash') return '/login';
      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (context, state) => const Scaffold(
          body: Center(child: CircularProgressIndicator()),
        ),
      ),
      GoRoute(path: '/login', builder: (context, state) => const LoginPage()),
      GoRoute(
        path: '/forgot-password',
        builder: (context, state) => const ForgotPasswordPage(),
      ),
      GoRoute(
        path: '/verify-code',
        builder: (context, state) => const VerifyCodePage(),
      ),
      GoRoute(
        path: '/reset-password',
        builder: (context, state) => const ResetPasswordPage(),
      ),
      GoRoute(
        path: '/reset-success',
        builder: (context, state) => const ResetSuccessPage(),
      ),
      ShellRoute(
        builder: (context, state, child) => AppShell(child: child),
        routes: [
          GoRoute(
            path: '/',
            pageBuilder: (context, state) => _fadePage(const DashboardPage()),
          ),
          GoRoute(
            path: '/cows',
            pageBuilder: (context, state) => _fadePage(const CowsListPage()),
          ),
          GoRoute(
            path: '/cows/add',
            pageBuilder: (context, state) => _fadePage(const AddCowPage()),
          ),
          GoRoute(
            path: '/cows/:id',
            pageBuilder: (context, state) =>
                _fadePage(CowDetailPage(id: state.pathParameters['id']!)),
          ),
          GoRoute(
            path: '/cows/:id/edit',
            pageBuilder: (context, state) =>
                _fadePage(EditCowPage(id: state.pathParameters['id']!)),
          ),
          GoRoute(
            path: '/alerts',
            pageBuilder: (context, state) => _fadePage(const AlertsPage()),
          ),
          GoRoute(
            path: '/reports',
            pageBuilder: (context, state) => _fadePage(const ReportsPage()),
          ),
          GoRoute(
            path: '/users',
            pageBuilder: (context, state) => _fadePage(const UsersListPage()),
          ),
          GoRoute(
            path: '/users/add',
            pageBuilder: (context, state) => _fadePage(const AddUserPage()),
          ),
          GoRoute(
            path: '/users/:id/edit',
            pageBuilder: (context, state) =>
                _fadePage(EditUserPage(id: state.pathParameters['id']!)),
          ),
        ],
      ),
    ],
  );

  // re-evaluate redirects when auth init completes or auth state changes
  ref.listen(authInitProvider, (_, __) => router.refresh());
  ref.listen(authProvider, (_, __) => router.refresh());

  return router;
});
