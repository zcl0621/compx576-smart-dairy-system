import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/mock/appMockData.dart';
import '../../core/models/appModels.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'alertsWidgets.dart';

class AlertsPage extends StatefulWidget {
  const AlertsPage({super.key});

  @override
  State<AlertsPage> createState() => _AlertsPageState();
}

class _AlertsPageState extends State<AlertsPage> {
  AlertSeverity? severity;
  int page = 1;

  static const pageSize = 8;

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final compact = width < 760;
    final alertsList = _getAlertsList();
    final totalPages = alertsList.isEmpty
        ? 1
        : (alertsList.length / pageSize).ceil();
    final safePage = page.clamp(1, totalPages);
    final pageItems = alertsList
        .skip((safePage - 1) * pageSize)
        .take(pageSize)
        .toList();

    return SingleChildScrollView(
      child: PageSection(
        padding: EdgeInsets.fromLTRB(
          compact ? 16 : 24,
          24,
          compact ? 16 : 24,
          32,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const PageIntro(
              title: 'Alerts',
              subtitle: 'Monitor and manage system alerts',
            ),
            const SizedBox(height: 20),
            GridView.count(
              crossAxisCount: width > 1080 ? 4 : 2,
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
              childAspectRatio: compact ? 1.55 : 1.75,
              children: [
                AlertsStatCard(
                  label: 'Active',
                  value: alertsList
                      .where((item) => item.status == AlertStatus.active)
                      .length,
                  color: AppColors.foreground,
                ),
                AlertsStatCard(
                  label: 'Warning',
                  value: alertsList
                      .where((item) => item.severity == AlertSeverity.warning)
                      .length,
                  color: AppColors.warning,
                ),
                AlertsStatCard(
                  label: 'Critical',
                  value: alertsList
                      .where((item) => item.severity == AlertSeverity.critical)
                      .length,
                  color: AppColors.critical,
                ),
                AlertsStatCard(
                  label: 'Offline',
                  value: alertsList
                      .where((item) => item.severity == AlertSeverity.offline)
                      .length,
                  color: AppColors.offline,
                ),
              ],
            ),
            const SizedBox(height: 16),
            SurfaceCard(
              child: SizedBox(
                width: double.infinity,
                child: Wrap(
                  spacing: 12,
                  runSpacing: 12,
                  crossAxisAlignment: WrapCrossAlignment.center,
                  children: [
                    const Icon(
                      Icons.filter_alt_outlined,
                      size: 18,
                      color: AppColors.mutedForeground,
                    ),
                    const Text(
                      'Severity:',
                      style: TextStyle(fontWeight: FontWeight.w600),
                    ),
                    SizedBox(
                      width: 160,
                      child: AppSelect<AlertSeverity?>(
                        value: severity,
                        hint: 'All',
                        prefixIcon: Icons.filter_alt_outlined,
                        items: const [
                          (text: 'All', value: null),
                          (text: 'Warning', value: AlertSeverity.warning),
                          (text: 'Critical', value: AlertSeverity.critical),
                          (text: 'Offline', value: AlertSeverity.offline),
                        ],
                        onChanged: _changeSeverity,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),
            SurfaceCard(
              padding: const EdgeInsets.all(0),
              child: Column(
                children: [
                  if (pageItems.isEmpty)
                    const Padding(
                      padding: EdgeInsets.all(24),
                      child: EmptyStateCard(
                        message: 'No alerts found',
                        icon: Icons.notifications_off_outlined,
                      ),
                    )
                  else
                    for (final item in pageItems)
                      AlertRow(
                        alert: item,
                        onTap: () => context.go('/cows/${item.cow_id}'),
                      ),
                  Container(
                    padding: const EdgeInsets.fromLTRB(16, 14, 16, 16),
                    decoration: const BoxDecoration(
                      border: Border(top: BorderSide(color: AppColors.border)),
                    ),
                    child: compact
                        ? Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                _resultLabel(
                                  safePage,
                                  pageItems.length,
                                  alertsList.length,
                                ),
                                style: const TextStyle(
                                  color: AppColors.mutedForeground,
                                ),
                              ),
                              const SizedBox(height: 12),
                              AppPagination(
                                currentPage: safePage,
                                totalPages: totalPages,
                                onChanged: (value) =>
                                    setState(() => page = value),
                              ),
                            ],
                          )
                        : Row(
                            children: [
                              Expanded(
                                child: Text(
                                  _resultLabel(
                                    safePage,
                                    pageItems.length,
                                    alertsList.length,
                                  ),
                                  style: const TextStyle(
                                    color: AppColors.mutedForeground,
                                  ),
                                ),
                              ),
                              AppPagination(
                                currentPage: safePage,
                                totalPages: totalPages,
                                onChanged: (value) =>
                                    setState(() => page = value),
                              ),
                            ],
                          ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  List<AlertItem> _getAlertsList() {
    return alerts.where((item) {
      if (severity == null) return true;
      return item.severity == severity;
    }).toList();
  }

  void _changeSeverity(AlertSeverity? value) {
    setState(() {
      severity = value;
      page = 1;
    });
  }

  String _resultLabel(int currentPage, int currentCount, int totalCount) {
    if (totalCount == 0) return 'Showing 0 results';
    final start = ((currentPage - 1) * pageSize) + 1;
    final end = start + currentCount - 1;
    return 'Showing $start to $end of $totalCount results';
  }
}
