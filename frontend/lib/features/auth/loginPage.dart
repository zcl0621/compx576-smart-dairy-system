import 'package:flutter/material.dart';

import '../../core/theme/appTheme.dart';
import 'authWidgets.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final emailController = TextEditingController(
    text: 'manager@smartdairy.local',
  );
  final passwordController = TextEditingController(text: 'password123');
  bool hidePassword = true;

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(16, 24, 16, 24),
          child: Center(
            child: LoginCard(
              emailController: emailController,
              passwordController: passwordController,
              hidePassword: hidePassword,
              onTogglePassword: () =>
                  setState(() => hidePassword = !hidePassword),
            ),
          ),
        ),
      ),
    );
  }
}
