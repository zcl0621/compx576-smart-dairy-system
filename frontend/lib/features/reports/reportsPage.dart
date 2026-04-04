import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/models/appModels.dart';
import '../../core/providers/data_providers.dart';
import '../../shared/appWidgets.dart';
import 'reportWidgets.dart';

class ReportsPage extends ConsumerStatefulWidget {
  const ReportsPage({super.key});

  @override
  ConsumerState<ReportsPage> createState() => _ReportsPageState();
}

class _ReportsPageState extends ConsumerState<ReportsPage> {
  int page = 1;
  ReportItem? selected;

  @override
  Widget build(BuildContext context) {
    final compact = MediaQuery.sizeOf(context).width < 1100;
    final listAsync = ref.watch(reportListProvider(page));

    return listAsync.when(
      loading: () => const PageSection(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            PageIntro(
              title: 'Reports',
              subtitle: 'View and analyze health reports across your herd',
            ),
            SizedBox(height: 20),
            LoadingStateCard(message: 'Loading reports', lines: 6),
          ],
        ),
      ),
      error: (e, _) => PageSection(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const PageIntro(
              title: 'Reports',
              subtitle: 'View and analyze health reports across your herd',
            ),
            const SizedBox(height: 20),
            EmptyStateCard(message: 'Failed to load: $e'),
          ],
        ),
      ),
      data: (result) {
        final allReports = result.list;
        final totalPages = result.totalPages;

        if (allReports.isEmpty) {
          return const PageSection(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                PageIntro(
                  title: 'Reports',
                  subtitle: 'View and analyze health reports across your herd',
                ),
                SizedBox(height: 20),
                EmptyStateCard(message: 'No reports available yet'),
              ],
            ),
          );
        }

        final currentItem = (selected != null &&
                allReports.any((r) => r.id == selected!.id))
            ? allReports.firstWhere((r) => r.id == selected!.id)
            : allReports.first;

        return SingleChildScrollView(
          child: PageSection(
            padding: EdgeInsets.fromLTRB(
              compact ? 16 : 24,
              24,
              compact ? 16 : 24,
              32,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const PageIntro(
                  title: 'Reports',
                  subtitle: 'View and analyze health reports across your herd',
                ),
                const SizedBox(height: 20),
                Builder(
                  builder: (context) {
                    final listCard = ReportsListCard(
                      reportsList: allReports,
                      selected: currentItem,
                      onSelect: (item) => setState(() => selected = item),
                      currentPage: page,
                      totalPages: totalPages,
                      onPageChanged: _changePage,
                      totalCount: result.total,
                      pageSize: 8,
                    );
                    final detailCard =
                        ReportDetailsCard(selected: currentItem);

                    if (compact) {
                      return Column(
                        children: [
                          listCard,
                          const SizedBox(height: 16),
                          detailCard,
                        ],
                      );
                    }
                    return Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(flex: 5, child: listCard),
                        const SizedBox(width: 20),
                        Expanded(flex: 5, child: detailCard),
                      ],
                    );
                  },
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  void _changePage(int value) {
    setState(() {
      page = value;
      selected = null;
    });
  }
}
