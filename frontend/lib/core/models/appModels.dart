// ignore_for_file: non_constant_identifier_names, constant_identifier_names

enum CowStatus { in_farm, sold, inactive }

enum CowCondition { normal, warning, critical, offline }

enum AlertSeverity { warning, critical, offline }

enum AlertStatus { active, resolved }

enum MetricRange { h24, d7, d30, all }

class DashboardSummary {
  const DashboardSummary({
    required this.total_cows,
    required this.normal,
    required this.warning,
    required this.critical,
    required this.offline,
  });

  final int total_cows;
  final int normal;
  final int warning;
  final int critical;
  final int offline;
}

class DashboardCowItem {
  const DashboardCowItem({
    required this.id,
    required this.name,
    required this.tag,
    required this.condition,
    required this.temperature,
    required this.heart_rate,
    required this.blood_oxygen,
    required this.alert_message,
    required this.updated_at,
  });

  final String id;
  final String name;
  final String tag;
  final CowCondition condition;
  final double? temperature;
  final double? heart_rate;
  final double? blood_oxygen;
  final String? alert_message;
  final DateTime updated_at;
}

class Cow {
  const Cow({
    required this.id,
    required this.name,
    required this.tag,
    required this.age,
    required this.can_milking,
    required this.status,
    required this.condition,
    required this.updated_at,
    required this.weight,
    required this.milk_amount,
    this.temperature,
    this.heart_rate,
    this.blood_oxygen,
  });

  final String id;
  final String name;
  final String tag;
  final int age;
  final bool can_milking;
  final CowStatus status;
  final CowCondition condition;
  final DateTime updated_at;
  final double? weight;
  final double? milk_amount;
  final double? temperature;
  final double? heart_rate;
  final double? blood_oxygen;
}

class AlertItem {
  const AlertItem({
    required this.id,
    required this.cow_id,
    required this.cow_name,
    required this.metric_key,
    required this.title,
    required this.message,
    required this.severity,
    required this.status,
    required this.resolved_at,
    required this.updated_at,
  });

  final String id;
  final String cow_id;
  final String cow_name;
  final String metric_key;
  final String title;
  final String message;
  final AlertSeverity severity;
  final AlertStatus status;
  final DateTime? resolved_at;
  final DateTime updated_at;
}

class ReportItem {
  const ReportItem({
    required this.id,
    required this.cow_id,
    required this.cow_name,
    required this.period_start,
    required this.period_end,
    required this.summary,
    required this.score,
    required this.details,
    required this.created_at,
  });

  final String id;
  final String cow_id;
  final String cow_name;
  final DateTime period_start;
  final DateTime period_end;
  final String summary;
  final double score;
  final String details;
  final DateTime created_at;
}

class UserItem {
  const UserItem({
    required this.id,
    required this.username,
    required this.email,
    required this.created_at,
    required this.updated_at,
  });

  final String id;
  final String username;
  final String email;
  final DateTime created_at;
  final DateTime updated_at;
}

class MetricPoint {
  const MetricPoint({required this.time, required this.value});

  final DateTime time;
  final double value;
}

class MovementPoint {
  const MovementPoint({required this.time, required this.distance_m});

  final DateTime time;
  final double distance_m;
}
