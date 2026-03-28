import 'package:flutter_test/flutter_test.dart';

import 'package:frontend/app.dart';

void main() {
  testWidgets('app builds login page', (tester) async {
    await tester.pumpWidget(const SmartDairyApp());
    await tester.pumpAndSettle();

    expect(find.text('Dairy Farm Management'), findsOneWidget);
    expect(find.text('Sign In'), findsOneWidget);
  });
}
