import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/mock/appMockData.dart';
import '../../core/models/appModels.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'cowWidgets.dart';

class CowsListPage extends StatefulWidget {
  const CowsListPage({super.key});

  @override
  State<CowsListPage> createState() => _CowsListPageState();
}

class _CowsListPageState extends State<CowsListPage> {
  String query = '';
  CowCondition? condition;
  CowStatus? status;
  String sortBy = 'Last Updated';
  int page = 1;

  static const pageSize = 10;

  List<CowTableRowData> get allRows => _buildCowTableRows();

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final filtered = _getRows();
    final totalPages = (filtered.length / pageSize).ceil().clamp(1, 999);
    final currentPage = page.clamp(1, totalPages);
    _syncPage(currentPage);
    final pageRows = filtered
        .skip((currentPage - 1) * pageSize)
        .take(pageSize)
        .toList();

    return SingleChildScrollView(
      child: PageSection(
        padding: EdgeInsets.fromLTRB(
          width < 700 ? 16 : 24,
          24,
          width < 700 ? 16 : 24,
          32,
        ),
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
                        'Cows',
                        style: Theme.of(context).textTheme.headlineMedium,
                      ),
                      const SizedBox(height: 6),
                      Text(
                        'Showing ${filtered.length} of ${allRows.length} cows',
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppColors.mutedForeground,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                SizedBox(
                  height: 42,
                  child: ElevatedButton.icon(
                    onPressed: () => context.go('/cows/add'),
                    style: ElevatedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                    ),
                    icon: const Icon(Icons.add_rounded, size: 18),
                    label: const Text('Add Cow'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),
            SurfaceCard(
              padding: const EdgeInsets.all(0),
              child: Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: width > 1100
                        ? Row(
                            children: [
                              Expanded(
                                child: CowSearchField(onChanged: _changeQuery),
                              ),
                              const SizedBox(width: 14),
                              Expanded(
                                child: AppSelect<CowCondition?>(
                                  prefixIcon: Icons.filter_alt_outlined,
                                  value: condition,
                                  hint: 'All Health Status',
                                  items: [
                                    (text: 'All Health Status', value: null),
                                    ...CowCondition.values.map(
                                      (item) => (
                                        text: titleCaseEnum(item.name),
                                        value: item as CowCondition?,
                                      ),
                                    ),
                                  ],
                                  onChanged: _changeCondition,
                                ),
                              ),
                              const SizedBox(width: 14),
                              Expanded(
                                child: AppSelect<CowStatus?>(
                                  prefixIcon: Icons.filter_alt_outlined,
                                  value: status,
                                  hint: 'All Farm Status',
                                  items: [
                                    (text: 'All Farm Status', value: null),
                                    ...CowStatus.values.map(
                                      (item) => (
                                        text: titleCaseEnum(item.name),
                                        value: item as CowStatus?,
                                      ),
                                    ),
                                  ],
                                  onChanged: _changeStatus,
                                ),
                              ),
                              const SizedBox(width: 14),
                              Expanded(
                                child: AppSelect<String>(
                                  value: sortBy,
                                  hint: 'Sort by Last Updated',
                                  items: const [
                                    (
                                      text: 'Sort by Last Updated',
                                      value: 'Last Updated',
                                    ),
                                    (text: 'Sort by Name', value: 'Name'),
                                  ],
                                  onChanged: _changeSortBy,
                                ),
                              ),
                            ],
                          )
                        : Wrap(
                            spacing: 14,
                            runSpacing: 14,
                            children: [
                              CowSearchField(
                                width: width > 760
                                    ? (width - 78) / 2
                                    : double.infinity,
                                onChanged: _changeQuery,
                              ),
                              AppSelect<CowCondition?>(
                                width: width > 760
                                    ? (width - 78) / 2
                                    : double.infinity,
                                prefixIcon: Icons.filter_alt_outlined,
                                value: condition,
                                hint: 'All Health Status',
                                items: [
                                  (text: 'All Health Status', value: null),
                                  ...CowCondition.values.map(
                                    (item) => (
                                      text: titleCaseEnum(item.name),
                                      value: item as CowCondition?,
                                    ),
                                  ),
                                ],
                                onChanged: _changeCondition,
                              ),
                              AppSelect<CowStatus?>(
                                width: width > 760
                                    ? (width - 78) / 2
                                    : double.infinity,
                                prefixIcon: Icons.filter_alt_outlined,
                                value: status,
                                hint: 'All Farm Status',
                                items: [
                                  (text: 'All Farm Status', value: null),
                                  ...CowStatus.values.map(
                                    (item) => (
                                      text: titleCaseEnum(item.name),
                                      value: item as CowStatus?,
                                    ),
                                  ),
                                ],
                                onChanged: _changeStatus,
                              ),
                              AppSelect<String>(
                                width: width > 760
                                    ? (width - 78) / 2
                                    : double.infinity,
                                value: sortBy,
                                hint: 'Sort by Last Updated',
                                items: const [
                                  (
                                    text: 'Sort by Last Updated',
                                    value: 'Last Updated',
                                  ),
                                  (text: 'Sort by Name', value: 'Name'),
                                ],
                                onChanged: _changeSortBy,
                              ),
                            ],
                          ),
                  ),
                  if (width < 980)
                    Padding(
                      padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
                      child: pageRows.isEmpty
                          ? const Padding(
                              padding: EdgeInsets.only(bottom: 4),
                              child: EmptyStateCard(
                                message: 'No cows match your filters',
                              ),
                            )
                          : Column(
                              children: pageRows
                                  .map(
                                    (cow) => Padding(
                                      padding: const EdgeInsets.only(
                                        bottom: 12,
                                      ),
                                      child: CowMobileCard(row: cow),
                                    ),
                                  )
                                  .toList(),
                            ),
                    )
                  else if (pageRows.isEmpty)
                    const Padding(
                      padding: EdgeInsets.fromLTRB(16, 0, 16, 16),
                      child: EmptyStateCard(
                        message: 'No cows match your filters',
                      ),
                    )
                  else
                    CowsTable(rows: pageRows),
                ],
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Text(
                    filtered.isEmpty
                        ? 'Showing 0 results'
                        : 'Showing ${(currentPage - 1) * pageSize + 1} to ${((currentPage - 1) * pageSize) + pageRows.length} of ${filtered.length} results',
                    style: const TextStyle(color: AppColors.mutedForeground),
                  ),
                ),
                if (filtered.isNotEmpty)
                  AppPagination(
                    currentPage: currentPage,
                    totalPages: totalPages,
                    onChanged: (p) => setState(() => page = p),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  List<CowTableRowData> _getRows() {
    final rows = allRows.where((item) {
      final matchQuery =
          query.isEmpty ||
          item.name.toLowerCase().contains(query.toLowerCase()) ||
          item.tag.toLowerCase().contains(query.toLowerCase());
      final matchCondition = condition == null || item.condition == condition;
      final matchStatus = status == null || item.status == status;
      return matchQuery && matchCondition && matchStatus;
    }).toList();
    rows.sort((a, b) {
      if (sortBy == 'Name') return a.name.compareTo(b.name);
      return b.updatedAt.compareTo(a.updatedAt);
    });
    return rows;
  }

  void _syncPage(int currentPage) {
    if (page == currentPage) return;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) setState(() => page = currentPage);
    });
  }

  void _changeQuery(String value) {
    setState(() {
      query = value;
      page = 1;
    });
  }

  void _changeCondition(CowCondition? value) {
    setState(() {
      condition = value;
      page = 1;
    });
  }

  void _changeStatus(CowStatus? value) {
    setState(() {
      status = value;
      page = 1;
    });
  }

  void _changeSortBy(String? value) {
    setState(() => sortBy = value ?? 'Last Updated');
  }
}

List<CowTableRowData> _buildCowTableRows() {
  return cows.map((cow) {
    return CowTableRowData(
      id: cow.id,
      tag: cow.tag,
      name: cow.name,
      canMilking: cow.can_milking,
      status: cow.status,
      condition: cow.condition,
      updatedAt: cow.updated_at,
    );
  }).toList();
}

class AddCowPage extends StatelessWidget {
  const AddCowPage({super.key});

  @override
  Widget build(BuildContext context) => const CowFormPage(
    title: 'Add Cow',
    subtitle: 'Create a new cow record',
    backLabel: 'Back to cows',
    backPath: '/cows',
  );
}

class EditCowPage extends StatelessWidget {
  const EditCowPage({super.key, required this.id});

  final String id;

  @override
  Widget build(BuildContext context) {
    final cow = getCow(id);
    if (cow == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) => context.go('/cows'));
      return const SizedBox.shrink();
    }
    return CowFormPage(
      title: 'Edit Cow',
      subtitle: 'Update the details for ${cow.name}',
      backLabel: 'Back to cow details',
      backPath: '/cows/$id',
      cow: cow,
    );
  }
}

class CowFormPage extends StatefulWidget {
  const CowFormPage({
    super.key,
    required this.title,
    required this.subtitle,
    required this.backLabel,
    required this.backPath,
    this.cow,
  });

  final String title;
  final String subtitle;
  final String backLabel;
  final String backPath;
  final Cow? cow;

  @override
  State<CowFormPage> createState() => _CowFormPageState();
}

class _CowFormPageState extends State<CowFormPage> {
  late final nameController = TextEditingController(
    text: widget.cow?.name ?? '',
  );
  late final tagController = TextEditingController(text: widget.cow?.tag ?? '');
  late final ageController = TextEditingController(
    text: widget.cow?.age.toString() ?? '',
  );
  late bool canMilking = widget.cow?.can_milking ?? true;
  late CowStatus status = widget.cow?.status ?? CowStatus.in_farm;

  @override
  void dispose() {
    nameController.dispose();
    tagController.dispose();
    ageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final formWidth = width < 1100 ? double.infinity : 680.0;

    return SingleChildScrollView(
      child: PageSection(
        padding: EdgeInsets.fromLTRB(
          width < 700 ? 16 : 24,
          24,
          width < 700 ? 16 : 24,
          32,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextButton.icon(
              onPressed: () => context.go(widget.backPath),
              icon: const Icon(Icons.arrow_back_rounded),
              label: Text(widget.backLabel),
              style: TextButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 8),
                foregroundColor: AppColors.mutedForeground,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              widget.title,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 6),
            Text(
              widget.subtitle,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: AppColors.mutedForeground,
              ),
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: formWidth,
              child: SurfaceCard(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    AppTextField(
                      label: 'Name',
                      controller: nameController,
                      hintText: 'Cow name',
                    ),
                    const SizedBox(height: 18),
                    AppTextField(
                      label: 'Tag',
                      controller: tagController,
                      hintText: 'DC-2408',
                    ),
                    const SizedBox(height: 18),
                    AppTextField(
                      label: 'Age (years)',
                      controller: ageController,
                      hintText: '4',
                      keyboardType: TextInputType.number,
                    ),
                    const SizedBox(height: 22),
                    Text(
                      'Can Milking',
                      style: Theme.of(context).textTheme.labelLarge,
                    ),
                    const SizedBox(height: 10),
                    _RadioLine<bool>(
                      value: canMilking,
                      items: const [(true, 'Yes'), (false, 'No')],
                      onChanged: (value) => setState(() => canMilking = value),
                    ),
                    const SizedBox(height: 22),
                    Text(
                      'Status',
                      style: Theme.of(context).textTheme.labelLarge,
                    ),
                    const SizedBox(height: 10),
                    _RadioLine<CowStatus>(
                      value: status,
                      items: const [
                        (CowStatus.in_farm, 'In Farm'),
                        (CowStatus.sold, 'Sold'),
                        (CowStatus.inactive, 'Inactive'),
                      ],
                      onChanged: (value) => setState(() => status = value),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            Row(
              children: [
                ElevatedButton(
                  onPressed: _saveCow,
                  child: Text(widget.cow == null ? 'Add Cow' : 'Save Changes'),
                ),
                const SizedBox(width: 14),
                TextButton(
                  onPressed: () => context.go(widget.backPath),
                  style: TextButton.styleFrom(
                    foregroundColor: AppColors.foreground,
                  ),
                  child: const Text('Cancel'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _saveCow() {
    final name = nameController.text.trim();
    final tag = tagController.text.trim();
    final age = int.tryParse(ageController.text.trim()) ?? 0;
    if (widget.cow == null) {
      cows.add(
        Cow(
          id: 'cow_${DateTime.now().millisecondsSinceEpoch}',
          name: name.isEmpty ? 'Unnamed' : name,
          tag: tag.isEmpty ? '-' : tag,
          age: age,
          can_milking: canMilking,
          status: status,
          condition: CowCondition.normal,
          updated_at: DateTime.now(),
          weight: null,
          milk_amount: null,
        ),
      );
    } else {
      final index = cows.indexWhere((item) => item.id == widget.cow!.id);
      if (index != -1) {
        final oldCow = widget.cow!;
        cows[index] = Cow(
          id: oldCow.id,
          name: name.isEmpty ? oldCow.name : name,
          tag: tag.isEmpty ? oldCow.tag : tag,
          age: age > 0 ? age : oldCow.age,
          can_milking: canMilking,
          status: status,
          condition: oldCow.condition,
          updated_at: DateTime.now(),
          weight: oldCow.weight,
          milk_amount: oldCow.milk_amount,
        );
      }
    }
    context.go(widget.backPath);
  }
}

class _RadioLine<T> extends StatelessWidget {
  const _RadioLine({
    required this.value,
    required this.items,
    required this.onChanged,
  });

  final T value;
  final List<(T, String)> items;
  final ValueChanged<T> onChanged;

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 18,
      runSpacing: 10,
      children: items.map((item) {
        final selected = item.$1 == value;
        return InkWell(
          onTap: () => onChanged(item.$1),
          borderRadius: BorderRadius.circular(999),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 18,
                height: 18,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: selected
                        ? const Color(0xFF2B7FFF)
                        : AppColors.mutedForeground,
                    width: 1.5,
                  ),
                ),
                padding: const EdgeInsets.all(3),
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: selected
                        ? const Color(0xFF2B7FFF)
                        : Colors.transparent,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Text(
                item.$2,
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }
}

class CowDetailPage extends StatefulWidget {
  const CowDetailPage({super.key, required this.id});

  final String id;

  @override
  State<CowDetailPage> createState() => _CowDetailPageState();
}

class _CowDetailPageState extends State<CowDetailPage> {
  MetricRange range = MetricRange.h24;

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final cow = getCow(widget.id);
    if (cow == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) => context.go('/cows'));
      return const SizedBox.shrink();
    }
    final cowAlerts = getCowAlerts(widget.id);
    final cowReports = getCowReports(widget.id);
    final temperatureSeries = buildMetricSeries(
      range: range,
      base: 38.6,
      step: 0.7,
    );
    final heartRateSeries = buildMetricSeries(range: range, base: 73, step: 5);
    final bloodOxygenSeries = buildMetricSeries(
      range: range,
      base: 95.5,
      step: 1.8,
    );
    final weightSeries = buildMetricSeries(
      range: range,
      base: cow.weight ?? 600,
      step: 4.3,
    );
    final milkSeries = buildMetricSeries(
      range: range,
      base: cow.milk_amount ?? 24,
      step: 2.2,
    );
    final movementSeries = buildMovementSeries(range: range);
    final charts = [
      (
        title: 'Temperature',
        values: temperatureSeries.map((item) => item.value).toList(),
        labels: cowChartLabels(
          temperatureSeries.map((item) => item.time).toList(),
          range,
        ),
        color: AppColors.primary,
        minY: 37.0,
        maxY: 41.0,
        fractionDigits: 0,
      ),
      (
        title: 'Heart Rate',
        values: heartRateSeries.map((item) => item.value).toList(),
        labels: cowChartLabels(
          heartRateSeries.map((item) => item.time).toList(),
          range,
        ),
        color: AppColors.warning,
        minY: 60.0,
        maxY: 100.0,
        fractionDigits: 0,
      ),
      (
        title: 'Blood Oxygen',
        values: bloodOxygenSeries.map((item) => item.value).toList(),
        labels: cowChartLabels(
          bloodOxygenSeries.map((item) => item.time).toList(),
          range,
        ),
        color: AppColors.normal,
        minY: 90.0,
        maxY: 100.0,
        fractionDigits: 0,
      ),
      (
        title: 'Weight',
        values: weightSeries.map((item) => item.value).toList(),
        labels: cowChartLabels(
          weightSeries.map((item) => item.time).toList(),
          range,
        ),
        color: AppColors.accent,
        minY: (cow.weight ?? 600) - 10,
        maxY: (cow.weight ?? 600) + 10,
        fractionDigits: 0,
      ),
      (
        title: 'Milk Production',
        values: milkSeries.map((item) => item.value).toList(),
        labels: cowChartLabels(
          milkSeries.map((item) => item.time).toList(),
          range,
        ),
        color: AppColors.warning,
        minY: 0.0,
        maxY: (cow.milk_amount ?? 24) + 5,
        fractionDigits: 0,
      ),
      (
        title: 'Movement Distance',
        values: movementSeries.map((item) => item.distance_m).toList(),
        labels: cowChartLabels(
          movementSeries.map((item) => item.time).toList(),
          range,
        ),
        color: AppColors.offline,
        minY: 0.0,
        maxY: switch (range) {
          MetricRange.h24 => 2.0,
          MetricRange.d7 => 10.0,
          MetricRange.d30 => 20.0,
          MetricRange.all => 200.0,
        },
        fractionDigits: 0,
      ),
    ];

    return SingleChildScrollView(
      child: PageSection(
        padding: EdgeInsets.fromLTRB(
          width < 700 ? 16 : 24,
          24,
          width < 700 ? 16 : 24,
          32,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextButton.icon(
              onPressed: () => context.go('/cows'),
              icon: const Icon(Icons.arrow_back_rounded),
              label: const Text('Back to cows'),
              style: TextButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 8),
                foregroundColor: AppColors.mutedForeground,
                textStyle: const TextStyle(fontSize: 14),
              ),
            ),
            const SizedBox(height: 16),
            CowHeaderSection(cow: cow),
            const SizedBox(height: 24),
            Text(
              'Overview',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontSize: 20),
            ),
            const SizedBox(height: 16),
            width > 1080
                ? Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(child: CowProfileCard(cow: cow)),
                      const SizedBox(width: 24),
                      Expanded(child: CurrentMetricsCard(cow: cow)),
                    ],
                  )
                : Column(
                    children: [
                      CowProfileCard(cow: cow),
                      const SizedBox(height: 16),
                      CurrentMetricsCard(cow: cow),
                    ],
                  ),
            const SizedBox(height: 20),
            Text(
              'Active Alerts',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontSize: 20),
            ),
            const SizedBox(height: 16),
            DetailCardFrame(
              child: cowAlerts.isEmpty
                  ? const EmptyStateCard(
                      message: 'No active alert for this cow',
                    )
                  : Column(
                      children: cowAlerts
                          .map(
                            (alert) => Padding(
                              padding: const EdgeInsets.only(bottom: 12),
                              child: AlertSummaryCard(alert: alert),
                            ),
                          )
                          .toList(),
                    ),
            ),
            const SizedBox(height: 20),
            Text(
              'Health Report Summary',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontSize: 20),
            ),
            const SizedBox(height: 16),
            DetailCardFrame(
              child: cowReports.isEmpty
                  ? const EmptyStateCard(message: 'No reports available')
                  : Column(
                      children: cowReports
                          .map(
                            (report) => Padding(
                              padding: const EdgeInsets.only(bottom: 14),
                              child: HealthReportSummaryCard(report: report),
                            ),
                          )
                          .toList(),
                    ),
            ),
            const SizedBox(height: 28),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Charts', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 16),
                width < 760
                    ? Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Time Period:',
                            style: TextStyle(
                              color: AppColors.mutedForeground,
                              fontSize: 14,
                              height: 20 / 14,
                              letterSpacing: -0.15,
                            ),
                          ),
                          const SizedBox(height: 12),
                          SingleChildScrollView(
                            scrollDirection: Axis.horizontal,
                            child: ChartRangeGroup(
                              selected: range,
                              onChanged: (item) => setState(() => range = item),
                            ),
                          ),
                        ],
                      )
                    : Row(
                        children: [
                          const Text(
                            'Time Period:',
                            style: TextStyle(
                              color: AppColors.mutedForeground,
                              fontSize: 14,
                              height: 20 / 14,
                              letterSpacing: -0.15,
                            ),
                          ),
                          const SizedBox(width: 12),
                          ChartRangeGroup(
                            selected: range,
                            onChanged: (item) => setState(() => range = item),
                          ),
                        ],
                      ),
                const SizedBox(height: 24),
                ListView.separated(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: charts.length,
                  separatorBuilder: (_, _) => const SizedBox(height: 24),
                  itemBuilder: (context, index) {
                    final chart = charts[index];
                    return ChartCard(
                      title: chart.title,
                      values: chart.values,
                      labels: chart.labels,
                      color: chart.color,
                      minY: chart.minY,
                      maxY: chart.maxY,
                      fractionDigits: chart.fractionDigits,
                    );
                  },
                ),
              ],
            ),
            const SizedBox(height: 28),
            Text('Map', style: Theme.of(context).textTheme.titleLarge),
            const SizedBox(height: 16),
            SurfaceCard(
              child: Container(
                height: 380,
                decoration: BoxDecoration(
                  color: AppColors.surfaceMuted,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(
                        Icons.location_on_outlined,
                        size: 48,
                        color: AppColors.mutedForeground,
                      ),
                      const SizedBox(height: 12),
                      const Text(
                        'Location: -37.7833, 175.2833',
                        style: TextStyle(color: AppColors.mutedForeground),
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        'Map integration would display here',
                        style: TextStyle(
                          fontSize: 12,
                          color: AppColors.mutedForeground,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
