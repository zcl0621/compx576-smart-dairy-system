import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/appTheme.dart';
import '../../shared/appWidgets.dart';
import 'userWidgets.dart';

// mutable list so add/edit/delete persist across pages
var kUserRows = [
  UserRowData(
    id: 'row_001',
    name: 'John Manager',
    username: 'johnm',
    email: 'john.manager@farm.com',
    created: '2025/1/15',
  ),
  UserRowData(
    id: 'row_002',
    name: 'Sarah Smith',
    username: 'sarahs',
    email: 'sarah.smith@farm.com',
    created: '2025/2/10',
  ),
  UserRowData(
    id: 'row_003',
    name: 'Mike Johnson',
    username: 'mikej',
    email: 'mike.johnson@farm.com',
    created: '2025/3/1',
  ),
  UserRowData(
    id: 'row_004',
    name: 'Lisa Chen',
    username: 'lisac',
    email: 'lisa.chen@farm.com',
    created: '2025/3/5',
  ),
  UserRowData(
    id: 'row_005',
    name: 'David Brown',
    username: 'davidb',
    email: 'david.brown@farm.com',
    created: '2024/12/20',
  ),
  UserRowData(
    id: 'row_006',
    name: 'Emily White',
    username: 'emilyw',
    email: 'emily.white@farm.com',
    created: '2025/1/8',
  ),
  UserRowData(
    id: 'row_007',
    name: 'Robert Taylor',
    username: 'robertt',
    email: 'robert.taylor@farm.com',
    created: '2025/2/18',
  ),
  UserRowData(
    id: 'row_008',
    name: 'Anna Garcia',
    username: 'annag',
    email: 'anna.garcia@farm.com',
    created: '2024/11/15',
  ),
  UserRowData(
    id: 'row_009',
    name: 'James Wilson',
    username: 'jamesw',
    email: 'james.wilson@farm.com',
    created: '2025/2/25',
  ),
  UserRowData(
    id: 'row_010',
    name: 'Maria Rodriguez',
    username: 'mariar',
    email: 'maria.rodriguez@farm.com',
    created: '2025/1/22',
  ),
  UserRowData(
    id: 'row_011',
    name: 'Daniel Lee',
    username: 'daniell',
    email: 'daniel.lee@farm.com',
    created: '2025/2/4',
  ),
  UserRowData(
    id: 'row_012',
    name: 'Olivia Martin',
    username: 'oliviam',
    email: 'olivia.martin@farm.com',
    created: '2025/3/2',
  ),
  UserRowData(
    id: 'row_013',
    name: 'Thomas Clark',
    username: 'thomasc',
    email: 'thomas.clark@farm.com',
    created: '2025/1/30',
  ),
  UserRowData(
    id: 'row_014',
    name: 'Grace Hall',
    username: 'graceh',
    email: 'grace.hall@farm.com',
    created: '2024/12/29',
  ),
  UserRowData(
    id: 'row_015',
    name: 'Henry Walker',
    username: 'henryw',
    email: 'henry.walker@farm.com',
    created: '2025/2/12',
  ),
];

class UsersListPage extends StatefulWidget {
  const UsersListPage({super.key});

  @override
  State<UsersListPage> createState() => _UsersListPageState();
}

class _UsersListPageState extends State<UsersListPage> {
  String query = '';
  UserRowData? selectedForPassword;
  UserRowData? selectedForDelete;
  int page = 1;

  static const pageSize = 10;

  List<UserRowData> get allRows => kUserRows;

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final filtered = _getRows();
    final totalPages = (filtered.length / pageSize).ceil().clamp(1, 999);
    final currentPage = page.clamp(1, totalPages);
    _syncPage(currentPage);
    final pageRows = filtered
        .skip((currentPage - 1) * pageSize)
        .take(pageSize)
        .toList();

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
                          hintText: 'Search by name, username, or email...',
                          onChanged: _changeQuery,
                          suffixIcon: const Icon(Icons.search_rounded),
                        ),
                      ),
                      if (width < 860)
                        pageRows.isEmpty
                            ? const Padding(
                                padding: EdgeInsets.fromLTRB(16, 0, 16, 16),
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
                                        child: Container(
                                          padding: const EdgeInsets.all(16),
                                          decoration: BoxDecoration(
                                            color: AppColors.surfaceMuted,
                                            borderRadius: BorderRadius.circular(
                                              18,
                                            ),
                                          ),
                                          child: Row(
                                            crossAxisAlignment:
                                                CrossAxisAlignment.start,
                                            children: [
                                              const CircleAvatar(
                                                backgroundColor:
                                                    AppColors.surface,
                                                child: Icon(
                                                  Icons.person_outline_rounded,
                                                ),
                                              ),
                                              const SizedBox(width: 12),
                                              Expanded(
                                                child: Column(
                                                  crossAxisAlignment:
                                                      CrossAxisAlignment.start,
                                                  children: [
                                                    Text(
                                                      user.name,
                                                      style: const TextStyle(
                                                        fontWeight:
                                                            FontWeight.w700,
                                                      ),
                                                    ),
                                                    const SizedBox(height: 4),
                                                    Text(
                                                      '${user.username} · ${user.email}',
                                                    ),
                                                    const SizedBox(height: 4),
                                                    Text(
                                                      user.created,
                                                      style: const TextStyle(
                                                        color: AppColors
                                                            .mutedForeground,
                                                      ),
                                                    ),
                                                  ],
                                                ),
                                              ),
                                              PopupMenuButton<String>(
                                                onSelected: (value) =>
                                                    _changeUserAction(
                                                      value,
                                                      user,
                                                      context,
                                                    ),
                                                itemBuilder: (_) => const [
                                                  PopupMenuItem(
                                                    value: 'edit',
                                                    child: Text('Edit'),
                                                  ),
                                                  PopupMenuItem(
                                                    value: 'password',
                                                    child: Text(
                                                      'Change Password',
                                                    ),
                                                  ),
                                                  PopupMenuItem(
                                                    value: 'delete',
                                                    child: Text('Delete'),
                                                  ),
                                                ],
                                              ),
                                            ],
                                          ),
                                        ),
                                      ),
                                    )
                                    .toList(),
                              )
                      else
                        Column(
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
                                      'Name',
                                      style: TextStyle(
                                        fontWeight: FontWeight.w700,
                                        color: AppColors.mutedForeground,
                                      ),
                                    ),
                                  ),
                                  Expanded(
                                    flex: 2,
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
                                padding: EdgeInsets.fromLTRB(16, 12, 16, 12),
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
                                      top: BorderSide(color: AppColors.border),
                                    ),
                                  ),
                                  child: Row(
                                    children: [
                                      Expanded(flex: 3, child: Text(user.name)),
                                      Expanded(
                                        flex: 2,
                                        child: Text(user.username),
                                      ),
                                      Expanded(
                                        flex: 4,
                                        child: Text(user.email),
                                      ),
                                      Expanded(
                                        flex: 2,
                                        child: Text(user.created),
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
                                                () => selectedForDelete = user,
                                              ),
                                            ),
                                            ActionText(
                                              label: 'Change Password',
                                              color: AppColors.primary,
                                              onTap: () => setState(
                                                () =>
                                                    selectedForPassword = user,
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
                                      filtered.isEmpty
                                          ? 'Showing 0 results'
                                          : 'Showing ${(currentPage - 1) * pageSize + 1} to ${((currentPage - 1) * pageSize) + pageRows.length} of ${filtered.length} results',
                                      style: const TextStyle(
                                        color: AppColors.mutedForeground,
                                      ),
                                    ),
                                  ),
                                  if (filtered.isNotEmpty)
                                    AppPagination(
                                      currentPage: currentPage,
                                      totalPages: totalPages,
                                      onChanged: (value) =>
                                          setState(() => page = value),
                                    ),
                                ],
                              ),
                            ),
                          ],
                        ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        if (selectedForPassword != null)
          ChangePasswordDialog(
            user: selectedForPassword!,
            onClose: () => setState(() => selectedForPassword = null),
          ),
        if (selectedForDelete != null)
          DeleteUserDialog(
            user: selectedForDelete!,
            onClose: () => setState(() => selectedForDelete = null),
            onConfirm: () => setState(() {
              allRows.remove(selectedForDelete);
              selectedForDelete = null;
            }),
          ),
      ],
    );
  }

  List<UserRowData> _getRows() {
    return allRows.where((item) {
      final q = query.toLowerCase();
      return q.isEmpty ||
          item.name.toLowerCase().contains(q) ||
          item.username.toLowerCase().contains(q) ||
          item.email.toLowerCase().contains(q);
    }).toList();
  }

  void _syncPage(int currentPage) {
    if (page == currentPage) return;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) setState(() => page = currentPage);
    });
  }

  void _changeQuery(String value) {
    setState(() {
      query = value;
      page = 1;
    });
  }

  void _changeUserAction(String value, UserRowData user, BuildContext context) {
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
}

class AddUserPage extends StatelessWidget {
  const AddUserPage({super.key});

  @override
  Widget build(BuildContext context) => const UserFormPage(
    title: 'Add User',
    subtitle: 'Create user account details',
    includePassword: true,
  );
}

class EditUserPage extends StatelessWidget {
  const EditUserPage({super.key, required this.id});

  final String id;

  @override
  Widget build(BuildContext context) {
    final matches = kUserRows.where((item) => item.id == id);
    if (matches.isEmpty) {
      WidgetsBinding.instance.addPostFrameCallback((_) => context.go('/users'));
      return const SizedBox.shrink();
    }
    final user = matches.first;
    return UserFormPage(
      title: 'Edit User',
      subtitle: 'Update user account details',
      userId: user.id,
      name: user.name,
      username: user.username,
      email: user.email,
      includePassword: true,
      passwordHint: 'Leave blank to keep current password',
    );
  }
}

class UserFormPage extends StatefulWidget {
  const UserFormPage({
    super.key,
    required this.title,
    required this.subtitle,
    this.userId,
    this.name = '',
    this.username = '',
    this.email = '',
    this.includePassword = false,
    this.passwordHint,
  });

  final String title;
  final String subtitle;
  final String? userId;
  final String name;
  final String username;
  final String email;
  final bool includePassword;
  final String? passwordHint;

  @override
  State<UserFormPage> createState() => _UserFormPageState();
}

class _UserFormPageState extends State<UserFormPage> {
  late final nameController = TextEditingController(text: widget.name);
  late final usernameController = TextEditingController(text: widget.username);
  late final emailController = TextEditingController(text: widget.email);
  late final passwordController = TextEditingController();

  @override
  void dispose() {
    nameController.dispose();
    usernameController.dispose();
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final formWidth = width < 1100 ? double.infinity : 680.0;

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
            const SizedBox(height: 24),
            SizedBox(
              width: formWidth,
              child: SurfaceCard(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    AppTextField(
                      label: 'Name',
                      controller: nameController,
                      hintText: 'John Manager',
                    ),
                    const SizedBox(height: 18),
                    AppTextField(
                      label: 'Username',
                      controller: usernameController,
                      hintText: 'johnm',
                    ),
                    if (widget.includePassword) ...[
                      const SizedBox(height: 18),
                      AppTextField(
                        label: 'Password',
                        controller: passwordController,
                        hintText: '',
                        obscureText: true,
                      ),
                      if (widget.passwordHint != null) ...[
                        const SizedBox(height: 6),
                        Text(
                          widget.passwordHint!,
                          style: const TextStyle(
                            color: AppColors.mutedForeground,
                            fontSize: 13,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ],
                    const SizedBox(height: 18),
                    AppTextField(
                      label: 'Email',
                      controller: emailController,
                      hintText: 'john.manager@farm.com',
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            Row(
              children: [
                ElevatedButton(
                  onPressed: _saveUser,
                  child: Text(
                    widget.title == 'Add User' ? 'Create User' : 'Save Changes',
                  ),
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

  void _saveUser() {
    final name = nameController.text.trim();
    final username = usernameController.text.trim();
    final email = emailController.text.trim();
    if (widget.userId == null) {
      kUserRows.add(
        UserRowData(
          id: 'row_${DateTime.now().millisecondsSinceEpoch}',
          name: name.isEmpty ? 'New User' : name,
          username: username.isEmpty ? 'user' : username,
          email: email.isEmpty ? '' : email,
          created: DateTime.now().toString().substring(0, 10),
        ),
      );
    } else {
      final index = kUserRows.indexWhere((item) => item.id == widget.userId);
      if (index != -1) {
        final oldUser = kUserRows[index];
        kUserRows[index] = UserRowData(
          id: oldUser.id,
          name: name.isEmpty ? oldUser.name : name,
          username: username.isEmpty ? oldUser.username : username,
          email: email.isEmpty ? oldUser.email : email,
          created: oldUser.created,
        );
      }
    }
    context.go('/users');
  }
}
