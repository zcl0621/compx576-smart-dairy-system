import 'package:flutter/material.dart';
import '../../core/mock/appMockData.dart';
import '../../core/models/appModels.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'dashboardWidgets.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  int page = 1;
  bool isRefreshing = false;

  Future<void> refreshData() async {
    setState(() => isRefreshing = true);
    await Future<void>.delayed(const Duration(milliseconds: 520));
    if (!mounted) return;
    setState(() => isRefreshing = false);
  }

  @override
  Widget build(BuildContext context) {
    const perPage = 18;
    final items = dashboardCows;
    final totalPages = (items.length / perPage).ceil();
    final pageItems = items.skip((page - 1) * perPage).take(perPage).toList();
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
                  onPressed: isRefreshing ? null : refreshData,
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.refresh_rounded, size: 18),
                  label: Text(isRefreshing ? 'Refreshing' : 'Refresh data'),
                ),
              ),
            ),
            const SizedBox(height: 24),
            LayoutBuilder(
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
                      value: items.length,
                      color: AppColors.foreground,
                    ),
                    DashboardStatCard(
                      label: 'Normal',
                      value: items
                          .where(
                            (item) => item.condition == CowCondition.normal,
                          )
                          .length,
                      color: AppColors.normal,
                    ),
                    DashboardStatCard(
                      label: 'Warning',
                      value: items
                          .where(
                            (item) => item.condition == CowCondition.warning,
                          )
                          .length,
                      color: AppColors.warning,
                    ),
                    DashboardStatCard(
                      label: 'Critical',
                      value: items
                          .where(
                            (item) => item.condition == CowCondition.critical,
                          )
                          .length,
                      color: AppColors.critical,
                    ),
                    DashboardStatCard(
                      label: 'Offline',
                      value: items
                          .where(
                            (item) => item.condition == CowCondition.offline,
                          )
                          .length,
                      color: AppColors.offline,
                    ),
                  ],
                );
              },
            ),
            const SizedBox(height: 28),
            isRefreshing
                ? const LoadingStateCard(
                    message: 'Refreshing dashboard cards',
                    lines: 12,
                  )
                : LayoutBuilder(
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
                      return Wrap(
                        spacing: 16,
                        runSpacing: 16,
                        children: pageItems
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
                    items.isEmpty
                        ? 'Showing 0 results'
                        : 'Showing ${(page - 1) * perPage + 1} to ${((page - 1) * perPage) + pageItems.length} of ${items.length} results',
                    style: const TextStyle(color: AppColors.mutedForeground),
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
        ),
      ),
    );
  }
}
