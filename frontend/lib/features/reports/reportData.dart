import '../../core/mock/appMockData.dart';
import '../../core/models/appModels.dart';

List<ReportItem> buildReportList() {
  return [
    ...reports,
    ReportItem(
      id: 'report_004',
      cow_id: 'cow_007',
      cow_name: 'Bella',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Slightly elevated temperature. Monitor for next 24 hours.',
      score: 72,
      details:
          'Temperature trending upward. Heart rate normal. Blood oxygen stable.',
      created_at: DateTime(2026, 3, 22, 16, 12, 34),
    ),
    ReportItem(
      id: 'report_005',
      cow_id: 'cow_008',
      cow_name: 'Rose',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Critical blood oxygen levels. Immediate attention required.',
      score: 65,
      details:
          'Blood oxygen stayed low. Need urgent physical check and device verification.',
      created_at: DateTime(2026, 3, 23, 9, 22, 0),
    ),
    ReportItem(
      id: 'report_006',
      cow_id: 'cow_009',
      cow_name: 'Sophie',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Excellent health. All metrics within normal range.',
      score: 92,
      details: 'Stable readings across all major metrics. No action needed.',
      created_at: DateTime(2026, 3, 20, 10, 5, 0),
    ),
    ReportItem(
      id: 'report_007',
      cow_id: 'cow_010',
      cow_name: 'Rosie',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Good health with stable vitals and strong production.',
      score: 88,
      details:
          'Production is stable. Continue current feeding and observation plan.',
      created_at: DateTime(2026, 3, 19, 8, 45, 0),
    ),
    ReportItem(
      id: 'report_008',
      cow_id: 'cow_011',
      cow_name: 'Maggie',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Elevated heart rate detected. Continue monitoring.',
      score: 78,
      details: 'Need closer watch on heart rate pattern during next 24 hours.',
      created_at: DateTime(2026, 3, 23, 11, 0, 0),
    ),
    ReportItem(
      id: 'report_009',
      cow_id: 'cow_012',
      cow_name: 'Buttercup',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Temperature slightly elevated. Continue observation.',
      score: 75,
      details:
          'Early warning only. Keep routine checks and hydration monitoring.',
      created_at: DateTime(2026, 3, 21, 7, 30, 0),
    ),
    ReportItem(
      id: 'report_010',
      cow_id: 'cow_013',
      cow_name: 'Penny',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Excellent overall health with strong production.',
      score: 90,
      details: 'No sign of anomaly. Good milk output and stable movement.',
      created_at: DateTime(2026, 3, 18, 13, 0, 0),
    ),
    ReportItem(
      id: 'report_011',
      cow_id: 'cow_014',
      cow_name: 'Hazel',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary: 'Temperature and heart rate slightly elevated.',
      score: 73,
      details: 'Need closer watch on temperature and heart rate trend.',
      created_at: DateTime(2026, 3, 20, 17, 0, 0),
    ),
    ReportItem(
      id: 'report_012',
      cow_id: 'cow_015',
      cow_name: 'Poppy',
      period_start: DateTime.now().subtract(const Duration(days: 7)),
      period_end: DateTime.now(),
      summary:
          'Condition stays stable. Milk amount is good and daily movement is normal.',
      score: 88,
      details:
          'No unusual temperature trend. Heart rate is stable. Keep current feed plan.',
      created_at: DateTime(2026, 3, 17, 9, 0, 0),
    ),
  ];
}
