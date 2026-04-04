import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../services/api_client.dart';

/// single ApiClient instance shared across the app
final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient();
});
