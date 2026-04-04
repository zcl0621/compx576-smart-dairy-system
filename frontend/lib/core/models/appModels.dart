// ignore_for_file: non_constant_identifier_names, constant_identifier_names

enum CowStatus { in_farm, sold, inactive }

enum CowCondition { normal, warning, critical, offline }

enum AlertSeverity { warning, critical, offline }

enum AlertStatus { active, resolved }

enum MetricRange { h24, d7, d30, all }

CowStatus parseCowStatus(String? v) => CowStatus.values.firstWhere(
      (e) => e.name == v,
      orElse: () => CowStatus.in_farm,
    );

CowCondition parseCowCondition(String? v) => CowCondition.values.firstWhere(
      (e) => e.name == v,
      orElse: () => CowCondition.normal,
    );

AlertSeverity parseAlertSeverity(String? v) =>
    AlertSeverity.values.firstWhere(
      (e) => e.name == v,
      orElse: () => AlertSeverity.warning,
    );

AlertStatus parseAlertStatus(String? v) => AlertStatus.values.firstWhere(
      (e) => e.name == v,
      orElse: () => AlertStatus.active,
    );

String metricRangeToQuery(MetricRange r) => switch (r) {
      MetricRange.h24 => '24h',
      MetricRange.d7 => '7d',
      MetricRange.d30 => '30d',
      MetricRange.all => 'all',
    };

class DashboardSummary {
  const DashboardSummary({
    required this.total_cows,
    required this.normal,
    required this.warning,
    required this.critical,
    required this.offline,
  });

  factory DashboardSummary.fromJson(Map<String, dynamic> json) {
    return DashboardSummary(
      total_cows: (json['total_cows'] as num?)?.toInt() ?? 0,
      normal: (json['normal'] as num?)?.toInt() ?? 0,
      warning: (json['warning'] as num?)?.toInt() ?? 0,
      critical: (json['critical'] as num?)?.toInt() ?? 0,
      offline: (json['offline'] as num?)?.toInt() ?? 0,
    );
  }

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

  factory DashboardCowItem.fromJson(Map<String, dynamic> json) {
    return DashboardCowItem(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      tag: json['tag'] as String? ?? '',
      condition: parseCowCondition(json['condition'] as String?),
      temperature: (json['temperature'] as num?)?.toDouble(),
      heart_rate: (json['heart_rate'] as num?)?.toDouble(),
      blood_oxygen: (json['blood_oxygen'] as num?)?.toDouble(),
      alert_message: json['alert_message'] as String?,
      updated_at: DateTime.tryParse(json['updated_at']?.toString() ?? '') ??
          DateTime.now(),
    );
  }

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
    this.milk_amount,
    this.temperature,
    this.heart_rate,
    this.blood_oxygen,
  });

  /// parse from cow/info response
  factory Cow.fromJson(Map<String, dynamic> json) {
    return Cow(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      tag: json['tag'] as String? ?? '',
      age: (json['age'] as num?)?.toInt() ?? 0,
      can_milking: json['can_milking'] as bool? ?? false,
      status: parseCowStatus(json['status'] as String?),
      condition: parseCowCondition(json['condition'] as String?),
      updated_at: DateTime.tryParse(json['updated_at']?.toString() ?? '') ??
          DateTime.now(),
      weight: (json['weight'] as num?)?.toDouble(),
      temperature: (json['temperature'] as num?)?.toDouble(),
      heart_rate: (json['heart_rate'] as num?)?.toDouble(),
      blood_oxygen: (json['blood_oxygen'] as num?)?.toDouble(),
      milk_amount: (json['milk_amount'] as num?)?.toDouble(),
    );
  }

  /// parse from cow/list item (same shape minus weight)
  factory Cow.fromListJson(Map<String, dynamic> json) {
    return Cow(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      tag: json['tag'] as String? ?? '',
      age: (json['age'] as num?)?.toInt() ?? 0,
      can_milking: json['can_milking'] as bool? ?? false,
      status: parseCowStatus(json['status'] as String?),
      condition: parseCowCondition(json['condition'] as String?),
      updated_at: DateTime.tryParse(json['updated_at']?.toString() ?? '') ??
          DateTime.now(),
      weight: null,
    );
  }

  Map<String, dynamic> toCreateJson() => {
        'name': name,
        'tag': tag,
        'age': age,
        'can_milking': can_milking,
        'status': status.name,
      };

  Map<String, dynamic> toUpdateJson() => {
        'id': id,
        'name': name,
        'tag': tag,
        'age': age,
        'can_milking': can_milking,
        'status': status.name,
      };

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

  factory AlertItem.fromJson(Map<String, dynamic> json) {
    return AlertItem(
      id: json['id'] as String? ?? '',
      cow_id: json['cow_id'] as String? ?? '',
      cow_name: json['cow_name'] as String? ?? '',
      metric_key: json['metric_key'] as String? ?? '',
      title: json['title'] as String? ?? '',
      message: json['message'] as String? ?? '',
      severity: parseAlertSeverity(json['severity'] as String?),
      status: parseAlertStatus(json['status'] as String?),
      resolved_at:
          DateTime.tryParse(json['resolved_at']?.toString() ?? ''),
      updated_at: DateTime.tryParse(json['updated_at']?.toString() ?? '') ??
          DateTime.now(),
    );
  }

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

class AlertSummary {
  const AlertSummary({
    required this.active,
    required this.warning,
    required this.critical,
    required this.offline,
  });

  factory AlertSummary.fromJson(Map<String, dynamic> json) {
    return AlertSummary(
      active: (json['active'] as num?)?.toInt() ?? 0,
      warning: (json['warning'] as num?)?.toInt() ?? 0,
      critical: (json['critical'] as num?)?.toInt() ?? 0,
      offline: (json['offline'] as num?)?.toInt() ?? 0,
    );
  }

  final int active;
  final int warning;
  final int critical;
  final int offline;
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

  factory ReportItem.fromJson(Map<String, dynamic> json) {
    // details comes as a JSON object from backend
    final rawDetails = json['details'];
    String detailsStr;
    if (rawDetails is Map) {
      detailsStr = rawDetails['note'] as String? ?? '';
    } else {
      detailsStr = rawDetails?.toString() ?? '';
    }

    return ReportItem(
      id: json['id'] as String? ?? '',
      cow_id: json['cow_id'] as String? ?? '',
      cow_name: json['cow_name'] as String? ?? '',
      period_start:
          DateTime.tryParse(json['period_start']?.toString() ?? '') ??
              DateTime.now(),
      period_end:
          DateTime.tryParse(json['period_end']?.toString() ?? '') ??
              DateTime.now(),
      summary: json['summary'] as String? ?? '',
      score: (json['score'] as num?)?.toDouble() ?? 0,
      details: detailsStr,
      created_at:
          DateTime.tryParse(json['created_at']?.toString() ?? '') ??
              DateTime.now(),
    );
  }

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

  factory UserItem.fromJson(Map<String, dynamic> json) {
    return UserItem(
      id: json['id'] as String? ?? '',
      username: json['username'] as String? ?? '',
      email: json['email'] as String? ?? '',
      created_at:
          DateTime.tryParse(json['created_at']?.toString() ?? '') ??
              DateTime.now(),
      updated_at:
          DateTime.tryParse(json['updated_at']?.toString() ?? '') ??
              DateTime.now(),
    );
  }

  final String id;
  final String username;
  final String email;
  final DateTime created_at;
  final DateTime updated_at;
}

class MetricPoint {
  const MetricPoint({required this.time, required this.value});

  factory MetricPoint.fromJson(Map<String, dynamic> json) {
    return MetricPoint(
      time: DateTime.tryParse(json['time']?.toString() ?? '') ?? DateTime.now(),
      value: (json['value'] as num?)?.toDouble() ?? 0,
    );
  }

  final DateTime time;
  final double value;
}

class MovementPoint {
  const MovementPoint({required this.time, required this.distance_m});

  factory MovementPoint.fromJson(Map<String, dynamic> json) {
    return MovementPoint(
      time: DateTime.tryParse(json['time']?.toString() ?? '') ?? DateTime.now(),
      distance_m: (json['distance_m'] as num?)?.toDouble() ?? 0,
    );
  }

  final DateTime time;
  final double distance_m;
}

/// generic paginated response wrapper
class PaginatedList<T> {
  const PaginatedList({
    required this.list,
    required this.page,
    required this.total,
    required this.totalPages,
  });

  final List<T> list;
  final int page;
  final int total;
  final int totalPages;
}

/// login response
class LoginResult {
  const LoginResult({
    required this.token,
    required this.expiresAt,
    required this.user,
  });

  factory LoginResult.fromJson(Map<String, dynamic> json) {
    return LoginResult(
      token: json['token'] as String? ?? '',
      expiresAt: DateTime.tryParse(json['expires_at']?.toString() ?? '') ??
          DateTime.now(),
      user: UserItem.fromJson(json['user'] as Map<String, dynamic>? ?? {}),
    );
  }

  final String token;
  final DateTime expiresAt;
  final UserItem user;
}

/// metric response with summary + series (for temperature, heart_rate, blood_oxygen, weight)
class StandardMetricResponse {
  const StandardMetricResponse({
    required this.current,
    required this.avg,
    required this.min,
    required this.max,
    required this.series,
  });

  factory StandardMetricResponse.fromJson(Map<String, dynamic> json) {
    final summary = json['summary'] as Map<String, dynamic>? ?? {};
    final rawSeries = json['series'] as List<dynamic>? ?? [];
    return StandardMetricResponse(
      current: (summary['current'] as num?)?.toDouble(),
      avg: (summary['avg'] as num?)?.toDouble(),
      min: (summary['min'] as num?)?.toDouble(),
      max: (summary['max'] as num?)?.toDouble(),
      series: rawSeries
          .map((e) => MetricPoint.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  final double? current;
  final double? avg;
  final double? min;
  final double? max;
  final List<MetricPoint> series;
}

/// milk amount metric response
class MilkMetricResponse {
  const MilkMetricResponse({
    required this.total,
    required this.avgPerSession,
    required this.sessionCount,
    required this.series,
  });

  factory MilkMetricResponse.fromJson(Map<String, dynamic> json) {
    final summary = json['summary'] as Map<String, dynamic>? ?? {};
    final rawSeries = json['series'] as List<dynamic>? ?? [];
    return MilkMetricResponse(
      total: (summary['total'] as num?)?.toDouble() ?? 0,
      avgPerSession: (summary['avg_per_session'] as num?)?.toDouble() ?? 0,
      sessionCount: (summary['session_count'] as num?)?.toInt() ?? 0,
      series: rawSeries
          .map((e) => MetricPoint.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  final double total;
  final double avgPerSession;
  final int sessionCount;
  final List<MetricPoint> series;
}

/// movement metric response
class MovementMetricResponse {
  const MovementMetricResponse({
    required this.distanceM,
    required this.pointCount,
    required this.series,
  });

  factory MovementMetricResponse.fromJson(Map<String, dynamic> json) {
    final summary = json['summary'] as Map<String, dynamic>? ?? {};
    final rawSeries = json['series'] as List<dynamic>? ?? [];
    return MovementMetricResponse(
      distanceM: (summary['distance_m'] as num?)?.toDouble() ?? 0,
      pointCount: (summary['point_count'] as num?)?.toInt() ?? 0,
      series: rawSeries
          .map((e) => MovementPoint.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  final double distanceM;
  final int pointCount;
  final List<MovementPoint> series;
}
