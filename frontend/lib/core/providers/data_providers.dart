// ignore_for_file: non_constant_identifier_names

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/appModels.dart';
import '../services/api_client.dart';
import 'api_provider.dart';

// ── dashboard ──

final dashboardSummaryProvider = FutureProvider<DashboardSummary>((ref) async {
  final api = ref.watch(apiClientProvider);
  final json = await api.get('/api/dashboard/summary');
  return DashboardSummary.fromJson(json);
});

final dashboardListProvider =
    FutureProvider.family<PaginatedList<DashboardCowItem>, int>(
  (ref, page) async {
    final api = ref.watch(apiClientProvider);
    final json = await api.get('/api/dashboard/list', query: {
      'page': '$page',
      'page_size': '18',
    });
    final list = (json['list'] as List<dynamic>? ?? [])
        .map((e) => DashboardCowItem.fromJson(e as Map<String, dynamic>))
        .toList();
    return PaginatedList(
      list: list,
      page: json['page'] as int? ?? page,
      total: (json['total'] as num?)?.toInt() ?? 0,
      totalPages: (json['total_pages'] as num?)?.toInt() ?? 0,
    );
  },
);

// ── cow ──

class CowListParams {
  const CowListParams({
    this.page = 1,
    this.pageSize = 10,
    this.name,
    this.condition,
    this.status,
    this.sort = 'updated_at',
  });

  final int page;
  final int pageSize;
  final String? name;
  final String? condition;
  final String? status;
  final String sort;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is CowListParams &&
          page == other.page &&
          pageSize == other.pageSize &&
          name == other.name &&
          condition == other.condition &&
          status == other.status &&
          sort == other.sort;

  @override
  int get hashCode => Object.hash(page, pageSize, name, condition, status, sort);
}

final cowListProvider =
    FutureProvider.family<PaginatedList<Cow>, CowListParams>(
  (ref, params) async {
    final api = ref.watch(apiClientProvider);
    final query = <String, String>{
      'page': '${params.page}',
      'page_size': '${params.pageSize}',
      'sort': params.sort,
    };
    if (params.name != null && params.name!.isNotEmpty) {
      query['name'] = params.name!;
    }
    if (params.condition != null) query['condition'] = params.condition!;
    if (params.status != null) query['status'] = params.status!;

    final json = await api.get('/api/cow/list', query: query);
    final list = (json['list'] as List<dynamic>? ?? [])
        .map((e) => Cow.fromListJson(e as Map<String, dynamic>))
        .toList();
    return PaginatedList(
      list: list,
      page: json['page'] as int? ?? params.page,
      total: (json['total'] as num?)?.toInt() ?? 0,
      totalPages: (json['total_pages'] as num?)?.toInt() ?? 0,
    );
  },
);

final cowInfoProvider = FutureProvider.family<Cow, String>((ref, id) async {
  final api = ref.watch(apiClientProvider);
  final json = await api.get('/api/cow/info', query: {'id': id});
  return Cow.fromJson(json);
});

Future<void> createCow(ApiClient api, Map<String, dynamic> body) async {
  await api.post('/api/cow/create', body: body);
}

Future<void> updateCow(ApiClient api, Map<String, dynamic> body) async {
  await api.post('/api/cow/update', body: body);
}

// cow metrics

class CowMetricParams {
  const CowMetricParams({required this.cowId, required this.range});

  final String cowId;
  final MetricRange range;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is CowMetricParams &&
          cowId == other.cowId &&
          range == other.range;

  @override
  int get hashCode => Object.hash(cowId, range);
}

FutureProviderFamily<StandardMetricResponse, CowMetricParams>
    _stdMetricProvider(String endpoint) {
  return FutureProvider.family<StandardMetricResponse, CowMetricParams>(
    (ref, params) async {
      final api = ref.watch(apiClientProvider);
      final json = await api.get(endpoint, query: {
        'cow_id': params.cowId,
        'range': metricRangeToQuery(params.range),
      });
      return StandardMetricResponse.fromJson(json);
    },
  );
}

final temperatureMetricProvider =
    _stdMetricProvider('/api/cow/metric/temperature');
final heartRateMetricProvider =
    _stdMetricProvider('/api/cow/metric/heart_rate');
final bloodOxygenMetricProvider =
    _stdMetricProvider('/api/cow/metric/blood_oxygen');
final weightMetricProvider = _stdMetricProvider('/api/cow/metric/weight');

final milkMetricProvider =
    FutureProvider.family<MilkMetricResponse, CowMetricParams>(
  (ref, params) async {
    final api = ref.watch(apiClientProvider);
    final json = await api.get('/api/cow/metric/milk_amount', query: {
      'cow_id': params.cowId,
      'range': metricRangeToQuery(params.range),
    });
    return MilkMetricResponse.fromJson(json);
  },
);

final movementMetricProvider =
    FutureProvider.family<MovementMetricResponse, CowMetricParams>(
  (ref, params) async {
    final api = ref.watch(apiClientProvider);
    final json = await api.get('/api/cow/metric/movement', query: {
      'cow_id': params.cowId,
      'range': metricRangeToQuery(params.range),
    });
    return MovementMetricResponse.fromJson(json);
  },
);

// ── alert ──

final alertSummaryProvider = FutureProvider<AlertSummary>((ref) async {
  final api = ref.watch(apiClientProvider);
  final json = await api.get('/api/alert/summary');
  return AlertSummary.fromJson(json);
});

class AlertListParams {
  const AlertListParams({this.page = 1, this.severity});

  final int page;
  final String? severity;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is AlertListParams &&
          page == other.page &&
          severity == other.severity;

  @override
  int get hashCode => Object.hash(page, severity);
}

final alertListProvider =
    FutureProvider.family<PaginatedList<AlertItem>, AlertListParams>(
  (ref, params) async {
    final api = ref.watch(apiClientProvider);
    final query = <String, String>{
      'page': '${params.page}',
      'page_size': '8',
    };
    if (params.severity != null) query['severity'] = params.severity!;

    final json = await api.get('/api/alert/list', query: query);
    final list = (json['list'] as List<dynamic>? ?? [])
        .map((e) => AlertItem.fromJson(e as Map<String, dynamic>))
        .toList();
    return PaginatedList(
      list: list,
      page: json['page'] as int? ?? params.page,
      total: (json['total'] as num?)?.toInt() ?? 0,
      totalPages: (json['total_pages'] as num?)?.toInt() ?? 0,
    );
  },
);

// active alerts for top bar bell (status=active, up to 5)
final activeAlertsProvider = FutureProvider<List<AlertItem>>((ref) async {
  final api = ref.watch(apiClientProvider);
  final json = await api.get('/api/alert/list', query: {
    'status': 'active',
    'page_size': '5',
  });
  return (json['list'] as List<dynamic>? ?? [])
      .map((e) => AlertItem.fromJson(e as Map<String, dynamic>))
      .toList();
});

// cow alerts (for detail page)
final cowAlertsProvider =
    FutureProvider.family<List<AlertItem>, String>((ref, cowId) async {
  final api = ref.watch(apiClientProvider);
  final json = await api.get('/api/alert/list', query: {
    'cow_id': cowId,
    'page_size': '50',
  });
  return (json['list'] as List<dynamic>? ?? [])
      .map((e) => AlertItem.fromJson(e as Map<String, dynamic>))
      .toList();
});

// ── report ──

final reportListProvider =
    FutureProvider.family<PaginatedList<ReportItem>, int>(
  (ref, page) async {
    final api = ref.watch(apiClientProvider);
    final json = await api.get('/api/report/list', query: {
      'page': '$page',
      'page_size': '8',
    });
    final list = (json['list'] as List<dynamic>? ?? [])
        .map((e) => ReportItem.fromJson(e as Map<String, dynamic>))
        .toList();
    return PaginatedList(
      list: list,
      page: json['page'] as int? ?? page,
      total: (json['total'] as num?)?.toInt() ?? 0,
      totalPages: (json['total_pages'] as num?)?.toInt() ?? 0,
    );
  },
);

// cow reports (for detail page)
final cowReportLatestProvider =
    FutureProvider.family<ReportItem?, String>((ref, cowId) async {
  final api = ref.watch(apiClientProvider);
  try {
    final json = await api.get('/api/report/latest', query: {'cow_id': cowId});
    if (json.isEmpty) return null;
    return ReportItem.fromJson(json);
  } on ApiException catch (e) {
    if (e.statusCode == 404) return null;
    rethrow;
  }
});

// ── user ──

class UserListParams {
  const UserListParams({this.page = 1, this.name});

  final int page;
  final String? name;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is UserListParams && page == other.page && name == other.name;

  @override
  int get hashCode => Object.hash(page, name);
}

final userInfoProvider = FutureProvider.family<UserItem, String>((ref, id) async {
  final api = ref.watch(apiClientProvider);
  final json = await api.get('/api/user/info', query: {'id': id});
  return UserItem.fromJson(json);
});

final userListProvider =
    FutureProvider.family<PaginatedList<UserItem>, UserListParams>(
  (ref, params) async {
    final api = ref.watch(apiClientProvider);
    final query = <String, String>{
      'page': '${params.page}',
      'page_size': '10',
    };
    if (params.name != null && params.name!.isNotEmpty) {
      query['name'] = params.name!;
    }

    final json = await api.get('/api/user/list', query: query);
    final list = (json['list'] as List<dynamic>? ?? [])
        .map((e) => UserItem.fromJson(e as Map<String, dynamic>))
        .toList();
    return PaginatedList(
      list: list,
      page: json['page'] as int? ?? params.page,
      total: (json['total'] as num?)?.toInt() ?? 0,
      totalPages: (json['total_pages'] as num?)?.toInt() ?? 0,
    );
  },
);

Future<void> createUser(ApiClient api, Map<String, dynamic> body) async {
  await api.post('/api/user/create', body: body);
}

Future<void> updateUser(ApiClient api, Map<String, dynamic> body) async {
  await api.post('/api/user/update', body: body);
}

Future<void> updateUserPassword(ApiClient api, String id, String password) async {
  await api.post('/api/user/update_password', body: {
    'id': id,
    'password': password,
  });
}

Future<void> deleteUser(ApiClient api, String id) async {
  await api.post('/api/user/delete', body: {'id': id});
}
