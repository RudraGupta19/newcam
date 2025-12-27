import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:async';
import 'dart:math';

class MonitorApp extends StatelessWidget {
  const MonitorApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CV40 Monitor Display',
      theme: ThemeData.dark(),
      home: const MonitorDisplayPage(),
      debugShowCheckedModeBanner: false,
    );
  }
}

class MonitorDisplayPage extends StatefulWidget {
  const MonitorDisplayPage({super.key});

  @override
  State<MonitorDisplayPage> createState() => _MonitorDisplayPageState();
}

class _MonitorDisplayPageState extends State<MonitorDisplayPage> with TickerProviderStateMixin {
  bool _isRecording = false;
  bool _showParameterSlider = false;
  String _currentParameter = 'Zoom';
  double _parameterValue = 1.0;
  Timer? _hideSliderTimer;
  Timer? _statusTimer;
  
  late AnimationController _recordDotController;
  late AnimationController _photoPopupController;
  late AnimationController _sliderController;

  @override
  void initState() {
    super.initState();
    _recordDotController = AnimationController(
      duration: const Duration(milliseconds: 1000),
      vsync: this,
    )..repeat(reverse: true);
    
    _photoPopupController = AnimationController(
      duration: const Duration(milliseconds: 300),
      vsync: this,
    );
    
    _sliderController = AnimationController(
      duration: const Duration(milliseconds: 200),
      vsync: this,
    );

    _startStatusPolling();
  }

  @override
  void dispose() {
    _recordDotController.dispose();
    _photoPopupController.dispose();
    _sliderController.dispose();
    _hideSliderTimer?.cancel();
    _statusTimer?.cancel();
    super.dispose();
  }

  void _startStatusPolling() {
    _statusTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      // In a real implementation, this would poll the backend for recording status
      // For now, we'll simulate status changes
    });
  }

  void _showParameterSlider(String parameter, double value) {
    setState(() {
      _currentParameter = parameter;
      _parameterValue = value;
      _showParameterSlider = true;
    });
    _sliderController.forward();
    
    _hideSliderTimer?.cancel();
    _hideSliderTimer = Timer(const Duration(seconds: 3), () {
      if (mounted) {
        _sliderController.reverse().then((_) {
          if (mounted) {
            setState(() => _showParameterSlider = false);
          }
        });
      }
    });
  }

  void _simulatePhotoCapture() {
    _photoPopupController.forward().then((_) {
      Future.delayed(const Duration(seconds: 2), () {
        if (mounted) _photoPopupController.reverse();
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      body: Stack(
        children: [
          // Main camera feed area (simulated)
          Container(
            width: double.infinity,
            height: double.infinity,
            child: _buildSimulatedCameraFeed(),
          ),

          // Recording indicator (top-right)
          if (_isRecording)
            Positioned(
              top: 40,
              right: 40,
              child: AnimatedBuilder(
                animation: _recordDotController,
                builder: (context, child) {
                  return Container(
                    width: 20,
                    height: 20,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: Color.lerp(Colors.red, Colors.red[300], _recordDotController.value),
                    ),
                  );
                },
              ),
            ),

          // Parameter slider overlay (center)
          if (_showParameterSlider)
            Positioned(
              top: MediaQuery.of(context).size.height * 0.4,
              left: 100,
              right: 100,
              child: SlideTransition(
                position: Tween<Offset>(
                  begin: const Offset(0, -1),
                  end: Offset.zero,
                ).animate(_sliderController),
                child: Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: Colors.black.withOpacity(0.8),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.blue, width: 2),
                  ),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        _currentParameter,
                        style: const TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(height: 16),
                      Row(
                        children: [
                          Expanded(
                            child: SliderTheme(
                              data: SliderTheme.of(context).copyWith(
                                activeTrackColor: Colors.blue,
                                inactiveTrackColor: Colors.grey[700],
                                thumbColor: Colors.blue,
                                trackHeight: 8,
                                thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 16),
                              ),
                              child: Slider(
                                value: _parameterValue,
                                min: _currentParameter == 'Zoom' ? 1.0 : -100,
                                max: _currentParameter == 'Zoom' ? 4.0 : 100,
                                onChanged: (value) => setState(() => _parameterValue = value),
                              ),
                            ),
                          ),
                          Container(
                            width: 80,
                            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                            decoration: BoxDecoration(
                              color: Colors.blue.withOpacity(0.2),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              _currentParameter == 'Zoom' 
                                ? '${_parameterValue.toStringAsFixed(1)}x'
                                : '${_parameterValue.round()}',
                              style: const TextStyle(color: Colors.blue, fontSize: 16, fontWeight: FontWeight.bold),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),

          // Photo capture popup
          Positioned(
            bottom: 100,
            right: 100,
            child: ScaleTransition(
              scale: _photoPopupController,
              child: Container(
                width: 150,
                height: 100,
                decoration: BoxDecoration(
                  color: Colors.grey[800],
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.white, width: 2),
                ),
                child: const Center(
                  child: Text(
                    'Photo\nCaptured',
                    style: TextStyle(color: Colors.white, fontSize: 14),
                    textAlign: TextAlign.center,
                  ),
                ),
              ),
            ),
          ),

          // Test controls (bottom-left, for simulation)
          Positioned(
            bottom: 20,
            left: 20,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                ElevatedButton(
                  onPressed: () => setState(() => _isRecording = !_isRecording),
                  child: Text(_isRecording ? 'Stop Recording' : 'Start Recording'),
                ),
                ElevatedButton(
                  onPressed: () => _showParameterSlider('Zoom', 2.5),
                  child: const Text('Show Zoom Slider'),
                ),
                ElevatedButton(
                  onPressed: () => _showParameterSlider('Brightness', 25),
                  child: const Text('Show Brightness Slider'),
                ),
                ElevatedButton(
                  onPressed: _simulatePhotoCapture,
                  child: const Text('Simulate Photo'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSimulatedCameraFeed() {
    return Container(
      decoration: BoxDecoration(
        gradient: RadialGradient(
          center: const Alignment(0.3, -0.2),
          radius: 1.2,
          colors: [
            Colors.pink[100]!,
            Colors.red[200]!,
            Colors.red[400]!,
            Colors.red[800]!,
          ],
        ),
      ),
      child: CustomPaint(
        painter: SimulatedTissuePainter(),
        size: Size.infinite,
      ),
    );
  }
}

class SimulatedTissuePainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()..style = PaintingStyle.fill;
    
    // Draw some organic shapes to simulate tissue
    final random = Random(42); // Fixed seed for consistent appearance
    
    for (int i = 0; i < 20; i++) {
      final center = Offset(
        random.nextDouble() * size.width,
        random.nextDouble() * size.height,
      );
      final radius = 20 + random.nextDouble() * 60;
      
      paint.color = Color.lerp(
        Colors.pink[200]!,
        Colors.red[300]!,
        random.nextDouble(),
      )!.withOpacity(0.6);
      
      canvas.drawCircle(center, radius, paint);
    }
    
    // Add some vessel-like lines
    paint.style = PaintingStyle.stroke;
    paint.strokeWidth = 3;
    paint.color = Colors.red[600]!;
    
    for (int i = 0; i < 10; i++) {
      final path = Path();
      path.moveTo(random.nextDouble() * size.width, random.nextDouble() * size.height);
      path.quadraticBezierTo(
        random.nextDouble() * size.width,
        random.nextDouble() * size.height,
        random.nextDouble() * size.width,
        random.nextDouble() * size.height,
      );
      canvas.drawPath(path, paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
