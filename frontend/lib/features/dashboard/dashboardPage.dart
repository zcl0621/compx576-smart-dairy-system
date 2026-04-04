import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/providers/data_providers.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'dashboardWidgets.dart';

class DashboardPage extends ConsumerStatefulWidget {
  const DashboardPage({super.key});

  @override
  ConsumerState<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends ConsumerState<DashboardPage> {
  int page = 1;

  void _refresh() {
    ref.invalidate(dashboardSummaryProvider);
    ref.invalidate(dashboardListProvider);
  }

  @override
  Widget build(BuildContext context) {
    final summaryAsync = ref.watch(dashboardSummaryProvider);
    final listAsync = ref.watch(dashboardListProvider(page));
    final width = MediaQuery.sizeOf(context).width;
    final horizontal = width < 700 ? 16.0 : 24.0;

    return SingleChildScrollView(
      child: PageSection(
        padding: EdgeInsets.fromLTRB(horizontal, 24, horizontal, 32),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            PageIntro(
              title: 'Dashboard',
              subtitle:
                  'Monitor herd health, active risks, and sensor coverage',
              trailing: SizedBox(
                height: 42,
                child: OutlinedButton.icon(
                  onPressed: _refresh,
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.refresh_rounded, size: 18),
                  label: const Text('Refresh data'),
                ),
              ),
            ),
            const SizedBox(height: 24),
            summaryAsync.when(
              loading: () => const LoadingStateCard(
                message: 'Loading summary',
                lines: 2,
              ),
              error: (e, _) => EmptyStateCard(message: 'Failed to load: $e'),
              data: (summary) => LayoutBuilder(
                builder: (context, constraints) {
                  final w = constraints.maxWidth;
                  final cols = w > 760
                      ? 5
                      : w > 480
                          ? 3
                          : 2;
                  return GridView(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: cols,
                      crossAxisSpacing: 16,
                      mainAxisSpacing: 16,
                      mainAxisExtent: 100,
                    ),
                    children: [
                      DashboardStatCard(
                        label: 'Total Cows',
                        value: summary.total_cows,
                        color: AppColors.foreground,
                      ),
                      DashboardStatCard(
                        label: 'Normal',
                        value: summary.normal,
                        color: AppColors.normal,
                      ),
                      DashboardStatCard(
                        label: 'Warning',
                        value: summary.warning,
                        color: AppColors.warning,
                      ),
                      DashboardStatCard(
                        label: 'Critical',
                        value: summary.critical,
                        color: AppColors.critical,
                      ),
                      DashboardStatCard(
                        label: 'Offline',
                        value: summary.offline,
                        color: AppColors.offline,
                      ),
                    ],
                  );
                },
              ),
            ),
            const SizedBox(height: 28),
            listAsync.when(
              loading: () => const LoadingStateCard(
                message: 'Refreshing dashboard cards',
                lines: 12,
              ),
              error: (e, _) => EmptyStateCard(message: 'Failed to load: $e'),
              data: (result) {
                final items = result.list;
                final totalPages = result.totalPages;
                return Column(
                  children: [
                    LayoutBuilder(
                      builder: (context, constraints) {
                        final w = constraints.maxWidth;
                        final cols = w > 900
                            ? 3
                            : w > 580
                                ? 2
                                : 1;
                        final cardWidth = cols == 3
                            ? (w - 32) / 3
                            : cols == 2
                                ? (w - 16) / 2
                                : w;
                        if (items.isEmpty) {
                          return const EmptyStateCard(
                            message: 'No cows found',
                          );
                        }
                        return Wrap(
                          spacing: 16,
                          runSpacing: 16,
                          children: items
                              .map(
                                (item) => SizedBox(
                                  width: cardWidth,
                                  child: DashboardCowCard(item: item),
                                ),
                              )
                              .toList(),
                        );
                      },
                    ),
                    const SizedBox(height: 18),
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            result.total == 0
                                ? 'Showing 0 results'
                                : 'Showing ${(page - 1) * 18 + 1} to ${(page - 1) * 18 + items.length} of ${result.total} results',
                            style: const TextStyle(
                              color: AppColors.mutedForeground,
                            ),
                          ),
                        ),
                        AppPagination(
                          currentPage: page,
                          totalPages: totalPages,
                          onChanged: (value) => setState(() => page = value),
                        ),
                      ],
                    ),
                  ],
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}
