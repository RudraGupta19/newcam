import 'package:flutter/material.dart';
import 'services/camera_api.dart';

class SessionStartPage extends StatefulWidget {
  final CameraApi api;
  const SessionStartPage({super.key, required this.api});

  @override
  State<SessionStartPage> createState() => _SessionStartPageState();
}

class _SessionStartPageState extends State<SessionStartPage> {
  final _formKey = GlobalKey<FormState>();
  final _doctorController = TextEditingController();
  final _hospitalController = TextEditingController();
  final _surgeryController = TextEditingController();
  final _patientController = TextEditingController();
  final _technicianController = TextEditingController();

  @override
  void dispose() {
    _doctorController.dispose();
    _hospitalController.dispose();
    _surgeryController.dispose();
    _patientController.dispose();
    _technicianController.dispose();
    super.dispose();
  }

  Future<void> _startSession() async {
    if (!_formKey.currentState!.validate()) return;

    try {
      final body = {
        'doctor': _doctorController.text,
        'hospital': _hospitalController.text,
        'surgery': _surgeryController.text,
        'patient': _patientController.text,
        'technician': _technicianController.text,
      };

      final ok = await widget.api.startSession(body);
      if (ok) {
        Navigator.pop(context, true); // Return success
      } else {
        _showError('Failed to start session');
      }
    } catch (e) {
      _showError('Error: $e');
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  Widget _buildTextField(String label, TextEditingController controller, {bool required = true}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 12),
      child: TextFormField(
        controller: controller,
        style: const TextStyle(color: Colors.white, fontSize: 18),
        decoration: InputDecoration(
          labelText: label,
          labelStyle: const TextStyle(color: Colors.grey, fontSize: 16),
          enabledBorder: OutlineInputBorder(
            borderSide: const BorderSide(color: Colors.grey),
            borderRadius: BorderRadius.circular(8),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: const BorderSide(color: Colors.blue),
            borderRadius: BorderRadius.circular(8),
          ),
          errorBorder: OutlineInputBorder(
            borderSide: const BorderSide(color: Colors.red),
            borderRadius: BorderRadius.circular(8),
          ),
          focusedErrorBorder: OutlineInputBorder(
            borderSide: const BorderSide(color: Colors.red),
            borderRadius: BorderRadius.circular(8),
          ),
        ),
        validator: required ? (value) => value?.isEmpty == true ? 'Required' : null : null,
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
        title: const Text('Start Surgery', style: TextStyle(color: Colors.white, fontSize: 28)),
      ),
      body: Padding(
        padding: const EdgeInsets.all(32),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                'Enter surgery details:',
                style: TextStyle(color: Colors.white, fontSize: 20),
              ),
              const SizedBox(height: 32),
              Expanded(
                child: SingleChildScrollView(
                  child: Column(
                    children: [
                      _buildTextField('Doctor Name', _doctorController),
                      _buildTextField('Hospital Name', _hospitalController),
                      _buildTextField('Surgery Type', _surgeryController),
                      _buildTextField('Patient Name', _patientController),
                      _buildTextField('Technician Name', _technicianController, required: false),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 32),
              ElevatedButton(
                onPressed: _startSession,
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.blue,
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                ),
                child: const Text(
                  'Start Surgery',
                  style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
