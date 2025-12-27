import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:async';
import 'settings_page.dart';
import 'session_start_page.dart';

class TouchscreenApp extends StatelessWidget {
  const TouchscreenApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CV40 Touch Control',
      theme: ThemeData.dark().copyWith(
        scaffoldBackgroundColor: Colors.black,
        primaryColor: Colors.white,
      ),
      home: const TouchControlPage(),
      debugShowCheckedModeBanner: false,
    );
  }
}

class TouchControlPage extends StatefulWidget {
  const TouchControlPage({super.key});

  @override
  State<TouchControlPage> createState() => _TouchControlPageState();
}

class _TouchControlPageState extends State<TouchControlPage> with TickerProviderStateMixin {
  bool _isRecording = false;
  bool _isPaused = false;
  bool _sessionStarted = false;
  bool _redBoost = false;
  bool _whiteBalanceInProgress = false;
  String _whiteBalanceMessage = '';
  String _debugMessage = 'System Ready';
  String _connectionStatus = 'Disconnected';

  // API Configuration
  String _apiBase = 'http://localhost:8081';

  late AnimationController _recordAnimationController;
  late AnimationController _whiteBalanceController;

  @override
  void initState() {
    super.initState();
    _recordAnimationController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    _whiteBalanceController = AnimationController(
      duration: const Duration(milliseconds: 1000),
      vsync: this,
    );
    _checkConnection();
    // Check connection every 5 seconds
    Timer.periodic(const Duration(seconds: 5), (timer) => _checkConnection());
  }

  @override
  void dispose() {
    _recordAnimationController.dispose();
    _whiteBalanceController.dispose();
    super.dispose();
  }

  Future<void> _checkConnection() async {
    try {
      final response = await http.get(Uri.parse('$_apiBase/api/settings')).timeout(const Duration(seconds: 3));
      if (response.statusCode == 200) {
        setState(() {
          _connectionStatus = 'Connected';
          _debugMessage = 'Backend connected successfully';
        });
      } else {
        setState(() {
          _connectionStatus = 'Error ${response.statusCode}';
          _debugMessage = 'Backend error: ${response.statusCode}';
        });
      }
    } catch (e) {
      setState(() {
        _connectionStatus = 'Disconnected';
        _debugMessage = 'Connection failed: ${e.toString()}';
      });
    }
  }

  Future<void> _startSession() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (_) => const SessionStartPage()),
    );
    if (result == true) {
      setState(() {
        _sessionStarted = true;
        _debugMessage = 'Surgery session started';
      });
    }
  }

  Future<void> _captureStillImage() async {
    setState(() => _debugMessage = 'Capturing still image...');
    try {
      final response = await http.post(
        Uri.parse('$_apiBase/api/capture'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        setState(() => _debugMessage = 'Still image captured: ${data['results']?.length ?? 0} files');
        _showMessage('Still image captured');
      } else {
        setState(() => _debugMessage = 'Capture failed: ${response.statusCode}');
        _showMessage('Capture failed', isError: true);
      }
    } catch (e) {
      setState(() => _debugMessage = 'Capture error: ${e.toString()}');
      _showMessage('Error capturing image', isError: true);
    }
  }

  Future<void> _toggleRecording() async {
    if (!_sessionStarted) {
      await _startSession();
      if (!_sessionStarted) return;
    }
    
    try {
      if (!_isRecording) {
        setState(() => _debugMessage = 'Starting recording...');
        final response = await http.post(
          Uri.parse('$_apiBase/api/recording/start'),
          headers: {'Content-Type': 'application/json'},
        );
        
        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          setState(() {
            _isRecording = true;
            _isPaused = false;
            _debugMessage = 'Recording started: ${data['workers']?.length ?? 0} workers';
          });
          _recordAnimationController.forward();
        } else {
          setState(() => _debugMessage = 'Recording start failed: ${response.statusCode}');
        }
      } else if (_isPaused) {
        setState(() => _debugMessage = 'Resuming recording...');
        final response = await http.post(
          Uri.parse('$_apiBase/api/recording/resume'),
          headers: {'Content-Type': 'application/json'},
        );
        
        if (response.statusCode == 200) {
          setState(() => _isPaused = false);
        }
      } else {
        final response = await http.post(
          Uri.parse('http://localhost:8081/api/recording/pause'),
          headers: {'Content-Type': 'application/json'},
        );
        
        if (response.statusCode == 200) {
          setState(() => _isPaused = true);
        }
      }
    } catch (e) {
      _showMessage('Recording error', isError: true);
    }
  }

  Future<void> _stopRecording() async {
    setState(() {
      _isRecording = false;
      _isPaused = false;
    });
    _recordAnimationController.reverse();
    
    try {
      final response = await http.post(
        Uri.parse('http://localhost:8081/api/recording/stop'),
        headers: {'Content-Type': 'application/json'},
      );
      
      if (response.statusCode == 200) {
        _showMessage('Recording saved to device(s)');
      }
    } catch (e) {
      _showMessage('Error stopping recording', isError: true);
    }
  }

  Future<void> _executeWhiteBalance() async {
    setState(() {
      _whiteBalanceInProgress = true;
      _whiteBalanceMessage = '';
    });
    _whiteBalanceController.forward();
    
    try {
      final response = await http.post(
        Uri.parse('http://localhost:8081/api/white-balance'),
        headers: {'Content-Type': 'application/json'},
      );
      
      await Future.delayed(const Duration(milliseconds: 1000)); // Simulate 1-second process
      
      if (response.statusCode == 200) {
        setState(() => _whiteBalanceMessage = 'White balance complete');
        Future.delayed(const Duration(seconds: 2), () {
          if (mounted) setState(() => _whiteBalanceMessage = '');
        });
      }
    } catch (e) {
      _showMessage('White balance error', isError: true);
    } finally {
      setState(() => _whiteBalanceInProgress = false);
      _whiteBalanceController.reverse();
    }
  }

  Future<void> _toggleRedBoost() async {
    setState(() => _redBoost = !_redBoost);
    
    try {
      final preset = _redBoost ? 'red_boost' : 'arthroscopy';
      await http.post(
        Uri.parse('http://localhost:8081/api/presets/apply'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({'preset': preset}),
      );
    } catch (e) {
      _showMessage('Preset error', isError: true);
    }
  }

  void _showMessage(String message, {bool isError = false}) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: isError ? Colors.red : Colors.green,
        duration: const Duration(seconds: 2),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        width: 2048,
        height: 1536,
        color: Colors.black,
        child: Column(
          children: [
            // Top Bar - 2048x267
            Container(
              width: 2048,
              height: 267,
              color: Colors.black,
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 95),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    // Left side - Pencil icon and Arthroscopy text
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        // Pencil icon - 150x150
                        Container(
                          width: 150,
                          height: 150,
                          decoration: BoxDecoration(
                            color: _redBoost ? Colors.red[300] : Colors.blue,
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: const Icon(
                            Icons.edit,
                            color: Colors.white,
                            size: 80,
                          ),
                        ),
                        const SizedBox(width: 40),
                        // Arthroscopy text - 878x191
                        Container(
                          width: 878,
                          height: 191,
                          alignment: Alignment.centerLeft,
                          child: Text(
                            'Arthroscopy',
                            style: TextStyle(
                              fontSize: 120,
                              fontWeight: FontWeight.w300,
                              color: _redBoost ? Colors.red[300] : Colors.white,
                              height: 1.0,
                            ),
                          ),
                        ),
                      ],
                    ),
                    // Settings icon - 120x120
                    GestureDetector(
                      onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const SettingsPage())),
                      child: Container(
                        width: 120,
                        height: 120,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: Colors.white.withOpacity(0.1),
                        ),
                        child: const Icon(
                          Icons.settings,
                          color: Colors.white,
                          size: 60,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // Main Control Area - 2048x1249 with 235px top offset
            Container(
              width: 2048,
              height: 1249,
              margin: const EdgeInsets.only(top: 235),
              child: Padding(
                padding: const EdgeInsets.only(
                  top: 171,
                  right: 95,
                  bottom: 171,
                  left: 95,
                ),
                child: Column(
                  children: [
                    // Debug info
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.all(20),
                      decoration: BoxDecoration(
                        color: Colors.grey[900],
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Debug Info:',
                            style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Status: $_connectionStatus',
                            style: TextStyle(
                              color: _connectionStatus == 'Connected' ? Colors.green : Colors.red,
                              fontSize: 14,
                            ),
                          ),
                          Text(
                            'Message: $_debugMessage',
                            style: const TextStyle(color: Colors.white70, fontSize: 14),
                          ),
                          Text(
                            'API: $_apiBase',
                            style: const TextStyle(color: Colors.blue, fontSize: 12),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 40),

                    // Red Boost Toggle (if not recording)
                    if (!_isRecording) ...[
                      Align(
                        alignment: Alignment.centerRight,
                        child: GestureDetector(
                          onTap: _toggleRedBoost,
                          child: Container(
                            padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 15),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(30),
                              color: _redBoost ? Colors.red.withOpacity(0.8) : Colors.white.withOpacity(0.1),
                              border: _redBoost ? Border.all(color: Colors.red, width: 3) : null,
                            ),
                            child: Text(
                              'Red Boost',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 24,
                                fontWeight: _redBoost ? FontWeight.bold : FontWeight.normal,
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 60),
                    ],

                    // Main Control Buttons Row
                    Expanded(
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          // Camera button - 600x600
                          _buildControlButton(
                            icon: Icons.camera_alt,
                            onTap: _captureStillImage,
                            size: 600,
                          ),

                          const SizedBox(width: 79), // Gap between buttons

                          // Record button (center) - 600x600
                          _buildRecordButton(),

                          const SizedBox(width: 79), // Gap between buttons

                          // White Balance button - 600x600
                          _buildWhiteBalanceButton(),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // White balance message overlay
            if (_whiteBalanceMessage.isNotEmpty)
              Positioned(
                bottom: 100,
                left: 0,
                right: 0,
                child: Center(
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 15),
                    decoration: BoxDecoration(
                      color: Colors.green.withOpacity(0.9),
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.3),
                          blurRadius: 10,
                          offset: const Offset(0, 5),
                        ),
                      ],
                    ),
                    child: Text(
                      _whiteBalanceMessage,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildControlButton({required IconData icon, required VoidCallback onTap, double size = 600}) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: const Color(0xFF404040),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.4),
              blurRadius: 30,
              offset: const Offset(0, 15),
            ),
          ],
        ),
        child: Icon(
          icon,
          color: Colors.white,
          size: size * 0.35, // Adjusted icon size ratio
        ),
      ),
    );
  }

  Widget _buildRecordButton() {
    return GestureDetector(
      onTap: _isRecording ? (_isPaused ? _toggleRecording : _stopRecording) : _toggleRecording,
      child: Container(
        width: 600,
        height: 600,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(80), // Rounded square
          color: const Color(0xFF404040),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.4),
              blurRadius: 30,
              offset: const Offset(0, 15),
            ),
          ],
        ),
        child: Center(
          child: Container(
            width: 400,
            height: 400,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: _isRecording ? (_isPaused ? Colors.orange : Colors.red) : Colors.red,
              boxShadow: [
                BoxShadow(
                  color: (_isRecording ? Colors.red : Colors.red).withOpacity(0.6),
                  blurRadius: 20,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            child: Center(
              child: _isRecording
                ? Icon(
                    _isPaused ? Icons.play_arrow : Icons.stop,
                    color: Colors.white,
                    size: 150,
                  )
                : Container(
                    width: 120,
                    height: 120,
                    decoration: const BoxDecoration(
                      shape: BoxShape.circle,
                      color: Colors.white,
                    ),
                  ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildWhiteBalanceButton() {
    return AnimatedBuilder(
      animation: _whiteBalanceController,
      builder: (context, child) {
        return GestureDetector(
          onTap: _whiteBalanceInProgress ? null : _executeWhiteBalance,
          child: Container(
            width: 600,
            height: 600,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: _whiteBalanceInProgress
                ? Color.lerp(const Color(0xFF6C757D), Colors.green, _whiteBalanceController.value)
                : const Color(0xFF6C757D),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.4),
                  blurRadius: 30,
                  offset: const Offset(0, 15),
                ),
              ],
            ),
            child: Center(
              child: Text(
                'WB',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 120,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
        );
      },
    );
  }
}
