import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/models/appModels.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';

class ReportsListCard extends StatelessWidget {
  const ReportsListCard({
    super.key,
    required this.reportsList,
    required this.selected,
    required this.onSelect,
    required this.currentPage,
    required this.totalPages,
    required this.onPageChanged,
    required this.totalCount,
    required this.pageSize,
  });

  final List<ReportItem> reportsList;
  final ReportItem selected;
  final ValueChanged<ReportItem> onSelect;
  final int currentPage;
  final int totalPages;
  final ValueChanged<int> onPageChanged;
  final int totalCount;
  final int pageSize;

  @override
  Widget build(BuildContext context) {
    return SurfaceCard(
      padding: const EdgeInsets.all(0),
      child: Column(
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            decoration: const BoxDecoration(
              color: AppColors.surfaceMuted,
              borderRadius: BorderRadius.vertical(top: Radius.circular(22)),
            ),
            child: Text(
              'All Reports',
              style: Theme.of(context).textTheme.titleMedium,
            ),
          ),
          for (final report in reportsList)
            ReportRow(
              report: report,
              selected: report.id == selected.id,
              onTap: () => onSelect(report),
            ),
          Container(
            padding: const EdgeInsets.fromLTRB(16, 14, 16, 16),
            decoration: const BoxDecoration(
              border: Border(top: BorderSide(color: AppColors.border)),
            ),
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    'Showing ${(currentPage - 1) * pageSize + 1} to ${((currentPage - 1) * pageSize) + reportsList.length} of $totalCount results',
                    style: const TextStyle(color: AppColors.mutedForeground),
                  ),
                ),
                AppPagination(
                  currentPage: currentPage,
                  totalPages: totalPages,
                  onChanged: onPageChanged,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class ReportRow extends StatelessWidget {
  const ReportRow({
    super.key,
    required this.report,
    required this.selected,
    required this.onTap,
  });

  final ReportItem report;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final scoreColor = reportScoreColor(report.score);
    return Material(
      color: selected ? const Color(0xFFF5F7F7) : Colors.transparent,
      child: InkWell(
        onTap: onTap,
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.fromLTRB(18, 16, 14, 14),
          decoration: BoxDecoration(
            border: Border(
              left: BorderSide(
                color: selected ? AppColors.primary : Colors.transparent,
                width: 4,
              ),
              top: const BorderSide(color: AppColors.border),
            ),
          ),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      report.cow_name,
                      style: const TextStyle(fontWeight: FontWeight.w700),
                    ),
                    const SizedBox(height: 4),
                    const Text(
                      'Last 7 days',
                      style: TextStyle(
                        color: AppColors.mutedForeground,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(report.summary, style: const TextStyle(height: 1.35)),
                    const SizedBox(height: 8),
                    Text(
                      reportDate(report.created_at),
                      style: const TextStyle(
                        color: AppColors.mutedForeground,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 12),
              Text(
                report.score.toStringAsFixed(0),
                style: TextStyle(
                  color: scoreColor,
                  fontSize: 22,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class ReportDetailsCard extends StatelessWidget {
  const ReportDetailsCard({super.key, required this.selected});

  final ReportItem selected;

  @override
  Widget build(BuildContext context) {
    final scoreColor = reportScoreColor(selected.score);
    return SurfaceCard(
      padding: const EdgeInsets.all(0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            decoration: const BoxDecoration(
              color: AppColors.surfaceMuted,
              borderRadius: BorderRadius.vertical(top: Radius.circular(22)),
            ),
            child: Text(
              'Report Details',
              style: Theme.of(context).textTheme.titleMedium,
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(22),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            selected.cow_name,
                            style: Theme.of(context).textTheme.headlineSmall,
                          ),
                          const SizedBox(height: 6),
                          const Text(
                            'Last 7 days',
                            style: TextStyle(
                              color: AppColors.mutedForeground,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                    TextButton(
                      onPressed: () => context.go('/cows/${selected.cow_id}'),
                      child: const Text('View Cow'),
                    ),
                  ],
                ),
                const SizedBox(height: 14),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(18),
                  decoration: BoxDecoration(
                    color: AppColors.surfaceMuted,
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Health Score',
                        style: TextStyle(color: AppColors.mutedForeground),
                      ),
                      const SizedBox(height: 10),
                      Text(
                        selected.score.toStringAsFixed(0),
                        style: TextStyle(
                          fontSize: 44,
                          fontWeight: FontWeight.w700,
                          color: scoreColor,
                        ),
                      ),
                      const SizedBox(height: 12),
                      ClipRRect(
                        borderRadius: BorderRadius.circular(999),
                        child: LinearProgressIndicator(
                          value: selected.score / 100,
                          minHeight: 12,
                          color: scoreColor,
                          backgroundColor: Colors.white,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 22),
                Text('Summary', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 10),
                Text(selected.summary, style: const TextStyle(height: 1.4)),
                const SizedBox(height: 22),
                Text('Details', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 10),
                Text(
                  selected.details,
                  style: const TextStyle(
                    color: AppColors.mutedForeground,
                    height: 1.45,
                  ),
                ),
                const SizedBox(height: 24),
                const Divider(height: 1),
                const SizedBox(height: 18),
                Row(
                  children: [
                    const Icon(
                      Icons.description_outlined,
                      size: 18,
                      color: AppColors.mutedForeground,
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        'Generated on ${reportDateTime(selected.created_at)}',
                        style: const TextStyle(
                          color: AppColors.mutedForeground,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

String reportDate(DateTime value) =>
    '${value.year}/${value.month}/${value.day}';

String reportDateTime(DateTime value) =>
    '${value.year}/${value.month}/${value.day} ${value.hour.toString().padLeft(2, '0')}:${value.minute.toString().padLeft(2, '0')}:${value.second.toString().padLeft(2, '0')}';

Color reportScoreColor(double score) {
  if (score >= 85) return AppColors.normal;
  if (score >= 70) return AppColors.warning;
  return AppColors.critical;
}
