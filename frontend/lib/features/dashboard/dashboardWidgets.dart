import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/models/appModels.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';

class DashboardStatCard extends StatelessWidget {
  const DashboardStatCard({
    super.key,
    required this.label,
    required this.value,
    required this.color,
  });

  final String label;
  final int value;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return SurfaceCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label, style: const TextStyle(color: AppColors.mutedForeground)),
          const SizedBox(height: 10),
          Text(
            '$value',
            style: Theme.of(
              context,
            ).textTheme.headlineMedium?.copyWith(color: color),
          ),
        ],
      ),
    );
  }
}

class DashboardCowCard extends StatelessWidget {
  const DashboardCowCard({super.key, required this.item});

  final DashboardCowItem item;

  @override
  Widget build(BuildContext context) {
    final color = item.condition.color;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: () => context.go('/cows/${item.id}'),
        child: Ink(
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: item.condition == CowCondition.normal
                  ? AppColors.border
                  : color.withValues(alpha: 0.92),
              width: item.condition == CowCondition.normal ? 1 : 1.5,
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(14),
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
                            item.name,
                            style: const TextStyle(
                              fontSize: 17,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            item.tag,
                            style: const TextStyle(
                              color: AppColors.mutedForeground,
                            ),
                          ),
                        ],
                      ),
                    ),
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Container(
                          width: 10,
                          height: 10,
                          decoration: BoxDecoration(
                            color: color,
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 8),
                        Text(
                          titleCaseEnum(item.condition.name),
                          style: TextStyle(color: color),
                        ),
                      ],
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    DashboardMetricCell(
                      icon: Icons.show_chart_rounded,
                      label: 'Temp',
                      value: item.temperature == null
                          ? '--'
                          : '${item.temperature}°C',
                    ),
                    DashboardMetricCell(
                      icon: Icons.favorite_border_rounded,
                      label: 'Heart',
                      value: item.heart_rate == null
                          ? '--'
                          : '${item.heart_rate} bpm',
                    ),
                    DashboardMetricCell(
                      icon: Icons.water_drop_outlined,
                      label: 'SpO2',
                      value: item.blood_oxygen == null
                          ? '--'
                          : '${item.blood_oxygen}%',
                    ),
                  ],
                ),
                const SizedBox(height: 10),
                if (item.alert_message != null)
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.symmetric(
                      horizontal: 11,
                      vertical: 10,
                    ),
                    decoration: BoxDecoration(
                      color: AppColors.surfaceMuted,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Icon(
                          item.condition == CowCondition.offline
                              ? Icons.wifi_off_rounded
                              : Icons.error_outline_rounded,
                          size: 18,
                          color: color,
                        ),
                        const SizedBox(width: 10),
                        Expanded(
                          child: Text(
                            _alertText(item),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(fontSize: 13),
                          ),
                        ),
                      ],
                    ),
                  )
                else
                  const SizedBox(height: 18),
                const SizedBox(height: 12),
                Text(
                  'Updated ${formatAgo(item.updated_at)}',
                  style: const TextStyle(
                    color: AppColors.mutedForeground,
                    fontSize: 13,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class DashboardMetricCell extends StatelessWidget {
  const DashboardMetricCell({
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
    return Expanded(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 17, color: AppColors.mutedForeground),
          const SizedBox(width: 6),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.mutedForeground,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  value,
                  style: const TextStyle(
                    fontWeight: FontWeight.w700,
                    fontSize: 13,
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

String _alertText(DashboardCowItem item) => item.alert_message ?? '';
