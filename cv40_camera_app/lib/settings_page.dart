import 'package:flutter/material.dart';
import 'services/camera_api.dart';

class SettingsPage extends StatefulWidget {
  final CameraApi api;
  const SettingsPage({super.key, required this.api});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> with TickerProviderStateMixin {
  late TabController _tabController;
  
  // Settings values
  double _brightness = 0.0;
  double _zoom = 1.0;
  double _contrast = 0.0;
  double _sharpness = 0.0;
  double _saturation = 0.0;
  double _temperature = 6500.0;
  double _hue = 0.0;
  bool _lowLightGain = false;
  bool _redBoost = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadSettings();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadSettings() async {
    try {
      final data = await widget.api.getSettings();
      if (data != null) {
        setState(() {
          _brightness = (data['colors']?['brightness'] ?? 0).toDouble();
          _zoom = (data['visuals']?['zoom'] ?? 1.0).toDouble();
          _contrast = (data['colors']?['contrast'] ?? 0).toDouble();
          _sharpness = (data['visuals']?['sharpness'] ?? 0.0).toDouble();
          _saturation = (data['colors']?['saturation'] ?? 0).toDouble();
          _temperature = (data['white']?['temperature'] ?? 6500).toDouble();
          _hue = (data['colors']?['hue'] ?? 0).toDouble();
          _lowLightGain = (data['exposure']?['lowLightGain'] ?? 0.0) > 0.0;
        });
      }
    } catch (e) {
      print('Error loading settings: $e');
    }
  }

  Future<void> _saveSettings() async {
    try {
      final body = {
        'visuals': {'zoom': _zoom, 'sharpness': _sharpness},
        'colors': {'brightness': _brightness.round(), 'contrast': _contrast.round(), 'saturation': _saturation.round(), 'hue': _hue.round()},
        'white': {'temperature': _temperature.round()},
        'exposure': {'lowLightGain': _lowLightGain ? 1.0 : 0.0}
      };
      await widget.api.postSettings(body);
    } catch (e) {
      print('Error saving settings: $e');
    }
  }

  Widget _buildSlider(String label, double value, double min, double max, ValueChanged<double> onChanged, {String? unit}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(label, style: const TextStyle(color: Colors.white, fontSize: 24)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: Colors.blue.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  '${value.round()}${unit ?? ''}',
                  style: const TextStyle(color: Colors.blue, fontSize: 18, fontWeight: FontWeight.bold),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          SliderTheme(
            data: SliderTheme.of(context).copyWith(
              activeTrackColor: Colors.blue,
              inactiveTrackColor: Colors.grey[700],
              thumbColor: Colors.blue,
              overlayColor: Colors.blue.withOpacity(0.2),
              trackHeight: 6,
              thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 12),
            ),
            child: Slider(
              value: value,
              min: min,
              max: max,
              onChanged: (v) {
                onChanged(v);
                _saveSettings();
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildToggle(String label, bool value, ValueChanged<bool> onChanged) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.white, fontSize: 24)),
          Switch(
            value: value,
            onChanged: (v) {
              onChanged(v);
              _saveSettings();
            },
            activeColor: Colors.blue,
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white, size: 32),
          onPressed: () => Navigator.pop(context),
        ),
        title: const Text('Settings', style: TextStyle(color: Colors.white, fontSize: 32)),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: Colors.blue,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.grey,
          labelStyle: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          tabs: const [
            Tab(text: 'Primary'),
            Tab(text: 'Colour'),
            Tab(text: 'Advanced'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          // Primary Tab
          Padding(
            padding: const EdgeInsets.all(32),
            child: Column(
              children: [
                _buildSlider('Brightness', _brightness, -100, 100, (v) => setState(() => _brightness = v)),
                _buildSlider('Zoom', _zoom, 1.0, 4.0, (v) => setState(() => _zoom = v), unit: 'x'),
                _buildSlider('Contrast', _contrast, -100, 100, (v) => setState(() => _contrast = v)),
                _buildToggle('Low light gain', _lowLightGain, (v) => setState(() => _lowLightGain = v)),
              ],
            ),
          ),
          // Colour Tab
          Padding(
            padding: const EdgeInsets.all(32),
            child: Column(
              children: [
                _buildSlider('Saturation', _saturation, -100, 100, (v) => setState(() => _saturation = v)),
                _buildSlider('Temperature', _temperature, 2000, 10000, (v) => setState(() => _temperature = v), unit: 'K'),
                _buildSlider('Hue', _hue, -180, 180, (v) => setState(() => _hue = v)),
              ],
            ),
          ),
          // Advanced Tab
          Padding(
            padding: const EdgeInsets.all(32),
            child: Column(
              children: [
                _buildSlider('Sharpness', _sharpness, 0, 2.0, (v) => setState(() => _sharpness = v)),
                const SizedBox(height: 40),
                const Text('RGB', style: TextStyle(color: Colors.white, fontSize: 24)),
                const SizedBox(height: 20),
                // RGB controls would go here - simplified for now
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: Colors.grey[900],
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Text(
                    'RGB color gain controls\n(Advanced tuning)',
                    style: TextStyle(color: Colors.grey, fontSize: 16),
                    textAlign: TextAlign.center,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
