import 'package:http/http.dart' as http;
import 'dart:convert';

class CameraApi {
  final String baseUrl;
  CameraApi({String? baseUrl}) : baseUrl = baseUrl ?? 'http://localhost:8081';

  Future<bool> startSession(Map<String, dynamic> meta) async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/session/start'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(meta),
    );
    return res.statusCode == 200;
  }

  Future<bool> capturePhoto() async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/capture'),
      headers: {'Content-Type': 'application/json'},
    );
    return res.statusCode == 200;
  }

  Future<bool> startRecording() async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/recording/start'),
      headers: {'Content-Type': 'application/json'},
    );
    return res.statusCode == 200;
  }

  Future<bool> pauseRecording() async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/recording/pause'),
      headers: {'Content-Type': 'application/json'},
    );
    return res.statusCode == 200;
  }

  Future<bool> resumeRecording() async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/recording/resume'),
      headers: {'Content-Type': 'application/json'},
    );
    return res.statusCode == 200;
  }

  Future<bool> stopRecording() async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/recording/stop'),
      headers: {'Content-Type': 'application/json'},
    );
    return res.statusCode == 200;
  }

  Future<bool> whiteBalance() async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/white-balance'),
      headers: {'Content-Type': 'application/json'},
    );
    return res.statusCode == 200;
  }

  Future<bool> applyPreset(String preset) async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/presets/apply'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'preset': preset}),
    );
    return res.statusCode == 200;
  }

  Future<Map<String, dynamic>?> getSettings() async {
    final res = await http.get(Uri.parse('$baseUrl/api/settings'));
    if (res.statusCode == 200) {
      return json.decode(res.body) as Map<String, dynamic>;
    }
    return null;
  }

  Future<bool> postSettings(Map<String, dynamic> body) async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/settings'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(body),
    );
    return res.statusCode == 200;
  }
}
