import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/models/appModels.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';

class CowTableRowData {
  const CowTableRowData({
    required this.id,
    required this.tag,
    required this.name,
    required this.canMilking,
    required this.status,
    required this.condition,
    required this.updatedAt,
  });

  final String id;
  final String tag;
  final String name;
  final bool canMilking;
  final CowStatus status;
  final CowCondition condition;
  final DateTime updatedAt;
}

class CowSearchField extends StatelessWidget {
  const CowSearchField({super.key, this.width, required this.onChanged});

  final double? width;
  final ValueChanged<String> onChanged;

  @override
  Widget build(BuildContext context) {
    final border = OutlineInputBorder(
      borderRadius: BorderRadius.circular(12),
      borderSide: const BorderSide(color: AppColors.strongBorder),
    );

    return SizedBox(
      width: width,
      child: TextField(
        onChanged: onChanged,
        style: const TextStyle(fontSize: 14),
        decoration: InputDecoration(
          hintText: 'Search by name or tag...',
          hintStyle: const TextStyle(color: AppColors.mutedForeground),
          prefixIcon: const Icon(
            Icons.search_rounded,
            size: 20,
            color: AppColors.mutedForeground,
          ),
          filled: true,
          fillColor: const Color(0xFFF7F5F1),
          isDense: true,
          contentPadding: const EdgeInsets.symmetric(
            horizontal: 14,
            vertical: 12,
          ),
          enabledBorder: border,
          focusedBorder: border.copyWith(
            borderSide: const BorderSide(color: AppColors.strongBorder),
          ),
          border: border,
        ),
      ),
    );
  }
}

class CowsTable extends StatelessWidget {
  const CowsTable({super.key, required this.rows});

  final List<CowTableRowData> rows;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 22, vertical: 14),
          decoration: const BoxDecoration(
            color: AppColors.surfaceMuted,
            border: Border(
              top: BorderSide(color: AppColors.border),
              bottom: BorderSide(color: AppColors.border),
            ),
          ),
          child: const Row(
            children: [
              Expanded(
                flex: 13,
                child: Text('Tag Code', style: _CowHeaderText.style),
              ),
              Expanded(
                flex: 10,
                child: Text('Name', style: _CowHeaderText.style),
              ),
              Expanded(
                flex: 12,
                child: Text('Can Milking', style: _CowHeaderText.style),
              ),
              Expanded(
                flex: 14,
                child: Text('Farm Status', style: _CowHeaderText.style),
              ),
              Expanded(
                flex: 14,
                child: Text('Health Status', style: _CowHeaderText.style),
              ),
              Expanded(
                flex: 18,
                child: Text('Last Updated', style: _CowHeaderText.style),
              ),
              Expanded(
                flex: 10,
                child: Text('Actions', style: _CowHeaderText.style),
              ),
            ],
          ),
        ),
        for (final row in rows) CowDesktopRow(row: row),
      ],
    );
  }
}

class _CowHeaderText {
  static const style = TextStyle(
    fontWeight: FontWeight.w700,
    color: AppColors.mutedForeground,
  );
}

class CowDesktopRow extends StatelessWidget {
  const CowDesktopRow({super.key, required this.row});

  final CowTableRowData row;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 22, vertical: 16),
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: AppColors.border)),
      ),
      child: Row(
        children: [
          Expanded(flex: 13, child: Text(row.tag)),
          Expanded(flex: 10, child: Text(row.name)),
          Expanded(
            flex: 12,
            child: Align(
              alignment: Alignment.centerLeft,
              child: MilkingPill(value: row.canMilking),
            ),
          ),
          Expanded(flex: 14, child: Text(cowStatusLabel(row.status))),
          Expanded(
            flex: 14,
            child: Row(
              children: [
                Container(
                  width: 8,
                  height: 8,
                  decoration: BoxDecoration(
                    color: row.condition.color,
                    shape: BoxShape.circle,
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  cowConditionLabel(row.condition),
                  style: TextStyle(color: row.condition.color),
                ),
              ],
            ),
          ),
          Expanded(flex: 18, child: Text(cowDateTimeText(row.updatedAt))),
          Expanded(
            flex: 10,
            child: Wrap(
              spacing: 10,
              children: [
                CowActionText(
                  label: 'View',
                  onTap: () => context.go('/cows/${row.id}'),
                ),
                CowActionText(
                  label: 'Edit',
                  onTap: () => context.go('/cows/${row.id}/edit'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class MilkingPill extends StatelessWidget {
  const MilkingPill({super.key, required this.value});

  final bool value;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        value ? 'Yes' : 'No',
        style: const TextStyle(
          color: AppColors.mutedForeground,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class CowActionText extends StatelessWidget {
  const CowActionText({super.key, required this.label, required this.onTap});

  final String label;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Text(
        label,
        style: const TextStyle(
          color: AppColors.primary,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class CowMobileCard extends StatelessWidget {
  const CowMobileCard({super.key, required this.row});

  final CowTableRowData row;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: BorderRadius.circular(18),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  row.name,
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
              ),
              MilkingPill(value: row.canMilking),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            row.tag,
            style: const TextStyle(color: AppColors.mutedForeground),
          ),
          const SizedBox(height: 8),
          Text('Farm status: ${cowStatusLabel(row.status)}'),
          const SizedBox(height: 4),
          Text('Health status: ${cowConditionLabel(row.condition)}'),
          const SizedBox(height: 4),
          Text('Last updated: ${cowDateTimeText(row.updatedAt)}'),
        ],
      ),
    );
  }
}

class ChartRangeGroup extends StatelessWidget {
  const ChartRangeGroup({
    super.key,
    required this.selected,
    required this.onChanged,
  });

  final MetricRange selected;
  final ValueChanged<MetricRange> onChanged;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        _ChartRangeButton(
          label: '24 Hours',
          width: 93,
          selected: selected == MetricRange.h24,
          onTap: () => onChanged(MetricRange.h24),
        ),
        const SizedBox(width: 8),
        _ChartRangeButton(
          label: '7 Days',
          width: 76,
          selected: selected == MetricRange.d7,
          onTap: () => onChanged(MetricRange.d7),
        ),
        const SizedBox(width: 8),
        _ChartRangeButton(
          label: '30 Days',
          width: 86,
          selected: selected == MetricRange.d30,
          onTap: () => onChanged(MetricRange.d30),
        ),
        const SizedBox(width: 8),
        _ChartRangeButton(
          label: 'All',
          width: 49,
          selected: selected == MetricRange.all,
          onTap: () => onChanged(MetricRange.all),
        ),
      ],
    );
  }
}

class _ChartRangeButton extends StatelessWidget {
  const _ChartRangeButton({
    required this.label,
    required this.selected,
    required this.onTap,
    required this.width,
  });

  final String label;
  final bool selected;
  final VoidCallback onTap;
  final double width;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        width: width,
        height: 36,
        decoration: BoxDecoration(
          color: selected ? AppColors.primary : AppColors.surfaceSoft,
          borderRadius: BorderRadius.circular(8),
        ),
        alignment: Alignment.center,
        child: Text(
          label,
          style: TextStyle(
            color: selected ? Colors.white : AppColors.mutedForeground,
            fontSize: 14,
            height: 20 / 14,
            fontWeight: FontWeight.w500,
            letterSpacing: -0.15,
          ),
        ),
      ),
    );
  }
}

class ChartCard extends StatelessWidget {
  const ChartCard({
    super.key,
    required this.title,
    required this.values,
    required this.labels,
    required this.color,
    required this.minY,
    required this.maxY,
    required this.fractionDigits,
  });

  final String title;
  final List<double> values;
  final List<String> labels;
  final Color color;
  final double minY;
  final double maxY;
  final int fractionDigits;

  @override
  Widget build(BuildContext context) {
    return DetailCardFrame(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: const TextStyle(
              color: AppColors.foreground,
              fontSize: 18,
              height: 28 / 18,
              fontWeight: FontWeight.w500,
              letterSpacing: -0.44,
            ),
          ),
          const SizedBox(height: 16),
          AppLineChart(
            values: values,
            labels: labels,
            color: color,
            minY: minY,
            maxY: maxY,
            fractionDigits: fractionDigits,
          ),
        ],
      ),
    );
  }
}

class CowHeaderSection extends StatelessWidget {
  const CowHeaderSection({super.key, required this.cow});

  final Cow cow;

  @override
  Widget build(BuildContext context) {
    final color = cow.condition.color;
    return SizedBox(
      height: 64,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Flexible(
                      child: Text(
                        cow.name,
                        style: const TextStyle(
                          color: AppColors.foreground,
                          fontSize: 30,
                          height: 1.2,
                          fontWeight: FontWeight.w500,
                          letterSpacing: 0.4,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Container(
                      width: 12,
                      height: 12,
                      decoration: BoxDecoration(
                        color: color,
                        shape: BoxShape.circle,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      cow.condition.name,
                      style: TextStyle(
                        color: color,
                        fontSize: 14,
                        height: 20 / 14,
                        fontWeight: FontWeight.w400,
                        letterSpacing: -0.15,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Text(
                  'Tag: ${cow.tag}',
                  style: const TextStyle(
                    color: AppColors.mutedForeground,
                    fontSize: 14,
                    height: 20 / 14,
                    letterSpacing: -0.15,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 24),
          SizedBox(
            width: 84,
            height: 44,
            child: ElevatedButton.icon(
              onPressed: () => context.go('/cows/${cow.id}/edit'),
              icon: const Icon(Icons.edit_outlined, size: 16),
              label: const Text('Edit'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
                elevation: 0,
                padding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 10,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
                textStyle: const TextStyle(
                  fontSize: 16,
                  height: 24 / 16,
                  fontWeight: FontWeight.w400,
                  letterSpacing: -0.31,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class CowProfileCard extends StatelessWidget {
  const CowProfileCard({super.key, required this.cow});

  final Cow cow;

  @override
  Widget build(BuildContext context) {
    return DetailCardFrame(
      child: SizedBox(
        height: 160,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Profile',
              style: TextStyle(
                color: AppColors.foreground,
                fontSize: 18,
                height: 28 / 18,
                fontWeight: FontWeight.w500,
                letterSpacing: -0.44,
              ),
            ),
            const SizedBox(height: 16),
            DetailInfoRow(label: 'Age', value: '${cow.age} years'),
            const SizedBox(height: 12),
            DetailInfoRow(
              label: 'Can Milking',
              value: cow.can_milking ? 'Yes' : 'No',
            ),
            const SizedBox(height: 12),
            DetailInfoRow(
              label: 'Weight',
              value: cow.weight == null
                  ? '--'
                  : '${cow.weight!.toStringAsFixed(0)} kg',
            ),
            const SizedBox(height: 12),
            DetailInfoRow(
              label: 'Last Updated',
              value: cowDetailDate(cow.updated_at),
            ),
          ],
        ),
      ),
    );
  }
}

class CurrentMetricsCard extends StatelessWidget {
  const CurrentMetricsCard({super.key, required this.cow});

  final Cow cow;

  @override
  Widget build(BuildContext context) {
    return DetailCardFrame(
      child: SizedBox(
        height: 160,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Current Metrics',
              style: TextStyle(
                color: AppColors.foreground,
                fontSize: 18,
                height: 28 / 18,
                fontWeight: FontWeight.w500,
                letterSpacing: -0.44,
              ),
            ),
            const SizedBox(height: 16),
            SizedBox(
              height: 116,
              child: GridView(
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 2,
                  crossAxisSpacing: 16,
                  mainAxisSpacing: 16,
                  mainAxisExtent: 50,
                ),
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                children: [
                  MetricBox(
                    icon: Icons.show_chart_outlined,
                    label: 'Temperature',
                    value: cow.temperature != null
                        ? '${cow.temperature!.toStringAsFixed(1)}°C'
                        : '--',
                  ),
                  MetricBox(
                    icon: Icons.favorite_border_rounded,
                    label: 'Heart Rate',
                    value: cow.heart_rate != null
                        ? '${cow.heart_rate!.toStringAsFixed(0)} bpm'
                        : '--',
                  ),
                  MetricBox(
                    icon: Icons.opacity_outlined,
                    label: 'Blood Oxygen',
                    value: cow.blood_oxygen != null
                        ? '${cow.blood_oxygen!.toStringAsFixed(0)}%'
                        : '--',
                  ),
                  MetricBox(
                    icon: Icons.local_drink_outlined,
                    label: 'Milk (24h)',
                    value: '${cow.milk_amount?.toStringAsFixed(1) ?? '--'}L',
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class AlertSummaryCard extends StatelessWidget {
  const AlertSummaryCard({super.key, required this.alert});

  final AlertItem alert;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      constraints: const BoxConstraints(minHeight: 100),
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 16),
      decoration: BoxDecoration(
        color: const Color(0xFFE8E5DF),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.only(top: 2),
            child: Icon(
              alert.severity == AlertSeverity.offline
                  ? Icons.wifi_off_rounded
                  : Icons.error_outline_rounded,
              size: 20,
              color: alert.severity.color,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  alert.title,
                  style: const TextStyle(
                    color: AppColors.foreground,
                    fontSize: 14,
                    height: 20 / 14,
                    fontWeight: FontWeight.w400,
                    letterSpacing: -0.15,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  alert.message,
                  style: const TextStyle(
                    color: AppColors.mutedForeground,
                    fontSize: 14,
                    height: 20 / 14,
                    fontWeight: FontWeight.w400,
                    letterSpacing: -0.15,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  cowDetailDate(alert.updated_at),
                  style: const TextStyle(
                    color: AppColors.mutedForeground,
                    fontSize: 12,
                    height: 16 / 12,
                    fontWeight: FontWeight.w400,
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

class HealthReportSummaryCard extends StatelessWidget {
  const HealthReportSummaryCard({super.key, required this.report});

  final ReportItem report;

  @override
  Widget build(BuildContext context) {
    final scoreColor = report.score >= 80
        ? AppColors.normal
        : report.score >= 70
        ? AppColors.warning
        : AppColors.critical;

    return Container(
      width: double.infinity,
      constraints: const BoxConstraints(minHeight: 170),
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 16),
      decoration: BoxDecoration(
        border: Border.all(color: AppColors.border),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Health Report',
                      style: TextStyle(
                        color: AppColors.foreground,
                        fontSize: 16,
                        height: 24 / 16,
                        fontWeight: FontWeight.w500,
                        letterSpacing: -0.31,
                      ),
                    ),
                    SizedBox(height: 4),
                    Text(
                      'Last 7 days',
                      style: TextStyle(
                        color: AppColors.mutedForeground,
                        fontSize: 14,
                        height: 20 / 14,
                        fontWeight: FontWeight.w400,
                        letterSpacing: -0.15,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 16),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    report.score.toStringAsFixed(0),
                    style: TextStyle(
                      color: scoreColor,
                      fontSize: 24,
                      height: 32 / 24,
                      fontWeight: FontWeight.w400,
                      letterSpacing: 0.07,
                    ),
                  ),
                  const Text(
                    'Health Score',
                    style: TextStyle(
                      color: AppColors.mutedForeground,
                      fontSize: 12,
                      height: 16 / 12,
                      fontWeight: FontWeight.w400,
                    ),
                  ),
                ],
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            report.summary,
            style: const TextStyle(
              color: AppColors.foreground,
              fontSize: 14,
              height: 20 / 14,
              fontWeight: FontWeight.w400,
              letterSpacing: -0.15,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            report.details,
            style: const TextStyle(
              color: AppColors.mutedForeground,
              fontSize: 14,
              height: 20 / 14,
              fontWeight: FontWeight.w400,
              letterSpacing: -0.15,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            'Generated: ${cowDetailDate(report.created_at)}',
            style: const TextStyle(
              color: AppColors.mutedForeground,
              fontSize: 12,
              height: 16 / 12,
              fontWeight: FontWeight.w400,
            ),
          ),
        ],
      ),
    );
  }
}

class DetailCardFrame extends StatelessWidget {
  const DetailCardFrame({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: AppColors.border),
      ),
      padding: const EdgeInsets.all(25),
      child: child,
    );
  }
}

class DetailInfoRow extends StatelessWidget {
  const DetailInfoRow({super.key, required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 20,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: const TextStyle(
              color: AppColors.mutedForeground,
              fontSize: 14,
              height: 20 / 14,
              fontWeight: FontWeight.w400,
              letterSpacing: -0.15,
            ),
          ),
          const SizedBox(width: 24),
          Flexible(
            child: Text(
              value,
              textAlign: TextAlign.right,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(
                color: AppColors.foreground,
                fontSize: 14,
                height: 20 / 14,
                fontWeight: FontWeight.w400,
                letterSpacing: -0.15,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class MetricBox extends StatelessWidget {
  const MetricBox({
    super.key,
    required this.icon,
    required this.label,
    required this.value,
  });

  final IconData icon;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 44,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          SizedBox(
            width: 32,
            height: 32,
            child: Icon(icon, size: 20, color: AppColors.mutedForeground),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    color: AppColors.mutedForeground,
                    fontSize: 12,
                    height: 16 / 12,
                    fontWeight: FontWeight.w400,
                  ),
                ),
                Text(
                  value,
                  style: const TextStyle(
                    color: AppColors.foreground,
                    fontSize: 20,
                    height: 28 / 20,
                    fontWeight: FontWeight.w400,
                    letterSpacing: -0.45,
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

String cowStatusLabel(CowStatus value) => switch (value) {
  CowStatus.in_farm => 'In Farm',
  CowStatus.sold => 'Sold',
  CowStatus.inactive => 'Inactive',
};

String cowConditionLabel(CowCondition value) => switch (value) {
  CowCondition.normal => 'Normal',
  CowCondition.warning => 'Warning',
  CowCondition.critical => 'Critical',
  CowCondition.offline => 'Offline',
};

String cowDateTimeText(DateTime value) {
  final month = value.month.toString().padLeft(2, '0');
  final day = value.day.toString().padLeft(2, '0');
  final hour = value.hour.toString().padLeft(2, '0');
  final minute = value.minute.toString().padLeft(2, '0');
  final second = value.second.toString().padLeft(2, '0');
  return '${value.year}/$month/$day $hour:$minute:$second';
}

List<String> cowChartLabels(List<DateTime> times, MetricRange range) {
  return times.map((time) {
    final month = time.month.toString().padLeft(2, '0');
    final day = time.day.toString().padLeft(2, '0');
    final hour = time.hour.toString().padLeft(2, '0');
    return switch (range) {
      MetricRange.h24 => '$hour:00',
      MetricRange.d7 => '$month/$day',
      MetricRange.d30 => '$month/$day',
      MetricRange.all => '$month/$day/${time.year}',
    };
  }).toList();
}

String cowDetailDate(DateTime value) {
  final month = value.month.toString().padLeft(2, '0');
  final day = value.day.toString().padLeft(2, '0');
  final hour = value.hour.toString().padLeft(2, '0');
  final minute = value.minute.toString().padLeft(2, '0');
  final second = value.second.toString().padLeft(2, '0');
  return '${value.year}/$month/$day $hour:$minute:$second';
}
