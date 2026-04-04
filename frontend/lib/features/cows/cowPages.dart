import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/models/appModels.dart';
import '../../core/providers/api_provider.dart';
import '../../core/providers/data_providers.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'cowWidgets.dart';

class CowsListPage extends ConsumerStatefulWidget {
  const CowsListPage({super.key});

  @override
  ConsumerState<CowsListPage> createState() => _CowsListPageState();
}

class _CowsListPageState extends ConsumerState<CowsListPage> {
  String query = '';
  CowCondition? condition;
  CowStatus? status;
  String sortBy = 'Last Updated';
  int page = 1;

  CowListParams get _params => CowListParams(
        page: page,
        pageSize: 10,
        name: query.isEmpty ? null : query,
        condition: condition?.name,
        status: status?.name,
        sort: sortBy == 'Name' ? 'name' : 'updated_at',
      );

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final listAsync = ref.watch(cowListProvider(_params));

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
                      listAsync.when(
                        loading: () => Text(
                          'Loading...',
                          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: AppColors.mutedForeground,
                          ),
                        ),
                        error: (_, __) => const SizedBox.shrink(),
                        data: (result) => Text(
                          'Showing ${result.total} cows',
                          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: AppColors.mutedForeground,
                          ),
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
                  listAsync.when(
                    loading: () => const Padding(
                      padding: EdgeInsets.all(24),
                      child: LoadingStateCard(
                        message: 'Loading cows',
                        lines: 5,
                      ),
                    ),
                    error: (e, _) => Padding(
                      padding: const EdgeInsets.all(24),
                      child: EmptyStateCard(message: 'Failed to load: $e'),
                    ),
                    data: (result) {
                      final items = result.list;
                      final pageRows = items
                          .map(
                            (cow) => CowTableRowData(
                              id: cow.id,
                              tag: cow.tag,
                              name: cow.name,
                              canMilking: cow.can_milking,
                              status: cow.status,
                              condition: cow.condition,
                              updatedAt: cow.updated_at,
                            ),
                          )
                          .toList();

                      if (width < 980)
                        return Padding(
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
                        );
                      if (pageRows.isEmpty)
                        return const Padding(
                          padding: EdgeInsets.fromLTRB(16, 0, 16, 16),
                          child: EmptyStateCard(
                            message: 'No cows match your filters',
                          ),
                        );
                      return CowsTable(rows: pageRows);
                    },
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),
            listAsync.when(
              loading: () => const SizedBox.shrink(),
              error: (_, __) => const SizedBox.shrink(),
              data: (result) {
                final items = result.list;
                final totalPages = result.totalPages;
                return Row(
                  children: [
                    Expanded(
                      child: Text(
                        result.total == 0
                            ? 'Showing 0 results'
                            : 'Showing ${(page - 1) * 10 + 1} to ${(page - 1) * 10 + items.length} of ${result.total} results',
                        style: const TextStyle(color: AppColors.mutedForeground),
                      ),
                    ),
                    if (result.total > 0)
                      AppPagination(
                        currentPage: page,
                        totalPages: totalPages,
                        onChanged: (p) => setState(() => page = p),
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

class EditCowPage extends ConsumerWidget {
  const EditCowPage({super.key, required this.id});

  final String id;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cowAsync = ref.watch(cowInfoProvider(id));
    return cowAsync.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) {
        WidgetsBinding.instance.addPostFrameCallback((_) => context.go('/cows'));
        return const SizedBox.shrink();
      },
      data: (cow) => CowFormPage(
        title: 'Edit Cow',
        subtitle: 'Update the details for ${cow.name}',
        backLabel: 'Back to cow details',
        backPath: '/cows/$id',
        cow: cow,
      ),
    );
  }
}

class CowFormPage extends ConsumerStatefulWidget {
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
  ConsumerState<CowFormPage> createState() => _CowFormPageState();
}

class _CowFormPageState extends ConsumerState<CowFormPage> {
  late final nameController = TextEditingController(
    text: widget.cow?.name ?? '',
  );
  late final tagController = TextEditingController(text: widget.cow?.tag ?? '');
  late final ageController = TextEditingController(
    text: widget.cow?.age.toString() ?? '',
  );
  late bool canMilking = widget.cow?.can_milking ?? true;
  late CowStatus status = widget.cow?.status ?? CowStatus.in_farm;
  bool loading = false;
  String? error;

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
            if (error != null) ...[
              const SizedBox(height: 16),
              EmptyStateCard(message: error!),
            ],
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
                  onPressed: loading ? null : _saveCow,
                  child: loading
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : Text(widget.cow == null ? 'Add Cow' : 'Save Changes'),
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

  Future<void> _saveCow() async {
    final name = nameController.text.trim();
    final tag = tagController.text.trim();
    final age = int.tryParse(ageController.text.trim()) ?? 0;
    setState(() {
      loading = true;
      error = null;
    });
    try {
      final api = ref.read(apiClientProvider);
      if (widget.cow == null) {
        await createCow(api, {
          'name': name.isEmpty ? 'Unnamed' : name,
          'tag': tag.isEmpty ? '-' : tag,
          'age': age,
          'can_milking': canMilking,
          'status': status.name,
        });
      } else {
        await updateCow(api, {
          'id': widget.cow!.id,
          'name': name.isEmpty ? widget.cow!.name : name,
          'tag': tag.isEmpty ? widget.cow!.tag : tag,
          'age': age > 0 ? age : widget.cow!.age,
          'can_milking': canMilking,
          'status': status.name,
        });
        ref.invalidate(cowInfoProvider(widget.cow!.id));
      }
      ref.invalidate(cowListProvider);
      if (mounted) context.go(widget.backPath);
    } catch (e) {
      if (mounted) setState(() => error = 'Failed to save: $e');
    } finally {
      if (mounted) setState(() => loading = false);
    }
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

class CowDetailPage extends ConsumerStatefulWidget {
  const CowDetailPage({super.key, required this.id});

  final String id;

  @override
  ConsumerState<CowDetailPage> createState() => _CowDetailPageState();
}

class _CowDetailPageState extends ConsumerState<CowDetailPage> {
  MetricRange range = MetricRange.h24;

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final cowAsync = ref.watch(cowInfoProvider(widget.id));

    return cowAsync.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) {
        WidgetsBinding.instance.addPostFrameCallback((_) => context.go('/cows'));
        return const SizedBox.shrink();
      },
      data: (cow) => _buildDetail(context, cow, width),
    );
  }

  Widget _buildDetail(BuildContext context, Cow cow, double width) {
    final metricParams = CowMetricParams(cowId: widget.id, range: range);
    final tempAsync = ref.watch(temperatureMetricProvider(metricParams));
    final hrAsync = ref.watch(heartRateMetricProvider(metricParams));
    final boAsync = ref.watch(bloodOxygenMetricProvider(metricParams));
    final weightAsync = ref.watch(weightMetricProvider(metricParams));
    final milkAsync = ref.watch(milkMetricProvider(metricParams));
    final moveAsync = ref.watch(movementMetricProvider(metricParams));
    final alertsAsync = ref.watch(cowAlertsProvider(widget.id));
    final reportAsync = ref.watch(cowReportLatestProvider(widget.id));

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
            alertsAsync.when(
              loading: () => const LoadingStateCard(
                message: 'Loading alerts',
                lines: 2,
              ),
              error: (e, _) => EmptyStateCard(message: 'Failed to load: $e'),
              data: (cowAlerts) => DetailCardFrame(
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
            ),
            const SizedBox(height: 20),
            Text(
              'Health Report Summary',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontSize: 20),
            ),
            const SizedBox(height: 16),
            reportAsync.when(
              loading: () => const LoadingStateCard(
                message: 'Loading report',
                lines: 3,
              ),
              error: (e, _) => EmptyStateCard(message: 'Failed to load: $e'),
              data: (report) => DetailCardFrame(
                child: report == null
                    ? const EmptyStateCard(message: 'No reports available')
                    : Padding(
                        padding: const EdgeInsets.only(bottom: 14),
                        child: HealthReportSummaryCard(report: report),
                      ),
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
                _MetricChart(
                  title: 'Temperature',
                  asyncValue: tempAsync,
                  color: AppColors.primary,
                  minY: 37.0,
                  maxY: 41.0,
                  range: range,
                ),
                const SizedBox(height: 24),
                _MetricChart(
                  title: 'Heart Rate',
                  asyncValue: hrAsync,
                  color: AppColors.warning,
                  minY: 60.0,
                  maxY: 100.0,
                  range: range,
                ),
                const SizedBox(height: 24),
                _MetricChart(
                  title: 'Blood Oxygen',
                  asyncValue: boAsync,
                  color: AppColors.normal,
                  minY: 90.0,
                  maxY: 100.0,
                  range: range,
                ),
                const SizedBox(height: 24),
                _MetricChart(
                  title: 'Weight',
                  asyncValue: weightAsync,
                  color: AppColors.accent,
                  minY: (cow.weight ?? 600) - 10,
                  maxY: (cow.weight ?? 600) + 10,
                  range: range,
                ),
                const SizedBox(height: 24),
                _MilkChart(
                  asyncValue: milkAsync,
                  range: range,
                  milkAmount: cow.milk_amount,
                ),
                const SizedBox(height: 24),
                _MovementChart(asyncValue: moveAsync, range: range),
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

class _MetricChart extends StatelessWidget {
  const _MetricChart({
    required this.title,
    required this.asyncValue,
    required this.color,
    required this.minY,
    required this.maxY,
    required this.range,
  });

  final String title;
  final AsyncValue<StandardMetricResponse> asyncValue;
  final Color color;
  final double minY;
  final double maxY;
  final MetricRange range;

  @override
  Widget build(BuildContext context) {
    return asyncValue.when(
      loading: () => LoadingStateCard(message: 'Loading $title', lines: 3),
      error: (e, _) => EmptyStateCard(message: 'Failed to load $title'),
      data: (data) => ChartCard(
        title: title,
        values: data.series.map((p) => p.value).toList(),
        labels: cowChartLabels(
          data.series.map((p) => p.time).toList(),
          range,
        ),
        color: color,
        minY: minY,
        maxY: maxY,
        fractionDigits: 0,
      ),
    );
  }
}

class _MilkChart extends StatelessWidget {
  const _MilkChart({
    required this.asyncValue,
    required this.range,
    this.milkAmount,
  });

  final AsyncValue<MilkMetricResponse> asyncValue;
  final MetricRange range;
  final double? milkAmount;

  @override
  Widget build(BuildContext context) {
    return asyncValue.when(
      loading: () => const LoadingStateCard(
        message: 'Loading Milk Production',
        lines: 3,
      ),
      error: (e, _) =>
          const EmptyStateCard(message: 'Failed to load Milk Production'),
      data: (data) => ChartCard(
        title: 'Milk Production',
        values: data.series.map((p) => p.value).toList(),
        labels: cowChartLabels(
          data.series.map((p) => p.time).toList(),
          range,
        ),
        color: AppColors.warning,
        fractionDigits: 0,
      ),
    );
  }
}

class _MovementChart extends StatelessWidget {
  const _MovementChart({required this.asyncValue, required this.range});

  final AsyncValue<MovementMetricResponse> asyncValue;
  final MetricRange range;

  @override
  Widget build(BuildContext context) {
    return asyncValue.when(
      loading: () => const LoadingStateCard(
        message: 'Loading Movement Distance',
        lines: 3,
      ),
      error: (e, _) =>
          const EmptyStateCard(message: 'Failed to load Movement Distance'),
      data: (data) => ChartCard(
        title: 'Movement Distance (m)',
        values: data.series.map((p) => p.distance_m).toList(),
        labels: cowChartLabels(
          data.series.map((p) => p.time).toList(),
          range,
        ),
        color: AppColors.offline,
        minY: 0.0,
        fractionDigits: 0,
      ),
    );
  }
}

