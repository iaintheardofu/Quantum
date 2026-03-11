import 'dart:async';
import 'dart:convert';
import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

// Central Logger for SomaOS
class SomaLogger {
  static final channel = WebSocketChannel.connect(
    Uri.parse('ws://localhost:8082/log'),
  );

  static void log(String message) {
    try {
      channel.sink.add(message);
    } catch (e) {
      // Silently fail if websocket is down
    }
  }
}

void main() {
  // Override print to stream to our master logger
  runZonedGuarded(() {
    runApp(const SomaApp());
  }, (error, stackTrace) {
    SomaLogger.log('[ERROR] $error');
  }, zoneSpecification: ZoneSpecification(
    print: (Zone self, ZoneDelegate parent, Zone zone, String line) {
      SomaLogger.log(line);
      parent.print(zone, line);
    },
  ));
}

class SomaApp extends StatelessWidget {
  const SomaApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'SomaOS HPQC Dashboard',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        brightness: Brightness.dark,
        scaffoldBackgroundColor: const Color(0xFF0A0F1D),
        primaryColor: const Color(0xFF00FFCC),
        fontFamily: 'monospace',
      ),
      home: const DashboardPage(),
    );
  }
}

class HardwareState {
  final int register;
  final double thermalLoad;
  final double phaseField;
  final int activeCells;
  final String routingMode;

  HardwareState({
    required this.register,
    required this.thermalLoad,
    required this.phaseField,
    required this.activeCells,
    required this.routingMode,
  });

  factory HardwareState.fromJson(Map<String, dynamic> json) {
    return HardwareState(
      register: json['register'] ?? 0,
      thermalLoad: (json['thermal_load'] ?? 0).toDouble(),
      phaseField: (json['phase_field'] ?? 0).toDouble(),
      activeCells: json['active_cells'] ?? 8,
      routingMode: json['routing_mode'] ?? 'idle',
    );
  }
}

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  HardwareState _state = HardwareState(
    register: 0,
    thermalLoad: 35.0,
    phaseField: 0.0,
    activeCells: 8,
    routingMode: 'idle',
  );
  Timer? _timer;
  bool _isIdeOpen = false;

  @override
  void initState() {
    super.initState();
    print('>> SomaOS Dashboard Initializing...');
    _startPolling();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void _startPolling() {
    _timer = Timer.periodic(const Duration(milliseconds: 100), (timer) async {
      try {
        final response = await http.get(Uri.parse('http://localhost:8081/api/state'));
        if (response.statusCode == 200) {
          if (mounted) {
            setState(() {
              _state = HardwareState.fromJson(json.decode(response.body));
            });
          }
        }
      } catch (e) {
        // Silently fail
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          Row(
            children: [
              Container(
                width: 300,
                color: const Color(0xFF0F1423),
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('SomaOS v3.5', 
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Color(0xFF00FFCC))
                    ),
                    const SizedBox(height: 10),
                    const Text('HPQC VIRTUALIZATION', style: TextStyle(fontSize: 10, letterSpacing: 2, color: Colors.grey)),
                    const SizedBox(height: 40),
                    _buildStatCard('XADC Thermal', '${_state.thermalLoad.toStringAsFixed(2)} °C', Icons.thermostat, Colors.red),
                    _buildStatCard('SPHY Phase', 'Φ ${_state.phaseField.toStringAsFixed(3)}', Icons.radio, Colors.green),
                    _buildStatCard('Topology', 'd=2^${_state.activeCells}', Icons.memory, Colors.blue),
                    const Spacer(),
                    ElevatedButton.icon(
                      onPressed: () {
                        print('>> Opening ClojureV IDE...');
                        setState(() => _isIdeOpen = true);
                      },
                      icon: const Icon(Icons.code),
                      label: const Text('OPEN CLOJUREV IDE'),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.blueAccent,
                        minimumSize: const Size(double.infinity, 50),
                      ),
                    ),
                  ],
                ),
              ),
              Expanded(
                child: Container(
                  padding: const EdgeInsets.all(40),
                  child: Column(
                    children: [
                      Text(
                        _state.routingMode == 'station' 
                          ? 'FRACTAL HYPERCUBE: 64-QUBIT STATION HUB' 
                          : 'TOPOLOGICAL ENTANGLEMENT BUS: 8-QUBIT MACRO-CUBE',
                        style: const TextStyle(fontSize: 18, color: Color(0xFFAEBCE0)),
                      ),
                      const SizedBox(height: 20),
                      Expanded(
                        child: CustomPaint(
                          painter: MobiusPainter(_state),
                          child: Container(),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
          Positioned(
            top: 40,
            right: 40,
            child: FloatingFlowWindow(state: _state),
          ),
          if (_isIdeOpen) 
            Positioned.fill(
              child: ClojureVIDE(onClose: () => setState(() => _isIdeOpen = false)),
            ),
        ],
      ),
    );
  }

  Widget _buildStatCard(String label, String value, IconData icon, Color color) {
    return Container(
      margin: const EdgeInsets.only(bottom: 20),
      padding: const EdgeInsets.all(15),
      decoration: BoxDecoration(
        color: const Color(0xFF1A2238),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Row(
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(width: 15),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label, style: const TextStyle(fontSize: 12, color: Colors.grey)),
              Text(value, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            ],
          ),
        ],
      ),
    );
  }
}

class MobiusPainter extends CustomPainter {
  final HardwareState state;
  MobiusPainter(this.state);

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final paint = Paint()
      ..color = const Color(0xFF00FFCC).withOpacity(0.5)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2;

    final path = Path();
    for (double i = 0; i < 2 * math.pi; i += 0.05) {
      double x = math.sin(i) * 200;
      double y = math.sin(i) * math.cos(i) * 100;
      double rotatedX = x * math.cos(state.phaseField) - y * math.sin(state.phaseField);
      double rotatedY = x * math.sin(state.phaseField) + y * math.cos(state.phaseField);
      if (i == 0) {
        path.moveTo(center.dx + rotatedX, center.dy + rotatedY);
      } else {
        path.lineTo(center.dx + rotatedX, center.dy + rotatedY);
      }
    }
    path.close();
    canvas.drawPath(path, paint);

    for (int i = 0; i < 8; i++) {
      double angle = (i / 8) * 2 * math.pi;
      double x = math.sin(angle) * 200;
      double y = math.sin(angle) * math.cos(angle) * 100;
      double rotatedX = x * math.cos(state.phaseField) - y * math.sin(state.phaseField);
      double rotatedY = x * math.sin(state.phaseField) + y * math.cos(state.phaseField);
      bool isActive = (state.register & (1 << i)) != 0;
      final nodePaint = Paint()..color = isActive ? const Color(0xFFFFFF00) : const Color(0xFFFF5555)..maskFilter = const MaskFilter.blur(BlurStyle.normal, 10);
      canvas.drawCircle(Offset(center.dx + rotatedX, center.dy + rotatedY), 10, nodePaint);
      canvas.drawCircle(Offset(center.dx + rotatedX, center.dy + rotatedY), 5, Paint()..color = Colors.white);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}

class FloatingFlowWindow extends StatelessWidget {
  final HardwareState state;
  const FloatingFlowWindow({super.key, required this.state});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 250,
      height: 300,
      decoration: BoxDecoration(
        color: const Color(0xCC0F1423),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0x3300FFCC)),
        boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.5), blurRadius: 20)],
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Text('PHOTONIC FLOW: ${state.routingMode.toUpperCase()}', 
              style: const TextStyle(fontSize: 10, fontWeight: FontWeight.bold, color: Color(0xFF00FFCC))),
          ),
          Expanded(
            child: CustomPaint(painter: FlowPainter(state), child: Container()),
          ),
          const Padding(
            padding: EdgeInsets.all(12),
            child: Text('LIVE TELEMETRY', style: TextStyle(fontSize: 8, color: Colors.grey)),
          ),
        ],
      ),
    );
  }
}

class FlowPainter extends CustomPainter {
  final HardwareState state;
  FlowPainter(this.state);

  @override
  void paint(Canvas canvas, Size size) {
    final random = math.Random(42);
    final paint = Paint()..color = const Color(0xFF00FFCC).withOpacity(0.6);
    for (int i = 0; i < 100; i++) {
      double t = DateTime.now().millisecondsSinceEpoch / 1000.0;
      double x, y;
      if (state.routingMode == 'grover') {
        double angle = t * 2 + i * 0.1;
        double r = (i % 20) + 30;
        x = size.width/2 + math.cos(angle) * r;
        y = size.height/2 + math.sin(angle) * r;
      } else {
        x = random.nextDouble() * size.width;
        y = random.nextDouble() * size.height;
      }
      canvas.drawCircle(Offset(x, y), 1.5, paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}

class ClojureVIDE extends StatefulWidget {
  final VoidCallback onClose;
  const ClojureVIDE({super.key, required this.onClose});

  @override
  State<ClojureVIDE> createState() => _ClojureVIDEState();
}

class _ClojureVIDEState extends State<ClojureVIDE> {
  late TextEditingController _controller;
  List<String> _terminal = ['SomaOS Flutter IDE v1.0 initialized.', 'Ready for HPQC synthesis...'];
  bool _isCompiling = false;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(
      text: '(ns ClojureV.qurq)\n\n(defn-ai grover_oracle [clk rst_n in]\n  (let [target 0xABCDEF]\n    (if (= in target)\n      (qurq/phi-scale out in -1.0)\n      (qurq/assign out in))))'
    );
  }

  void _runSynthesis() async {
    print('>> Initiating Live Synthesis via IDE...');
    setState(() {
      _isCompiling = true;
      _terminal.add('> Initiating Live Synthesis...');
    });

    try {
      final response = await http.post(
        Uri.parse('http://localhost:8081/api/synthesize'),
        body: json.encode({'code': _controller.text, 'mode': 'grover'}),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        print('>> [SUCCESS] Synthesis Result: ${data["message"]}');
        setState(() {
          _terminal.add(data['output'] ?? '[SUCCESS] Manifest complete.');
        });
      }
    } catch (e) {
      print('>> [ERROR] Synthesis Failed: $e');
      setState(() => _terminal.add('[ERROR] Toolchain connection failed.'));
    } finally {
      setState(() => _isCompiling = false);
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.black.withOpacity(0.85),
      padding: const EdgeInsets.all(40),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: Scaffold(
          backgroundColor: const Color(0xFF0F1423),
          appBar: AppBar(
            backgroundColor: const Color(0xFF1A2238),
            title: const Text('ClojureV IDE (Flutter Edition)'),
            actions: [
              IconButton(onPressed: _runSynthesis, icon: const Icon(Icons.play_arrow, color: Colors.green)),
              IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
            ],
          ),
          body: Column(
            children: [
              Expanded(
                child: TextField(
                  maxLines: null,
                  controller: _controller,
                  style: const TextStyle(fontSize: 14, color: Colors.white, fontFamily: 'monospace'),
                  decoration: const InputDecoration(contentPadding: EdgeInsets.all(20), border: InputBorder.none),
                ),
              ),
              Container(
                height: 200,
                width: double.infinity,
                color: Colors.black,
                padding: const EdgeInsets.all(15),
                child: ListView.builder(
                  itemCount: _terminal.length,
                  itemBuilder: (context, i) => Text(_terminal[i], style: const TextStyle(color: Colors.green, fontSize: 12)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
