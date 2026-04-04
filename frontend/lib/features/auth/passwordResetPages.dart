import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/providers/auth_provider.dart';
import '../../core/services/api_client.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'authWidgets.dart';

/// shared state across the reset flow
final _resetEmailProvider = StateProvider<String>((ref) => '');
final _resetTokenProvider = StateProvider<String>((ref) => '');

class ForgotPasswordPage extends ConsumerStatefulWidget {
  const ForgotPasswordPage({super.key});

  @override
  ConsumerState<ForgotPasswordPage> createState() => _ForgotPasswordPageState();
}

class _ForgotPasswordPageState extends ConsumerState<ForgotPasswordPage> {
  final emailController = TextEditingController();
  bool success = false;
  bool loading = false;
  String? error;

  @override
  void dispose() {
    emailController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AuthFlowPage(
      step: 1,
      title: 'Reset your password',
      subtitle: 'Enter your email and we send a 6-digit code',
      children: [
        if (success)
          const StatusPanel(
            icon: Icons.mail_outline_rounded,
            color: AppColors.normal,
            title: 'Code sent successfully',
            message: 'Check your email for the latest verification code.',
          )
        else ...[
          if (error != null)
            Padding(
              padding: const EdgeInsets.only(bottom: 16),
              child: StatusPanel(
                icon: Icons.error_outline_rounded,
                color: AppColors.critical,
                title: 'Error',
                message: error!,
              ),
            ),
          AppTextField(
            label: 'Email',
            controller: emailController,
            hintText: 'admin@smartdairy.local',
            keyboardType: TextInputType.emailAddress,
          ),
          const SizedBox(height: 20),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: loading ? null : _sendCode,
              child: loading
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Text('Send code'),
            ),
          ),
        ],
        const SizedBox(height: 16),
        if (success)
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () => context.go('/verify-code'),
              child: const Text('Continue'),
            ),
          ),
      ],
    );
  }

  Future<void> _sendCode() async {
    final email = emailController.text.trim();
    if (!email.contains('@') || !email.contains('.')) {
      setState(() => error = 'Please enter a valid email address');
      return;
    }
    setState(() {
      loading = true;
      error = null;
    });
    try {
      await ref.read(authProvider.notifier).requestPasswordReset(email);
      ref.read(_resetEmailProvider.notifier).state = email;
      if (mounted) setState(() => success = true);
    } on ApiException catch (e) {
      if (mounted) setState(() => error = e.message);
    } catch (e) {
      if (mounted) setState(() => error = 'Connection failed');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }
}

class VerifyCodePage extends ConsumerStatefulWidget {
  const VerifyCodePage({super.key});

  @override
  ConsumerState<VerifyCodePage> createState() => _VerifyCodePageState();
}

class _VerifyCodePageState extends ConsumerState<VerifyCodePage> {
  final controllers = List.generate(6, (_) => TextEditingController());
  final focusNodes = List.generate(6, (_) => FocusNode());
  String? error;
  String? note;
  bool loading = false;

  @override
  void dispose() {
    for (final item in controllers) {
      item.dispose();
    }
    for (final item in focusNodes) {
      item.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final code = controllers.map((item) => item.text).join();

    return AuthFlowPage(
      step: 2,
      title: 'Verify code',
      subtitle: 'Enter the 6-digit code from your email',
      children: [
        Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: AppColors.surfaceMuted,
            borderRadius: BorderRadius.circular(16),
          ),
          child: const Row(
            children: [
              Icon(
                Icons.email_outlined,
                size: 18,
                color: AppColors.mutedForeground,
              ),
              SizedBox(width: 10),
              Expanded(
                child: Text(
                  'Sent to your email address',
                  style: TextStyle(color: AppColors.mutedForeground),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),
        if (error != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: StatusPanel(
              icon: Icons.error_outline_rounded,
              color: AppColors.critical,
              title: 'Verification failed',
              message: error!,
            ),
          ),
        if (note != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: StatusPanel(
              icon: Icons.info_outline_rounded,
              color: AppColors.warning,
              title: 'Notice',
              message: note!,
            ),
          ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: List.generate(
            6,
            (index) => SizedBox(
              width: 48,
              child: TextField(
                controller: controllers[index],
                focusNode: focusNodes[index],
                textAlign: TextAlign.center,
                keyboardType: TextInputType.number,
                maxLength: 1,
                decoration: const InputDecoration(counterText: ''),
                onChanged: (value) => _changeCode(index, value),
              ),
            ),
          ),
        ),
        const SizedBox(height: 16),
        Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: AppColors.surfaceMuted,
            borderRadius: BorderRadius.circular(16),
          ),
          child: const Center(
            child: Text(
              'Only latest code works.',
              style: TextStyle(color: AppColors.mutedForeground),
            ),
          ),
        ),
        const SizedBox(height: 20),
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: loading ? null : () => _verifyCode(code),
            child: loading
                ? const SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('Verify code'),
          ),
        ),
        const SizedBox(height: 10),
        Center(
          child: TextButton(
            onPressed: _resendCode,
            child: const Text('Resend code'),
          ),
        ),
      ],
    );
  }

  void _changeCode(int index, String value) {
    setState(() {
      error = null;
      note = null;
    });
    if (value.isNotEmpty && index < 5) {
      focusNodes[index + 1].requestFocus();
    }
  }

  Future<void> _verifyCode(String code) async {
    if (code.length != 6) {
      setState(() => error = 'Please enter the complete 6-digit code');
      return;
    }
    setState(() {
      loading = true;
      error = null;
    });
    try {
      final resetToken =
          await ref.read(authProvider.notifier).verifyResetCode(code);
      ref.read(_resetTokenProvider.notifier).state = resetToken;
      if (mounted) context.go('/reset-password');
    } on ApiException catch (e) {
      if (mounted) {
        setState(() => error = e.message);
        for (final item in controllers) {
          item.clear();
        }
        focusNodes.first.requestFocus();
      }
    } catch (e) {
      if (mounted) setState(() => error = 'Connection failed');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }

  Future<void> _resendCode() async {
    final email = ref.read(_resetEmailProvider);
    if (email.isEmpty) {
      setState(() => note = 'Go back and enter your email first');
      return;
    }
    try {
      await ref.read(authProvider.notifier).requestPasswordReset(email);
      if (mounted) setState(() => note = 'New code sent successfully');
    } catch (_) {
      if (mounted) setState(() => error = 'Failed to resend code');
    }
  }
}

class ResetPasswordPage extends ConsumerStatefulWidget {
  const ResetPasswordPage({super.key});

  @override
  ConsumerState<ResetPasswordPage> createState() => _ResetPasswordPageState();
}

class _ResetPasswordPageState extends ConsumerState<ResetPasswordPage> {
  final passwordController = TextEditingController();
  final confirmController = TextEditingController();
  bool hide1 = true;
  bool hide2 = true;
  bool loading = false;
  String? error;

  @override
  void dispose() {
    passwordController.dispose();
    confirmController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final password = passwordController.text;
    final confirm = confirmController.text;
    final hasMinLength = password.length >= 8;
    final hasLettersAndNumbers = RegExp(
      r'(?=.*[A-Za-z])(?=.*\d)',
    ).hasMatch(password);
    final matched = password.isNotEmpty && password == confirm;

    return AuthFlowPage(
      step: 3,
      title: 'Create new password',
      subtitle: 'Enter your new password',
      children: [
        if (error != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: StatusPanel(
              icon: Icons.error_outline_rounded,
              color: AppColors.critical,
              title: 'Password not ready',
              message: error!,
            ),
          ),
        AppTextField(
          label: 'New password',
          controller: passwordController,
          obscureText: hide1,
          suffixIcon: IconButton(
            onPressed: () => setState(() => hide1 = !hide1),
            icon: Icon(
              hide1 ? Icons.visibility_outlined : Icons.visibility_off_outlined,
            ),
          ),
          onChanged: (_) => setState(() => error = null),
        ),
        const SizedBox(height: 16),
        AppTextField(
          label: 'Confirm password',
          controller: confirmController,
          obscureText: hide2,
          suffixIcon: IconButton(
            onPressed: () => setState(() => hide2 = !hide2),
            icon: Icon(
              hide2 ? Icons.visibility_outlined : Icons.visibility_off_outlined,
            ),
          ),
          onChanged: (_) => setState(() => error = null),
        ),
        const SizedBox(height: 16),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.surfaceMuted,
            borderRadius: BorderRadius.circular(16),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Password must have',
                style: TextStyle(color: AppColors.mutedForeground),
              ),
              const SizedBox(height: 10),
              PasswordRule(label: 'At least 8 characters', ok: hasMinLength),
              PasswordRule(
                label: 'Use letters and numbers',
                ok: hasLettersAndNumbers,
              ),
              PasswordRule(label: 'Passwords match', ok: matched),
            ],
          ),
        ),
        const SizedBox(height: 20),
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: loading
                ? null
                : () =>
                    _updatePassword(hasMinLength, hasLettersAndNumbers, matched),
            child: loading
                ? const SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('Update password'),
          ),
        ),
      ],
    );
  }

  Future<void> _updatePassword(
    bool hasMinLength,
    bool hasLettersAndNumbers,
    bool matched,
  ) async {
    if (!hasMinLength || !hasLettersAndNumbers) {
      setState(
        () => error =
            'Password must be at least 8 characters and use letters with numbers',
      );
      return;
    }
    if (!matched) {
      setState(() => error = 'Passwords do not match');
      return;
    }
    final resetToken = ref.read(_resetTokenProvider);
    if (resetToken.isEmpty) {
      setState(() => error = 'Reset token missing. Start the flow again.');
      return;
    }
    setState(() {
      loading = true;
      error = null;
    });
    try {
      await ref.read(authProvider.notifier).confirmPasswordReset(
            resetToken,
            passwordController.text,
          );
      if (mounted) context.go('/reset-success');
    } on ApiException catch (e) {
      if (mounted) setState(() => error = e.message);
    } catch (e) {
      if (mounted) setState(() => error = 'Connection failed');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }
}

class ResetSuccessPage extends StatelessWidget {
  const ResetSuccessPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const AuthFlowPage(
      step: 4,
      title: 'Password updated',
      subtitle: 'You can sign in with the new password now',
      children: [
        Center(
          child: CircleAvatar(
            radius: 38,
            backgroundColor: Color(0x1F8A9A9A),
            child: Icon(
              Icons.check_circle_outline_rounded,
              size: 40,
              color: AppColors.normal,
            ),
          ),
        ),
        SizedBox(height: 20),
        _ResetDoneText(),
        SizedBox(height: 20),
        _BackLoginButton(),
      ],
    );
  }
}

class _ResetDoneText extends StatelessWidget {
  const _ResetDoneText();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: BorderRadius.circular(16),
      ),
      child: const Text(
        'The reset flow is complete. Next step is sign in again and continue to dashboard.',
        style: TextStyle(color: AppColors.mutedForeground),
      ),
    );
  }
}

class _BackLoginButton extends StatelessWidget {
  const _BackLoginButton();

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: () => context.go('/login'),
        child: const Text('Back to login'),
      ),
    );
  }
}
