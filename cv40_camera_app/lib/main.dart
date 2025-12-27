import 'package:flutter/material.dart';
import 'services/camera_api.dart';
import 'settings_page.dart';
import 'session_start_page.dart';

void main() {
  runApp(const CV40CameraApp());
}

class CV40CameraApp extends StatelessWidget {
  const CV40CameraApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CV40 Camera System',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      home: const HomePage(),
      debugShowCheckedModeBanner: false,
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> with TickerProviderStateMixin {
  bool _isRecording = false;
  bool _isPaused = false;
  bool _sessionStarted = false;
  bool _redBoost = false;
  late CameraApi _api;
  
  late AnimationController _recordAnimationController;
  late Animation<double> _recordSplitAnimation;
  late Animation<Offset> _pauseButtonSlideAnimation;
  late Animation<Offset> _stopButtonSlideAnimation;
  
  @override
  void initState() {
    super.initState();
    _api = CameraApi();
    
    _recordAnimationController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    
    _recordSplitAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _recordAnimationController,
      curve: Curves.elasticOut,
    ));
    
    _pauseButtonSlideAnimation = Tween<Offset>(
      begin: const Offset(0, 0),
      end: const Offset(-0.6, 0),
    ).animate(CurvedAnimation(
      parent: _recordAnimationController,
      curve: Curves.elasticOut,
    ));
    
    _stopButtonSlideAnimation = Tween<Offset>(
      begin: const Offset(0, 0),
      end: const Offset(0.6, 0),
    ).animate(CurvedAnimation(
      parent: _recordAnimationController,
      curve: Curves.elasticOut,
    ));
  }
  
  @override
  void dispose() {
    _recordAnimationController.dispose();
    super.dispose();
  }

  Future<void> _captureStillImage() async {
    final ok = await _api.capturePhoto();
    if (ok) {
      print('Still image captured successfully');
    }
  }

  Future<void> _startSession() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (_) => SessionStartPage(api: _api)),
    );
    if (result == true) {
      setState(() => _sessionStarted = true);
    }
  }

  Future<void> _toggleRecording() async {
    if (!_sessionStarted) {
      await _startSession();
      if (!_sessionStarted) return;
    }

    if (!_isRecording) {
      final ok = await _api.startRecording();
      if (ok) {
        setState(() {
          _isRecording = true;
          _isPaused = false;
        });
        _recordAnimationController.forward();
        print('Recording started successfully');
      }
    } else if (_isPaused) {
      final ok = await _api.resumeRecording();
      if (ok) {
        setState(() {
          _isPaused = false;
        });
        print('Recording resumed successfully');
      }
    } else {
      final ok = await _api.pauseRecording();
      if (ok) {
        setState(() {
          _isPaused = true;
        });
        print('Recording paused successfully');
      }
    }
  }

  Future<void> _stopRecording() async {
    setState(() {
      _isRecording = false;
      _isPaused = false;
    });
    
    _recordAnimationController.reverse();
    
    final ok = await _api.stopRecording();
    if (ok) {
      print('Recording stopped successfully');
    }
  }

  Future<void> _executeWhiteBalance() async {
    final ok = await _api.whiteBalance();
    if (ok) {
      print('White balance executed successfully');
    }
  }

  Future<void> _toggleRedBoost() async {
    setState(() => _redBoost = !_redBoost);
    final preset = _redBoost ? 'red_boost' : 'arthroscopy';
    final ok = await _api.applyPreset(preset);
    if (ok) {
      print('Preset $preset applied successfully');
    }
  }

  @override
  Widget build(BuildContext context) {
    final screenSize = MediaQuery.of(context).size;
    final buttonSize = screenSize.width > 600 ? 140.0 : 100.0;
    final recordButtonSize = screenSize.width > 600 ? 180.0 : 130.0;
    
    return Scaffold(
      body: Container(
        width: double.infinity,
        height: double.infinity,
        color: Colors.black,
        child: Stack(
          children: [
            // Top header
            Positioned(
              top: 60,
              left: 40,
              right: 40,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: [
                      if (_isRecording)
                        Padding(
                          padding: const EdgeInsets.only(right: 20),
                          child: GestureDetector(
                            onTap: _stopRecording,
                            child: Container(
                              padding: const EdgeInsets.all(12),
                              decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                color: Colors.white.withOpacity(0.1),
                              ),
                              child: const Icon(
                                Icons.arrow_back,
                                color: Colors.white,
                                size: 28,
                              ),
                            ),
                          ),
                        ),
                      Text(
                        'Arthroscopy',
                        style: TextStyle(
                          fontSize: screenSize.width > 600 ? 56 : 42,
                          fontWeight: FontWeight.w300,
                          color: Colors.white,
                        ),
                      ),
                      if (!_isRecording) ...[
                        const SizedBox(width: 20),
                        GestureDetector(
                          onTap: _toggleRedBoost,
                          child: Container(
                            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(20),
                              color: _redBoost ? Colors.red.withOpacity(0.8) : Colors.white.withOpacity(0.1),
                            ),
                            child: Text(
                              'Red Boost',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                                fontWeight: _redBoost ? FontWeight.bold : FontWeight.normal,
                              ),
                            ),
                          ),
                        ),
                      ],
                    ],
                  ),
                  GestureDetector(
                    onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => SettingsPage(api: _api))),
                    child: Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: Colors.white.withOpacity(0.1),
                      ),
                      child: const Icon(
                        Icons.settings,
                        color: Colors.white,
                        size: 28,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            
            // Recording controls (center)
            if (_isRecording)
              Positioned(
                top: 200,
                left: 0,
                right: 0,
                bottom: 200,
                child: AnimatedBuilder(
                  animation: _recordSplitAnimation,
                  builder: (context, child) {
                    return Center(
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          // Pause button (left side)
                          SlideTransition(
                            position: _pauseButtonSlideAnimation,
                            child: GestureDetector(
                              onTap: _toggleRecording,
                              child: Container(
                                width: recordButtonSize,
                                height: recordButtonSize,
                                decoration: BoxDecoration(
                                  shape: BoxShape.circle,
                                  color: const Color(0xFF404040),
                                  boxShadow: [
                                    BoxShadow(
                                      color: Colors.black.withOpacity(0.3),
                                      blurRadius: 20,
                                      offset: const Offset(0, 10),
                                    ),
                                  ],
                                ),
                                child: Center(
                                  child: Icon(
                                    _isPaused ? Icons.play_arrow : Icons.pause,
                                    color: Colors.white,
                                    size: recordButtonSize * 0.4,
                                  ),
                                ),
                              ),
                            ),
                          ),
                          
                          SizedBox(width: 60 * _recordSplitAnimation.value),
                          
                          // Stop button (right side)
                          SlideTransition(
                            position: _stopButtonSlideAnimation,
                            child: GestureDetector(
                              onTap: _stopRecording,
                              child: Container(
                                width: recordButtonSize,
                                height: recordButtonSize,
                                decoration: BoxDecoration(
                                  shape: BoxShape.circle,
                                  color: const Color(0xFFDC3545),
                                  boxShadow: [
                                    BoxShadow(
                                      color: Colors.black.withOpacity(0.3),
                                      blurRadius: 20,
                                      offset: const Offset(0, 10),
                                    ),
                                  ],
                                ),
                                child: Center(
                                  child: Icon(
                                    Icons.stop,
                                    color: Colors.white,
                                    size: recordButtonSize * 0.4,
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),
            
            // Bottom control buttons (only when not recording)
            if (!_isRecording)
              Positioned(
                bottom: 100,
                left: 0,
                right: 0,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    // Camera button
                    GestureDetector(
                      onTap: _captureStillImage,
                      child: Container(
                        width: buttonSize,
                        height: buttonSize,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: const Color(0xFF404040),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.3),
                              blurRadius: 20,
                              offset: const Offset(0, 10),
                            ),
                          ],
                        ),
                        child: Center(
                          child: Icon(
                            Icons.camera_alt,
                            color: Colors.white,
                            size: buttonSize * 0.4,
                          ),
                        ),
                      ),
                    ),
                    
                    // Record button
                    GestureDetector(
                      onTap: _toggleRecording,
                      child: Container(
                        width: recordButtonSize,
                        height: recordButtonSize,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: Colors.red,
                          boxShadow: [
                            BoxShadow(
                              color: Colors.red.withOpacity(0.4),
                              blurRadius: 30,
                              offset: const Offset(0, 10),
                            ),
                          ],
                        ),
                        child: Center(
                          child: Container(
                            width: recordButtonSize * 0.3,
                            height: recordButtonSize * 0.3,
                            decoration: const BoxDecoration(
                              shape: BoxShape.circle,
                              color: Colors.white,
                            ),
                          ),
                        ),
                      ),
                    ),
                    
                    // White Balance button
                    GestureDetector(
                      onTap: _executeWhiteBalance,
                      child: Container(
                        width: buttonSize,
                        height: buttonSize,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: const Color(0xFF6C757D),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.3),
                              blurRadius: 20,
                              offset: const Offset(0, 10),
                            ),
                          ],
                        ),
                        child: Center(
                          child: Text(
                            'WB',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: buttonSize * 0.25,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
          ],
        ),
      ),
    );
  }
}
