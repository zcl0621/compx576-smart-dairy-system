import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../core/providers/auth_provider.dart';
import '../core/providers/data_providers.dart';
import '../core/theme/appTheme.dart';
import 'appWidgets.dart';

class AppShell extends ConsumerWidget {
  const AppShell({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final compact = isCompact(context);

    return Scaffold(
      backgroundColor: AppColors.backgroundAccent,
      drawer: compact
          ? const Drawer(
              backgroundColor: AppColors.surfaceMuted,
              child: SafeArea(child: _SidebarContent()),
            )
          : null,
      body: SafeArea(
        bottom: false,
        child: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [AppColors.background, AppColors.backgroundAccent],
            ),
          ),
          child: Row(
            children: [
              if (!compact)
                const SizedBox(width: 262, child: _SidebarContent()),
              Expanded(
                child: Column(
                  children: [
                    _TopBar(compact: compact),
                    Expanded(child: child),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _SidebarContent extends StatelessWidget {
  const _SidebarContent();

  @override
  Widget build(BuildContext context) {
    final location = GoRouterState.of(context).uri.toString();
    final items = <({String label, String path, IconData icon})>[
      (label: 'Dashboard', path: '/', icon: Icons.dashboard_outlined),
      (label: 'Cows', path: '/cows', icon: Icons.assignment_outlined),
      (label: 'Alerts', path: '/alerts', icon: Icons.warning_amber_outlined),
      (label: 'Reports', path: '/reports', icon: Icons.description_outlined),
      (label: 'Users', path: '/users', icon: Icons.people_outline_rounded),
    ];

    return Container(
      decoration: const BoxDecoration(
        color: AppColors.surfaceMuted,
        border: Border(right: BorderSide(color: AppColors.border)),
      ),
      child: Column(
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.fromLTRB(20, 26, 20, 20),
            decoration: const BoxDecoration(
              color: AppColors.surfaceMuted,
              border: Border(bottom: BorderSide(color: AppColors.border)),
            ),
            child: const Row(
              children: [
                _SidebarBrandMark(),
                SizedBox(width: 12),
                Expanded(
                  child: Text(
                    'Farm Health',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
                  ),
                ),
              ],
            ),
          ),
          Expanded(
            child: ListView(
              padding: const EdgeInsets.fromLTRB(18, 18, 18, 18),
              children: [
                for (final item in items)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 10),
                    child: _SidebarTile(
                      item: item,
                      active:
                          location == item.path ||
                          (item.path != '/' && location.startsWith(item.path)),
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _SidebarBrandMark extends StatelessWidget {
  const _SidebarBrandMark();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 42,
      height: 42,
      decoration: BoxDecoration(
        color: const Color(0xFFF1E8CF),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: AppColors.border),
      ),
      child: const Icon(Icons.agriculture_rounded, color: AppColors.primary),
    );
  }
}

class _SidebarTile extends StatelessWidget {
  const _SidebarTile({required this.item, required this.active});

  final ({String label, String path, IconData icon}) item;
  final bool active;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: active ? AppColors.primary : Colors.transparent,
      borderRadius: BorderRadius.circular(10),
      child: InkWell(
        borderRadius: BorderRadius.circular(10),
        onTap: () {
          context.go(item.path);
          if (Scaffold.maybeOf(context)?.isDrawerOpen ?? false) {
            Navigator.of(context).pop();
          }
        },
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
          child: Row(
            children: [
              Icon(
                item.icon,
                size: 22,
                color: active ? Colors.white : AppColors.foreground,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  item.label,
                  style: TextStyle(
                    color: active ? Colors.white : AppColors.foreground,
                    fontWeight: active ? FontWeight.w700 : FontWeight.w500,
                    fontSize: 15,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _TopBar extends ConsumerWidget {
  const _TopBar({required this.compact});

  final bool compact;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final activeAlerts =
        ref.watch(activeAlertsProvider).valueOrNull ?? [];

    return Container(
      height: 68,
      padding: const EdgeInsets.symmetric(horizontal: 18),
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(bottom: BorderSide(color: AppColors.border)),
      ),
      child: Row(
        children: [
          if (compact)
            Builder(
              builder: (context) => IconButton(
                onPressed: () => Scaffold.of(context).openDrawer(),
                style: IconButton.styleFrom(
                  backgroundColor: AppColors.surfaceMuted,
                ),
                icon: const Icon(Icons.menu_rounded),
              ),
            ),
          const Spacer(),
          PopupMenuButton<String>(
            tooltip: 'Recent alerts',
            itemBuilder: (context) => [
              for (final alert in activeAlerts.take(5))
                PopupMenuItem<String>(
                  value: alert.cow_id,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 10,
                  ),
                  child: SizedBox(
                    width: 280,
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.only(top: 2),
                          child: Icon(
                            alert.severity.name == 'offline'
                                ? Icons.wifi_off_rounded
                                : Icons.warning_amber_rounded,
                            color: alert.severity.color,
                            size: 20,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                alert.cow_name,
                                style: const TextStyle(
                                  fontWeight: FontWeight.w700,
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                alert.message,
                                maxLines: 2,
                                overflow: TextOverflow.ellipsis,
                                style: const TextStyle(
                                  color: AppColors.mutedForeground,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
            ],
            onSelected: (value) => context.go('/cows/$value'),
            child: Stack(
              clipBehavior: Clip.none,
              children: [
                const Padding(
                  padding: EdgeInsets.all(8),
                  child: Icon(
                    Icons.notifications_none_rounded,
                    size: 22,
                    color: AppColors.foreground,
                  ),
                ),
                if (activeAlerts.isNotEmpty)
                  Positioned(
                    right: 4,
                    top: 4,
                    child: Container(
                      width: 8,
                      height: 8,
                      decoration: const BoxDecoration(
                        color: AppColors.critical,
                        shape: BoxShape.circle,
                      ),
                    ),
                  ),
              ],
            ),
          ),
          const SizedBox(width: 16),
          Container(width: 1, height: 36, color: AppColors.border),
          const SizedBox(width: 16),
          PopupMenuButton<String>(
            onSelected: (value) {
              if (value == 'logout') {
                ref.read(authProvider.notifier).logout();
                context.go('/login');
              }
            },
            itemBuilder: (_) => const [
              PopupMenuItem<String>(value: 'logout', child: Text('Logout')),
            ],
            child: Row(
              children: [
                CircleAvatar(
                  radius: 16,
                  backgroundColor: AppColors.surfaceMuted,
                  child: Icon(
                    Icons.person_outline_rounded,
                    color: AppColors.mutedForeground,
                    size: 18,
                  ),
                ),
                if (!compact) ...[
                  const SizedBox(width: 12),
                  const Text(
                    'Farm Manager',
                    style: TextStyle(fontWeight: FontWeight.w600),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}
