import 'package:flutter/material.dart';

import '../../core/models/appModels.dart';
import '../../shared/appWidgets.dart';
import 'reportData.dart';
import 'reportWidgets.dart';

class ReportsPage extends StatefulWidget {
  const ReportsPage({super.key});

  @override
  State<ReportsPage> createState() => _ReportsPageState();
}

class _ReportsPageState extends State<ReportsPage> {
  late final List<ReportItem> allReports = buildReportList();
  int page = 1;
  static const pageSize = 8;
  late ReportItem selected = allReports.first;

  @override
  Widget build(BuildContext context) {
    final compact = MediaQuery.sizeOf(context).width < 1100;

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

    final totalPages = (allReports.length / pageSize).ceil();
    final pageItems = allReports
        .skip((page - 1) * pageSize)
        .take(pageSize)
        .toList();
    final currentItem = pageItems.contains(selected)
        ? selected
        : pageItems.first;

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
                  reportsList: pageItems,
                  selected: currentItem,
                  onSelect: (item) => setState(() => selected = item),
                  currentPage: page,
                  totalPages: totalPages,
                  onPageChanged: _changePage,
                  totalCount: allReports.length,
                  pageSize: pageSize,
                );
                final detailCard = ReportDetailsCard(selected: currentItem);

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
  }

  void _changePage(int value) {
    setState(() {
      page = value;
      selected = allReports.skip((value - 1) * pageSize).first;
    });
  }
}
