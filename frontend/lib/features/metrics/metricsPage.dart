import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/providers/data_providers.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';

class MetricsPage extends ConsumerStatefulWidget {
  const MetricsPage({super.key});

  @override
  ConsumerState<MetricsPage> createState() => _MetricsPageState();
}

const _metricTypeItems = [
  (text: 'All', value: null),
  (text: 'Temperature', value: 'temperature'),
  (text: 'Heart Rate', value: 'heart_rate'),
  (text: 'Blood Oxygen', value: 'blood_oxygen'),
  (text: 'Weight', value: 'weight'),
  (text: 'Milk Amount', value: 'milk_amount'),
  (text: 'Milking Duration', value: 'milking_duration'),
  (text: 'Latitude', value: 'latitude'),
  (text: 'Longitude', value: 'longitude'),
];

const _pageSize = 20;

class _MetricsPageState extends ConsumerState<MetricsPage> {
  int _page = 1;
  String? _metricType;
  String? _cowId;

  final _cowIdController = TextEditingController();

  @override
  void dispose() {
    _cowIdController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final compact = width < 760;
    final params = MetricListParams(
      page: _page,
      pageSize: _pageSize,
      cowId: _cowId,
      metricType: _metricType,
    );
    final listAsync = ref.watch(metricListProvider(params));

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
              title: 'Metrics',
              subtitle: 'Browse raw metric records for all cows',
            ),
            const SizedBox(height: 20),
            SurfaceCard(
              child: compact
                  ? Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Type',
                          style: TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 13,
                            color: AppColors.mutedForeground,
                          ),
                        ),
                        const SizedBox(height: 6),
                        SizedBox(
                          width: double.infinity,
                          child: AppSelect<String?>(
                            value: _metricType,
                            hint: 'All',
                            prefixIcon: Icons.filter_alt_outlined,
                            items: _metricTypeItems,
                            onChanged: _changeMetricType,
                          ),
                        ),
                        const SizedBox(height: 12),
                        const Text(
                          'Cow',
                          style: TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 13,
                            color: AppColors.mutedForeground,
                          ),
                        ),
                        const SizedBox(height: 6),
                        TextField(
                          controller: _cowIdController,
                          style: const TextStyle(fontSize: 14),
                          decoration: const InputDecoration(
                            hintText: 'Filter by cow',
                            prefixIcon: Icon(
                              Icons.search_rounded,
                              size: 18,
                              color: AppColors.mutedForeground,
                            ),
                            isDense: true,
                            contentPadding: EdgeInsets.symmetric(
                              horizontal: 14,
                              vertical: 12,
                            ),
                          ),
                          onSubmitted: _changeCowId,
                          onChanged: (v) {
                            if (v.isEmpty) _changeCowId('');
                          },
                        ),
                      ],
                    )
                  : SizedBox(
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
                            'Type:',
                            style: TextStyle(fontWeight: FontWeight.w600),
                          ),
                          SizedBox(
                            width: 200,
                            child: AppSelect<String?>(
                              value: _metricType,
                              hint: 'All',
                              prefixIcon: Icons.filter_alt_outlined,
                              items: _metricTypeItems,
                              onChanged: _changeMetricType,
                            ),
                          ),
                          const Text(
                            'Cow:',
                            style: TextStyle(fontWeight: FontWeight.w600),
                          ),
                          SizedBox(
                            width: 160,
                            child: TextField(
                              controller: _cowIdController,
                              style: const TextStyle(fontSize: 14),
                              decoration: const InputDecoration(
                                hintText: 'Filter by cow',
                                prefixIcon: Icon(
                                  Icons.search_rounded,
                                  size: 18,
                                  color: AppColors.mutedForeground,
                                ),
                                isDense: true,
                                contentPadding: EdgeInsets.symmetric(
                                  horizontal: 14,
                                  vertical: 12,
                                ),
                              ),
                              onSubmitted: _changeCowId,
                              onChanged: (v) {
                                if (v.isEmpty) _changeCowId('');
                              },
                            ),
                          ),
                        ],
                      ),
                    ),
            ),
            const SizedBox(height: 16),
            listAsync.when(
              loading: () => const LoadingStateCard(
                message: 'Loading metrics',
                lines: 6,
              ),
              error: (e, _) => EmptyStateCard(message: 'Failed to load: $e'),
              data: (result) {
                final pageItems = result.list;
                final totalPages = result.totalPages;
                return SurfaceCard(
                  padding: const EdgeInsets.all(0),
                  child: Column(
                    children: [
                      if (pageItems.isEmpty)
                        const Padding(
                          padding: EdgeInsets.all(24),
                          child: EmptyStateCard(
                            message: 'No metrics found',
                            icon: Icons.analytics_outlined,
                          ),
                        )
                      else if (compact)
                        // mobile: card list
                        ...pageItems.map((item) => Container(
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                          decoration: const BoxDecoration(
                            border: Border(bottom: BorderSide(color: AppColors.border)),
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    item.cowName,
                                    style: const TextStyle(fontWeight: FontWeight.w600),
                                  ),
                                  Text(
                                    item.metricType.replaceAll('_', ' '),
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: AppColors.mutedForeground,
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 4),
                              Row(
                                children: [
                                  Text(
                                    '${item.metricValue.toStringAsFixed(2)} ${item.unit}',
                                    style: const TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.w700,
                                    ),
                                  ),
                                  const Spacer(),
                                  Text(
                                    item.source,
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: AppColors.mutedForeground,
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 4),
                              Text(
                                _formatTime(item.createdAt.toLocal()),
                                style: const TextStyle(
                                  fontFamily: 'monospace',
                                  fontSize: 12,
                                  color: AppColors.mutedForeground,
                                ),
                              ),
                            ],
                          ),
                        ))
                      else
                        // desktop: data table
                        SizedBox(
                          width: double.infinity,
                          child: DataTable(
                            dataRowMinHeight: 56,
                            dataRowMaxHeight: 64,
                            columns: const [
                              DataColumn(label: Text('Time')),
                              DataColumn(label: Text('Cow')),
                              DataColumn(label: Text('Type')),
                              DataColumn(label: Text('Value')),
                              DataColumn(label: Text('Unit')),
                              DataColumn(label: Text('Source')),
                            ],
                            rows: [
                              for (final item in pageItems)
                                DataRow(
                                  cells: [
                                    DataCell(
                                      Text(
                                        _formatTime(item.createdAt.toLocal()),
                                        style: const TextStyle(
                                          fontFamily: 'monospace',
                                          fontSize: 13,
                                          color: AppColors.mutedForeground,
                                        ),
                                      ),
                                    ),
                                    DataCell(
                                      Text.rich(
                                        TextSpan(
                                          children: [
                                            TextSpan(
                                              text: item.cowName,
                                              style: const TextStyle(
                                                fontWeight: FontWeight.w600,
                                              ),
                                            ),
                                            TextSpan(
                                              text: '\n${item.cowId}',
                                              style: const TextStyle(
                                                fontSize: 12,
                                                color: AppColors.mutedForeground,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                    ),
                                    DataCell(
                                      Text(
                                        item.metricType.replaceAll('_', ' '),
                                      ),
                                    ),
                                    DataCell(
                                      Text(
                                        item.metricValue.toStringAsFixed(2),
                                        style: const TextStyle(
                                          fontWeight: FontWeight.w600,
                                        ),
                                      ),
                                    ),
                                    DataCell(Text(item.unit)),
                                    DataCell(
                                      Text(
                                        item.source,
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
                      Container(
                        padding: const EdgeInsets.fromLTRB(16, 14, 16, 16),
                        decoration: const BoxDecoration(
                          border: Border(
                            top: BorderSide(color: AppColors.border),
                          ),
                        ),
                        child: compact
                            ? Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    _resultLabel(
                                      _page,
                                      pageItems.length,
                                      result.total,
                                    ),
                                    style: const TextStyle(
                                      color: AppColors.mutedForeground,
                                    ),
                                  ),
                                  const SizedBox(height: 12),
                                  AppPagination(
                                    currentPage: _page,
                                    totalPages: totalPages,
                                    onChanged: (v) => setState(() => _page = v),
                                  ),
                                ],
                              )
                            : Row(
                                children: [
                                  Expanded(
                                    child: Text(
                                      _resultLabel(
                                        _page,
                                        pageItems.length,
                                        result.total,
                                      ),
                                      style: const TextStyle(
                                        color: AppColors.mutedForeground,
                                      ),
                                    ),
                                  ),
                                  AppPagination(
                                    currentPage: _page,
                                    totalPages: totalPages,
                                    onChanged: (v) => setState(() => _page = v),
                                  ),
                                ],
                              ),
                      ),
                    ],
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  void _changeMetricType(String? value) {
    setState(() {
      _metricType = value;
      _page = 1;
    });
  }

  void _changeCowId(String value) {
    setState(() {
      _cowId = value.isEmpty ? null : value;
      _page = 1;
    });
  }

  String _formatTime(DateTime dt) {
    String p(int n, [int w = 2]) => n.toString().padLeft(w, '0');
    return '${dt.year}-${p(dt.month)}-${p(dt.day)} ${p(dt.hour)}:${p(dt.minute)}:${p(dt.second)}';
  }

  String _resultLabel(int currentPage, int currentCount, int totalCount) {
    if (totalCount == 0) return 'Showing 0 results';
    final start = ((currentPage - 1) * _pageSize) + 1;
    final end = start + currentCount - 1;
    return 'Showing $start to $end of $totalCount results';
  }
}
