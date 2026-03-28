import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../features/alerts/alertsPage.dart';
import '../../features/auth/authPages.dart';
import '../../features/cows/cowPages.dart';
import '../../features/dashboard/dashboardPage.dart';
import '../../features/reports/reportsPage.dart';
import '../../features/users/userPages.dart';
import '../../shared/appShell.dart';

Page<void> _fadePage(Widget child) => CustomTransitionPage<void>(
  child: child,
  transitionDuration: const Duration(milliseconds: 150),
  reverseTransitionDuration: const Duration(milliseconds: 150),
  transitionsBuilder: (context, animation, _, child) =>
      FadeTransition(opacity: animation, child: child),
);

final appRouter = GoRouter(
  initialLocation: '/login',
  routes: [
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
