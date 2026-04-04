import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/appModels.dart';
import '../services/api_client.dart';
import 'api_provider.dart';

/// token + user info
class AuthState {
  const AuthState({this.token, this.user});

  final String? token;
  final UserItem? user;

  bool get isLoggedIn => token != null && token!.isNotEmpty;
}

class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier(this._api) : super(const AuthState()) {
    _api.onUnauthorized = _handleUnauthorized;
  }

  final ApiClient _api;

  Future<void> login(String email, String password) async {
    final json = await _api.post('/api/auth/login', body: {
      'email': email,
      'password': password,
    });
    final result = LoginResult.fromJson(json);
    _api.setToken(result.token);
    state = AuthState(token: result.token, user: result.user);
  }

  /// restore session from saved token
  Future<bool> tryRefresh() async {
    final saved = await _api.loadPersistedToken();
    if (saved == null || saved.isEmpty) return false;
    try {
      final json = await _api.post('/api/auth/refresh');
      final result = LoginResult.fromJson(json);
      _api.setToken(result.token);
      state = AuthState(token: result.token, user: result.user);
      return true;
    } catch (_) {
      _api.setToken(null);
      state = const AuthState();
      return false;
    }
  }

  Future<void> requestPasswordReset(String email) async {
    await _api.post('/api/auth/password-reset/request', body: {
      'email': email,
    });
  }

  Future<String> verifyResetCode(String code) async {
    final json = await _api.post('/api/auth/password-reset/verify', body: {
      'code': code,
    });
    return json['reset_token'] as String? ?? '';
  }

  Future<void> confirmPasswordReset(String resetToken, String newPassword) async {
    await _api.post('/api/auth/password-reset/confirm', body: {
      'reset_token': resetToken,
      'new_password': newPassword,
    });
  }

  void logout() {
    _api.setToken(null);
    state = const AuthState();
  }

  void _handleUnauthorized() {
    debugPrint('auth: 401, clearing token');
    _api.setToken(null);
    state = const AuthState();
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final api = ref.watch(apiClientProvider);
  return AuthNotifier(api);
});
