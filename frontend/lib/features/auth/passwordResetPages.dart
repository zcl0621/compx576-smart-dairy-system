import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'authWidgets.dart';

class ForgotPasswordPage extends StatefulWidget {
  const ForgotPasswordPage({super.key});

  @override
  State<ForgotPasswordPage> createState() => _ForgotPasswordPageState();
}

class _ForgotPasswordPageState extends State<ForgotPasswordPage> {
  final emailController = TextEditingController();
  bool success = false;
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
                title: 'Email not valid',
                message: error!,
              ),
            ),
          AppTextField(
            label: 'Email',
            controller: emailController,
            hintText: 'manager@smartdairy.local',
            keyboardType: TextInputType.emailAddress,
          ),
          const SizedBox(height: 20),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _sendCode,
              child: const Text('Send code'),
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

  void _sendCode() {
    final email = emailController.text.trim();
    if (!email.contains('@') || !email.contains('.')) {
      setState(() => error = 'Please enter a valid email address');
      return;
    }
    setState(() {
      error = null;
      success = true;
    });
  }
}

class VerifyCodePage extends StatefulWidget {
  const VerifyCodePage({super.key});

  @override
  State<VerifyCodePage> createState() => _VerifyCodePageState();
}

class _VerifyCodePageState extends State<VerifyCodePage> {
  final controllers = List.generate(6, (_) => TextEditingController());
  final focusNodes = List.generate(6, (_) => FocusNode());
  String? error;
  String? note;

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
              'Only latest code works. Demo code is 123456.',
              style: TextStyle(color: AppColors.mutedForeground),
            ),
          ),
        ),
        const SizedBox(height: 20),
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: () => _verifyCode(code),
            child: const Text('Verify code'),
          ),
        ),
        const SizedBox(height: 10),
        Center(
          child: TextButton(
            onPressed: () =>
                setState(() => note = 'New code sent successfully'),
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

  void _verifyCode(String code) {
    if (code.length != 6) {
      setState(() => error = 'Please enter the complete 6-digit code');
      return;
    }
    if (code != '123456') {
      setState(() => error = 'Invalid or expired code. Please try again.');
      for (final item in controllers) {
        item.clear();
      }
      focusNodes.first.requestFocus();
      return;
    }
    context.go('/reset-password');
  }
}

class ResetPasswordPage extends StatefulWidget {
  const ResetPasswordPage({super.key});

  @override
  State<ResetPasswordPage> createState() => _ResetPasswordPageState();
}

class _ResetPasswordPageState extends State<ResetPasswordPage> {
  final passwordController = TextEditingController();
  final confirmController = TextEditingController();
  bool hide1 = true;
  bool hide2 = true;
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
            onPressed: () =>
                _updatePassword(hasMinLength, hasLettersAndNumbers, matched),
            child: const Text('Update password'),
          ),
        ),
      ],
    );
  }

  void _updatePassword(
    bool hasMinLength,
    bool hasLettersAndNumbers,
    bool matched,
  ) {
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
    context.go('/reset-success');
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
