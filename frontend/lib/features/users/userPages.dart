import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/models/appModels.dart';
import '../../core/providers/api_provider.dart';
import '../../core/providers/data_providers.dart';
import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'userWidgets.dart';

class UsersListPage extends ConsumerStatefulWidget {
  const UsersListPage({super.key});

  @override
  ConsumerState<UsersListPage> createState() => _UsersListPageState();
}

class _UsersListPageState extends ConsumerState<UsersListPage> {
  String query = '';
  UserItem? selectedForPassword;
  UserItem? selectedForDelete;
  int page = 1;

  UserListParams get _params => UserListParams(
        page: page,
        name: query.isEmpty ? null : query,
      );

  void _refresh() {
    ref.invalidate(userListProvider);
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final listAsync = ref.watch(userListProvider(_params));

    return Stack(
      children: [
        SingleChildScrollView(
          child: PageSection(
            padding: EdgeInsets.fromLTRB(
              width < 700 ? 16 : 24,
              24,
              width < 700 ? 16 : 24,
              32,
            ),
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
                            'Users',
                            style: Theme.of(context).textTheme.headlineMedium,
                          ),
                          const SizedBox(height: 6),
                          Text(
                            'Manage system users and permissions',
                            style: Theme.of(context).textTheme.bodyMedium
                                ?.copyWith(color: AppColors.mutedForeground),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 16),
                    Padding(
                      padding: const EdgeInsets.only(top: 8),
                      child: SizedBox(
                        height: 42,
                        child: ElevatedButton.icon(
                          onPressed: () => context.go('/users/add'),
                          style: ElevatedButton.styleFrom(
                            padding: const EdgeInsets.symmetric(horizontal: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(10),
                            ),
                          ),
                          icon: const Icon(Icons.add_rounded, size: 18),
                          label: const Text('Add User'),
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 24),
                SurfaceCard(
                  padding: const EdgeInsets.all(0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.all(16),
                        child: AppTextField(
                          label: '',
                          hintText: 'Search by name or email...',
                          onChanged: _changeQuery,
                          suffixIcon: const Icon(Icons.search_rounded),
                        ),
                      ),
                      listAsync.when(
                        loading: () => const Padding(
                          padding: EdgeInsets.all(24),
                          child: LoadingStateCard(
                            message: 'Loading users',
                            lines: 5,
                          ),
                        ),
                        error: (e, _) => Padding(
                          padding: const EdgeInsets.all(24),
                          child:
                              EmptyStateCard(message: 'Failed to load: $e'),
                        ),
                        data: (result) {
                          final pageRows = result.list;
                          if (width < 860) {
                            return pageRows.isEmpty
                                ? const Padding(
                                    padding:
                                        EdgeInsets.fromLTRB(16, 0, 16, 16),
                                    child: EmptyStateCard(
                                      message: 'No users match your search',
                                    ),
                                  )
                                : Column(
                                    children: pageRows
                                        .map(
                                          (user) => Padding(
                                            padding: const EdgeInsets.fromLTRB(
                                              16,
                                              0,
                                              16,
                                              12,
                                            ),
                                            child: _UserMobileCard(
                                              user: user,
                                              onAction: (value) =>
                                                  _handleAction(
                                                      value, user, context),
                                            ),
                                          ),
                                        )
                                        .toList(),
                                  );
                          }
                          return Column(
                            children: [
                              Container(
                                color: AppColors.surfaceMuted,
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 24,
                                  vertical: 14,
                                ),
                                child: const Row(
                                  children: [
                                    Expanded(
                                      flex: 3,
                                      child: Text(
                                        'Username',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w700,
                                          color: AppColors.mutedForeground,
                                        ),
                                      ),
                                    ),
                                    Expanded(
                                      flex: 4,
                                      child: Text(
                                        'Email',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w700,
                                          color: AppColors.mutedForeground,
                                        ),
                                      ),
                                    ),
                                    Expanded(
                                      flex: 2,
                                      child: Text(
                                        'Created',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w700,
                                          color: AppColors.mutedForeground,
                                        ),
                                      ),
                                    ),
                                    Expanded(
                                      flex: 4,
                                      child: Text(
                                        'Actions',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w700,
                                          color: AppColors.mutedForeground,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              if (pageRows.isEmpty)
                                const Padding(
                                  padding:
                                      EdgeInsets.fromLTRB(16, 12, 16, 12),
                                  child: EmptyStateCard(
                                    message: 'No users match your search',
                                  ),
                                )
                              else
                                for (final user in pageRows)
                                  Container(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 24,
                                      vertical: 16,
                                    ),
                                    decoration: const BoxDecoration(
                                      border: Border(
                                        top: BorderSide(
                                            color: AppColors.border),
                                      ),
                                    ),
                                    child: Row(
                                      children: [
                                        Expanded(
                                            flex: 3,
                                            child: Text(user.username)),
                                        Expanded(
                                          flex: 4,
                                          child: Text(user.email),
                                        ),
                                        Expanded(
                                          flex: 2,
                                          child: Text(_formatDate(
                                              user.created_at)),
                                        ),
                                        Expanded(
                                          flex: 4,
                                          child: Wrap(
                                            spacing: 12,
                                            runSpacing: 8,
                                            children: [
                                              ActionText(
                                                label: 'Edit',
                                                color: AppColors.primary,
                                                onTap: () => context.go(
                                                  '/users/${user.id}/edit',
                                                ),
                                              ),
                                              ActionText(
                                                label: 'Delete',
                                                color: AppColors.critical,
                                                onTap: () => setState(
                                                  () => selectedForDelete =
                                                      user,
                                                ),
                                              ),
                                              ActionText(
                                                label: 'Change Password',
                                                color: AppColors.primary,
                                                onTap: () => setState(
                                                  () =>
                                                      selectedForPassword =
                                                          user,
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                              Padding(
                                padding: const EdgeInsets.fromLTRB(
                                  24,
                                  16,
                                  24,
                                  18,
                                ),
                                child: Row(
                                  children: [
                                    Expanded(
                                      child: Text(
                                        result.total == 0
                                            ? 'Showing 0 results'
                                            : 'Showing ${(page - 1) * 10 + 1} to ${(page - 1) * 10 + pageRows.length} of ${result.total} results',
                                        style: const TextStyle(
                                          color: AppColors.mutedForeground,
                                        ),
                                      ),
                                    ),
                                    if (result.total > 0)
                                      AppPagination(
                                        currentPage: page,
                                        totalPages: result.totalPages,
                                        onChanged: (value) =>
                                            setState(() => page = value),
                                      ),
                                  ],
                                ),
                              ),
                            ],
                          );
                        },
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        if (selectedForPassword != null)
          _ChangePasswordOverlay(
            user: selectedForPassword!,
            onClose: () => setState(() => selectedForPassword = null),
          ),
        if (selectedForDelete != null)
          _DeleteUserOverlay(
            user: selectedForDelete!,
            onClose: () => setState(() => selectedForDelete = null),
            onConfirm: () async {
              final api = ref.read(apiClientProvider);
              try {
                await deleteUser(api, selectedForDelete!.id);
                if (mounted) {
                  setState(() => selectedForDelete = null);
                  _refresh();
                }
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Delete failed: $e')),
                  );
                }
              }
            },
          ),
      ],
    );
  }

  void _changeQuery(String value) {
    setState(() {
      query = value;
      page = 1;
    });
  }

  void _handleAction(String value, UserItem user, BuildContext context) {
    if (value == 'edit') {
      context.go('/users/${user.id}/edit');
      return;
    }
    if (value == 'password') {
      setState(() => selectedForPassword = user);
      return;
    }
    if (value == 'delete') {
      setState(() => selectedForDelete = user);
    }
  }

  String _formatDate(DateTime dt) {
    return '${dt.year}/${dt.month.toString().padLeft(2, '0')}/${dt.day.toString().padLeft(2, '0')}';
  }
}

class _UserMobileCard extends StatelessWidget {
  const _UserMobileCard({required this.user, required this.onAction});

  final UserItem user;
  final void Function(String) onAction;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: BorderRadius.circular(18),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const CircleAvatar(
            backgroundColor: AppColors.surface,
            child: Icon(Icons.person_outline_rounded),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  user.username,
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
                const SizedBox(height: 4),
                Text(user.email),
                const SizedBox(height: 4),
                Text(
                  '${user.created_at.year}/${user.created_at.month}/${user.created_at.day}',
                  style: const TextStyle(color: AppColors.mutedForeground),
                ),
              ],
            ),
          ),
          PopupMenuButton<String>(
            onSelected: onAction,
            itemBuilder: (_) => const [
              PopupMenuItem(value: 'edit', child: Text('Edit')),
              PopupMenuItem(value: 'password', child: Text('Change Password')),
              PopupMenuItem(value: 'delete', child: Text('Delete')),
            ],
          ),
        ],
      ),
    );
  }
}

class _ChangePasswordOverlay extends ConsumerStatefulWidget {
  const _ChangePasswordOverlay({required this.user, required this.onClose});

  final UserItem user;
  final VoidCallback onClose;

  @override
  ConsumerState<_ChangePasswordOverlay> createState() =>
      _ChangePasswordOverlayState();
}

class _ChangePasswordOverlayState
    extends ConsumerState<_ChangePasswordOverlay> {
  final passwordController = TextEditingController();
  final confirmController = TextEditingController();
  String? error;
  bool loading = false;

  @override
  void dispose() {
    passwordController.dispose();
    confirmController.dispose();
    super.dispose();
  }

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
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
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
                            child: const Icon(Icons.lock_outline_rounded,
                                color: AppColors.primary),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text('Change Password',
                                    style: Theme.of(context)
                                        .textTheme
                                        .titleLarge),
                                const SizedBox(height: 4),
                                Text(
                                  widget.user.username,
                                  style: const TextStyle(
                                      color: AppColors.mutedForeground),
                                ),
                              ],
                            ),
                          ),
                          IconButton(
                            onPressed: widget.onClose,
                            icon: const Icon(Icons.close_rounded),
                            color: AppColors.mutedForeground,
                          ),
                        ],
                      ),
                      const SizedBox(height: 18),
                      const Divider(height: 1, color: AppColors.border),
                      const SizedBox(height: 18),
                      if (error != null)
                        Padding(
                          padding: const EdgeInsets.only(bottom: 14),
                          child: Container(
                            padding: const EdgeInsets.all(14),
                            decoration: BoxDecoration(
                              color:
                                  AppColors.critical.withValues(alpha: 0.08),
                              borderRadius: BorderRadius.circular(16),
                              border: Border.all(
                                color: AppColors.critical
                                    .withValues(alpha: 0.2),
                              ),
                            ),
                            child: Text(error!,
                                style: const TextStyle(
                                    color: AppColors.foreground)),
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
                          style:
                              TextStyle(color: AppColors.mutedForeground),
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
                            onPressed: loading ? null : _changePassword,
                            child: loading
                                ? const SizedBox(
                                    width: 18,
                                    height: 18,
                                    child: CircularProgressIndicator(
                                        strokeWidth: 2),
                                  )
                                : const Text('Change Password'),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _changePassword() async {
    if (passwordController.text.length < 6) {
      setState(() => error = 'Password must be at least 6 characters');
      return;
    }
    if (passwordController.text != confirmController.text) {
      setState(() => error = 'Passwords do not match');
      return;
    }
    setState(() {
      loading = true;
      error = null;
    });
    try {
      final api = ref.read(apiClientProvider);
      await updateUserPassword(
          api, widget.user.id, passwordController.text);
      if (mounted) {
        widget.onClose();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Password changed successfully')),
        );
      }
    } catch (e) {
      if (mounted) setState(() => error = 'Failed: $e');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }
}

class _DeleteUserOverlay extends StatelessWidget {
  const _DeleteUserOverlay({
    required this.user,
    required this.onClose,
    required this.onConfirm,
  });

  final UserItem user;
  final VoidCallback onClose;
  final VoidCallback onConfirm;

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
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
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
                            child: const Icon(Icons.delete_outline_rounded,
                                color: AppColors.primary),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text('Delete User',
                                    style: Theme.of(context)
                                        .textTheme
                                        .titleLarge),
                                const SizedBox(height: 4),
                                Text(
                                  user.username,
                                  style: const TextStyle(
                                      color: AppColors.mutedForeground),
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
                      const SizedBox(height: 18),
                      Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color:
                              AppColors.critical.withValues(alpha: 0.08),
                          borderRadius: BorderRadius.circular(16),
                          border: Border.all(
                            color:
                                AppColors.critical.withValues(alpha: 0.2),
                          ),
                        ),
                        child: const Text(
                          'This action cannot be undone. The user account will be permanently deleted.',
                          style: TextStyle(height: 1.4),
                        ),
                      ),
                      const SizedBox(height: 18),
                      Text(
                        'Email: ${user.email}',
                        style: const TextStyle(
                            color: AppColors.mutedForeground),
                      ),
                      const SizedBox(height: 20),
                      Wrap(
                        spacing: 12,
                        runSpacing: 12,
                        children: [
                          OutlinedButton(
                              onPressed: onClose,
                              child: const Text('Cancel')),
                          ElevatedButton(
                            style: ElevatedButton.styleFrom(
                              backgroundColor: AppColors.critical,
                            ),
                            onPressed: () {
                              onConfirm();
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                    content: Text('User deleted')),
                              );
                            },
                            child: const Text('Delete User'),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class AddUserPage extends StatelessWidget {
  const AddUserPage({super.key});

  @override
  Widget build(BuildContext context) => const UserFormPage(
    title: 'Add User',
    subtitle: 'Create user account details',
  );

}

class EditUserPage extends ConsumerWidget {
  const EditUserPage({super.key, required this.id});

  final String id;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userAsync = ref.watch(userInfoProvider(id));
    return userAsync.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) {
        WidgetsBinding.instance.addPostFrameCallback((_) => context.go('/users'));
        return const SizedBox.shrink();
      },
      data: (user) => UserFormPage(
        title: 'Edit User',
        subtitle: 'Update user account details',
        user: user,
      ),
    );
  }
}

class UserFormPage extends ConsumerStatefulWidget {
  const UserFormPage({
    super.key,
    required this.title,
    required this.subtitle,
    this.user,
  });

  final String title;
  final String subtitle;
  final UserItem? user;

  @override
  ConsumerState<UserFormPage> createState() => _UserFormPageState();
}

class _UserFormPageState extends ConsumerState<UserFormPage> {
  late final TextEditingController usernameController;
  late final TextEditingController emailController;
  final passwordController = TextEditingController();
  bool loading = false;
  String? error;

  @override
  void initState() {
    super.initState();
    usernameController = TextEditingController(text: widget.user?.username ?? '');
    emailController = TextEditingController(text: widget.user?.email ?? '');
  }

  @override
  void dispose() {
    usernameController.dispose();
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final formWidth = width < 1100 ? double.infinity : 680.0;
    final isEdit = widget.user != null;

    return SingleChildScrollView(
      child: PageSection(
        padding: EdgeInsets.fromLTRB(
          width < 700 ? 16 : 24,
          24,
          width < 700 ? 16 : 24,
          32,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextButton.icon(
              onPressed: () => context.go('/users'),
              icon: const Icon(Icons.arrow_back_rounded),
              label: const Text('Back to users'),
              style: TextButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 8),
                foregroundColor: AppColors.mutedForeground,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              widget.title,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 6),
            Text(
              widget.subtitle,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: AppColors.mutedForeground,
              ),
            ),
            if (error != null) ...[
              const SizedBox(height: 16),
              EmptyStateCard(message: error!),
            ],
            const SizedBox(height: 24),
            SizedBox(
              width: formWidth,
              child: SurfaceCard(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    AppTextField(
                      label: 'Username',
                      controller: usernameController,
                      hintText: 'johnm',
                    ),
                    const SizedBox(height: 18),
                    AppTextField(
                      label: 'Email',
                      controller: emailController,
                      hintText: 'john.manager@farm.com',
                    ),
                    if (!isEdit) ...[
                      const SizedBox(height: 18),
                      AppTextField(
                        label: 'Password',
                        controller: passwordController,
                        hintText: 'Enter password',
                        obscureText: true,
                      ),
                    ],
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            Row(
              children: [
                ElevatedButton(
                  onPressed: loading ? null : _saveUser,
                  child: loading
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : Text(isEdit ? 'Save Changes' : 'Create User'),
                ),
                const SizedBox(width: 14),
                TextButton(
                  onPressed: () => context.go('/users'),
                  style: TextButton.styleFrom(
                    foregroundColor: AppColors.foreground,
                  ),
                  child: const Text('Cancel'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _saveUser() async {
    final username = usernameController.text.trim();
    final email = emailController.text.trim();
    final password = passwordController.text;

    setState(() {
      loading = true;
      error = null;
    });
    try {
      final api = ref.read(apiClientProvider);
      if (widget.user == null) {
        await createUser(api, {
          'username': username,
          'email': email,
          'password': password,
        });
      } else {
        await updateUser(api, {
          'id': widget.user!.id,
          'username': username,
          'email': email,
        });
        ref.invalidate(userInfoProvider(widget.user!.id));
      }
      if (mounted) {
        ref.invalidate(userListProvider);
        context.go('/users');
      }
    } catch (e) {
      if (mounted) setState(() => error = 'Failed to save: $e');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }
}
