import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:flutter_map_cancellable_tile_provider/flutter_map_cancellable_tile_provider.dart';
import 'package:latlong2/latlong.dart';
import 'package:frontend/core/models/appModels.dart';
import 'package:frontend/shared/appWidgets.dart';

class MovementMapCard extends StatelessWidget {
  const MovementMapCard({super.key, required this.response});

  final MovementPathResponse response;

  // farm center fallback
  static const _defaultLat = -37.7870;
  static const _defaultLng = 175.2793;

  @override
  Widget build(BuildContext context) {
    if (response.points.isEmpty) {
      return _emptyState(context);
    }

    final points = response.points;
    final latLngs = points.map((p) => LatLng(p.lat, p.lng)).toList();
    final bounds = LatLngBounds.fromPoints(latLngs);
    final stayPoints = points.where((p) => p.staySeconds > 0).toList();

    return SurfaceCard(
      child: SizedBox(
        height: 300,
        child: ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: FlutterMap(
            options: MapOptions(
              initialCameraFit: CameraFit.bounds(
                bounds: bounds,
                padding: const EdgeInsets.all(40),
              ),
            ),
            children: [
              TileLayer(
                urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                userAgentPackageName: 'com.compx576.smartdairy',
                tileProvider: CancellableNetworkTileProvider(),
              ),
              // movement path line
              PolylineLayer(
                polylines: [
                  Polyline(
                    points: latLngs,
                    color: Colors.blue,
                    strokeWidth: 3,
                  ),
                ],
              ),
              // stay heatmap circles
              CircleLayer(
                circles: stayPoints.map((p) => _stayCircle(p)).toList(),
              ),
              // start and end markers
              MarkerLayer(
                markers: [
                  _marker(latLngs.first, Colors.green, 'Start'),
                  _marker(latLngs.last, Colors.red, 'End'),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _emptyState(BuildContext context) {
    return SurfaceCard(
      child: SizedBox(
        height: 300,
        child: ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: Stack(
            children: [
              FlutterMap(
                options: MapOptions(
                  initialCenter: const LatLng(_defaultLat, _defaultLng),
                  initialZoom: 15,
                ),
                children: [
                  TileLayer(
                    urlTemplate:
                        'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                    userAgentPackageName: 'com.compx576.smartdairy',
                    tileProvider: CancellableNetworkTileProvider(),
                  ),
                ],
              ),
              Center(
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  decoration: BoxDecoration(
                    color: Colors.black54,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Text(
                    'No movement data',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  CircleMarker _stayCircle(MovementPathPoint p) {
    final minutes = p.staySeconds / 60;
    Color color;
    double radius;
    if (minutes > 15) {
      color = Colors.red.withOpacity(0.4);
      radius = 20;
    } else if (minutes > 5) {
      color = Colors.orange.withOpacity(0.4);
      radius = 14;
    } else {
      color = Colors.orange.withOpacity(0.25);
      radius = 8;
    }

    return CircleMarker(
      point: LatLng(p.lat, p.lng),
      radius: radius,
      color: color,
      borderColor: color.withOpacity(0.8),
      borderStrokeWidth: 1,
    );
  }

  Marker _marker(LatLng point, Color color, String label) {
    return Marker(
      point: point,
      width: 16,
      height: 16,
      child: Tooltip(
        message: label,
        child: Container(
          decoration: BoxDecoration(
            color: color,
            shape: BoxShape.circle,
            border: Border.all(color: Colors.white, width: 2),
          ),
        ),
      ),
    );
  }
}
