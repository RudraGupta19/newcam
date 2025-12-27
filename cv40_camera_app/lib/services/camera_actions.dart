import 'camera_api.dart';

class CameraActions {
  final CameraApi api;
  CameraActions(this.api);

  double _clamp(double v, double min, double max) {
    if (v < min) return min;
    if (v > max) return max;
    return v;
  }

  Future<bool> ensureSession(Map<String, dynamic> meta) async {
    return await api.startSession(meta);
  }

  Future<bool> toggleRecording(bool isRecording, bool isPaused, {Map<String, dynamic>? sessionMeta}) async {
    if (!isRecording) {
      if (sessionMeta != null) {
        await ensureSession(sessionMeta);
      }
      return await api.startRecording();
    }
    if (isPaused) {
      return await api.resumeRecording();
    }
    return await api.pauseRecording();
  }

  Future<bool> stopRecording() async {
    return await api.stopRecording();
  }

  Future<bool> capturePhoto() async {
    return await api.capturePhoto();
  }

  Future<bool> whiteBalance() async {
    return await api.whiteBalance();
  }

  Future<bool> applyPreset(String preset) async {
    return await api.applyPreset(preset);
  }

  Future<bool> setBrightness(double v) async {
    final b = _clamp(v, -100, 100).round();
    return await api.postSettings({"colors": {"brightness": b}});
  }

  Future<bool> setContrast(double v) async {
    final c = _clamp(v, -100, 100).round();
    return await api.postSettings({"colors": {"contrast": c}});
  }

  Future<bool> setZoom(double v) async {
    final z = _clamp(v, 1.0, 4.0);
    return await api.postSettings({"visuals": {"zoom": z}});
  }

  Future<bool> setSharpness(double v) async {
    final s = _clamp(v, 0.0, 2.0);
    return await api.postSettings({"visuals": {"sharpness": s}});
  }

  Future<bool> setSaturation(double v) async {
    final s = _clamp(v, -100.0, 100.0).round();
    return await api.postSettings({"colors": {"saturation": s}});
  }

  Future<bool> setTemperature(double v) async {
    final t = _clamp(v, 2000.0, 10000.0).round();
    return await api.postSettings({"white": {"temperature": t}});
  }

  Future<bool> setHue(double v) async {
    final h = _clamp(v, -180.0, 180.0).round();
    return await api.postSettings({"colors": {"hue": h}});
  }

  Future<bool> setLowLightGain(bool on) async {
    return await api.postSettings({"exposure": {"lowLightGain": on ? 1.0 : 0.0}});
  }
}
