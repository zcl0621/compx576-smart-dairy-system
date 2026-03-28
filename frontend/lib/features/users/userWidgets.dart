import 'package:flutter/material.dart';

import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';

class UserRowData {
  const UserRowData({
    required this.id,
    required this.name,
    required this.username,
    required this.email,
    required this.created,
  });

  final String id;
  final String name;
  final String username;
  final String email;
  final String created;
}

class ActionText extends StatelessWidget {
  const ActionText({
    super.key,
    required this.label,
    required this.color,
    required this.onTap,
  });

  final String label;
  final Color color;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Text(
        label,
        style: TextStyle(color: color, fontWeight: FontWeight.w600),
      ),
    );
  }
}

class ChangePasswordDialog extends StatefulWidget {
  const ChangePasswordDialog({
    super.key,
    required this.user,
    required this.onClose,
  });

  final UserRowData user;
  final VoidCallback onClose;

  @override
  State<ChangePasswordDialog> createState() => _ChangePasswordDialogState();
}

class _ChangePasswordDialogState extends State<ChangePasswordDialog> {
  final passwordController = TextEditingController();
  final confirmController = TextEditingController();
  String? error;

  @override
  void dispose() {
    passwordController.dispose();
    confirmController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return _DialogFrame(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _DialogHeader(
            icon: Icons.lock_outline_rounded,
            title: 'Change Password',
            subtitle: widget.user.name,
            onClose: widget.onClose,
            divider: true,
          ),
          const SizedBox(height: 18),
          if (error != null)
            Padding(
              padding: const EdgeInsets.only(bottom: 14),
              child: Container(
                padding: const EdgeInsets.all(14),
                decoration: BoxDecoration(
                  color: AppColors.critical.withValues(alpha: 0.08),
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(
                    color: AppColors.critical.withValues(alpha: 0.2),
                  ),
                ),
                child: Text(
                  error!,
                  style: const TextStyle(color: AppColors.foreground),
                ),
              ),
            ),
          AppTextField(
            label: 'New Password',
            controller: passwordController,
            hintText: 'Enter new password',
            obscureText: true,
          ),
          const SizedBox(height: 16),
          AppTextField(
            label: 'Confirm Password',
            controller: confirmController,
            hintText: 'Confirm new password',
            obscureText: true,
          ),
          const SizedBox(height: 16),
          Container(
            padding: const EdgeInsets.all(14),
            decoration: BoxDecoration(
              color: AppColors.surfaceMuted,
              borderRadius: BorderRadius.circular(16),
            ),
            child: const Text(
              'Password must be at least 6 characters long',
              style: TextStyle(color: AppColors.mutedForeground),
            ),
          ),
          const SizedBox(height: 18),
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              TextButton(
                onPressed: widget.onClose,
                child: const Text('Cancel'),
              ),
              const SizedBox(width: 12),
              ElevatedButton(
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 18,
                    vertical: 15,
                  ),
                ),
                onPressed: _changePassword,
                child: const Text('Change Password'),
              ),
            ],
          ),
        ],
      ),
    );
  }

  void _changePassword() {
    if (passwordController.text.length < 6) {
      setState(() => error = 'Password must be at least 6 characters');
      return;
    }
    if (passwordController.text != confirmController.text) {
      setState(() => error = 'Passwords do not match');
      return;
    }
    widget.onClose();
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Password changed successfully')),
    );
  }
}

class DeleteUserDialog extends StatelessWidget {
  const DeleteUserDialog({
    super.key,
    required this.user,
    required this.onClose,
    required this.onConfirm,
  });

  final UserRowData user;
  final VoidCallback onClose;
  final VoidCallback onConfirm;

  @override
  Widget build(BuildContext context) {
    return _DialogFrame(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _DialogHeader(
            icon: Icons.delete_outline_rounded,
            title: 'Delete User',
            subtitle: user.username,
            onClose: onClose,
          ),
          const SizedBox(height: 18),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.critical.withValues(alpha: 0.08),
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: AppColors.critical.withValues(alpha: 0.2),
              ),
            ),
            child: const Text(
              'This is still front-end only. Delete action is visual now and does not call backend.',
              style: TextStyle(height: 1.4),
            ),
          ),
          const SizedBox(height: 18),
          Text(
            'Email: ${user.email}',
            style: const TextStyle(color: AppColors.mutedForeground),
          ),
          const SizedBox(height: 20),
          Wrap(
            spacing: 12,
            runSpacing: 12,
            children: [
              OutlinedButton(onPressed: onClose, child: const Text('Cancel')),
              ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.critical,
                ),
                onPressed: () {
                  onConfirm();
                  ScaffoldMessenger.of(
                    context,
                  ).showSnackBar(const SnackBar(content: Text('User deleted')));
                },
                child: const Text('Delete User'),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _DialogFrame extends StatelessWidget {
  const _DialogFrame({required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    return Positioned.fill(
      child: GestureDetector(
        onTap: () => Navigator.of(context).focusNode.unfocus(),
        child: Material(
          color: Colors.black.withValues(alpha: 0.4),
          child: Center(
            child: SingleChildScrollView(
              padding: EdgeInsets.all(width < 700 ? 12 : 16),
              child: ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 520),
                child: SurfaceCard(
                  padding: EdgeInsets.all(width < 700 ? 18 : 24),
                  child: child,
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _DialogHeader extends StatelessWidget {
  const _DialogHeader({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onClose,
    this.divider = false,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onClose;
  final bool divider;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Row(
          children: [
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: AppColors.surfaceMuted,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: AppColors.primary),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: Theme.of(context).textTheme.titleLarge),
                  const SizedBox(height: 4),
                  Text(
                    subtitle,
                    style: const TextStyle(color: AppColors.mutedForeground),
                  ),
                ],
              ),
            ),
            IconButton(
              onPressed: onClose,
              icon: const Icon(Icons.close_rounded),
              color: AppColors.mutedForeground,
            ),
          ],
        ),
        if (divider) ...[
          const SizedBox(height: 18),
          const Divider(height: 1, color: AppColors.border),
        ],
      ],
    );
  }
}
