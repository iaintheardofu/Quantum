import 'dart:async';
import 'dart:convert';
import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

// Central Logger for SomaOS
class SomaLogger {
  static WebSocketChannel? _channel;
  static bool _isConnected = false;

  static void _init() {
    if (_isConnected) return;
    try {
      _channel = WebSocketChannel.connect(Uri.parse('ws://localhost:8082/log'));
      _isConnected = true;
    } catch (e) {
      _isConnected = false;
    }
  }

  static void log(String message) {
    if (!_isConnected) _init();
    try {
      _channel?.sink.add(message);
    } catch (e) {
      _isConnected = false;
    }
  }
}

void main() {
  runZonedGuarded(() {
    WidgetsFlutterBinding.ensureInitialized();
    runApp(const SomaApp());
  }, (error, stackTrace) {
    SomaLogger.log('[BROWSER-ERROR] $error');
  }, zoneSpecification: ZoneSpecification(
    print: (Zone self, ZoneDelegate parent, Zone zone, String line) {
      parent.print(zone, line);
      SomaLogger.log(line);
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

// --- UPDATED IDE WITH SENTINEL EYE & CORTEX AI ---

// --- EXAMPLES DATA ---
const Map<String, String> EXAMPLES = {
  'station_scaling.cljv': '''(ns ClojureV.qurq)

(defn-fractal HyperStation [clk rst_n in]
  "Manifesting a 64-qubit Fractal Hypercube via the Master Station Hub"
  (let [master_anchor (qurq/spawn-macro-cell :Station_Core :anchor true)
        station_bus (qurq/spawn-station-bus master_anchor)]
    
    ;; Braiding 8 independent 8-qubit cubes into a unified manifold
    (loop [cube_idx 0]
      (if (< cube_idx 8)
        (do
          (qurq/spawn-macro-cube cube_idx :connect-to station_bus)
          (recur (inc cube_idx)))))
    
    (qurq/assign out station_bus)))''',

  'grovers_search.cljv': '''(ns ClojureV.qurq)

(defn-ai grover_oracle [clk rst_n in]
  (let [target 0xABCDEF]
    (if (= in target)
      (qurq/phi-scale out in -1.0)
      (qurq/assign out in))))

(defn-ai grover_diffusion [clk rst_n in]
  (let [mean (qurq/read-average in)]
    (qurq/phi-scale out (- (* 2 mean) in))))

(defn-fractal grovers_search [clk rst_n in]
  (loop [depth 16 signal in]
    (if (zero? depth)
      (qurq/assign out signal)
      (let [marked (grover_oracle clk rst_n signal)]
        (recur (dec depth) (grover_diffusion clk rst_n marked))))))''',

  'shors_factorization.cljv': '''(ns ClojureV.qurq)

(defn-ai modular_exponentiation [clk rst_n base exp mod]
  (let [phi_base (qurq/phi-scale base)]
    (qurq/mod-exp out phi_base exp mod)))

(defn-fractal quantum_fourier_transform [clk rst_n in]
  (loop [q_idx 0 data in]
    (if (= q_idx 8)
      (qurq/assign out data)
      (let [h_gate (qurq/hadamard data q_idx)
            cp_gate (qurq/controlled-phase h_gate q_idx)]
        (recur (inc q_idx) cp_gate)))))

(defn-fractal shors_factorization [clk rst_n n]
  (let [a (qurq/select-coprime n)
        x (qurq/superposition 8)
        f_x (modular_exponentiation clk rst_n a x n)]
    (let [qft_result (quantum_fourier_transform clk rst_n f_x)]
      (qurq/collapsed-period out qft_result))))''',

  'bell_state.cljv': '''(ns ClojureV.qurq)

(defn-ai BellState [clk rst_n in]
  "Creating a maximally entangled state (Psi+) using a sum-split braid"
  (let [h_gate (qurq/hadamard in)
        cnot_gate (qurq/sum-split h_gate in)]
    (qurq/assign out cnot_gate)))'''
};


class ClojureVIDE extends StatefulWidget {
  final VoidCallback onClose;
  const ClojureVIDE({super.key, required this.onClose});

  @override
  State<ClojureVIDE> createState() => _ClojureVIDEState();
}

class _ClojureVIDEState extends State<ClojureVIDE> {
  late TextEditingController _controller;
  final ScrollController _terminalScrollController = ScrollController();
  List<String> _terminal = ['SomaOS Flutter IDE v1.0 initialized.', 'Ready for HPQC synthesis...'];
  bool _isCompiling = false;
  String _activeFile = 'grovers_search.cljv';
  
  // Sentinel Eye Stream
  WebSocketChannel? _eyeChannel;
  String _latestFrameBase64 = "";

  // Cortex AI Chat
  final TextEditingController _chatController = TextEditingController();
  List<String> _chatLog = ['[CORTEX AI] Connected to Vertex AI Multimodal Engine.'];
  bool _isThinking = false;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: EXAMPLES[_activeFile]);
    _connectSentinelEye();
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_terminalScrollController.hasClients) {
        _terminalScrollController.animateTo(
          _terminalScrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _runSynthesis() async {
    print('>> Initiating Live Synthesis via IDE...');
    setState(() {
      _isCompiling = true;
      _terminal.add('> Initiating Live Synthesis for $_activeFile...');
    });
    _scrollToBottom();

    // Determine routing mode based on active file
    String mode = 'idle';
    if (_activeFile.contains('grover')) mode = 'grover';
    else if (_activeFile.contains('shor')) mode = 'shor';
    else if (_activeFile.contains('bell')) mode = 'bell';
    else if (_activeFile.contains('station')) mode = 'station';

    try {
      final response = await http.post(
        Uri.parse('http://localhost:8081/api/synthesize'),
        body: json.encode({'code': _controller.text, 'mode': mode}),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        setState(() {
          if (data['output'] != null) {
            _terminal.addAll(data['output'].toString().split('\n').where((l) => l.trim().isNotEmpty));
          }
          _terminal.add('[SUCCESS] Manifest complete for mode: ${mode.toUpperCase()}');
        });
      }
    } catch (e) {
      setState(() => _terminal.add('[ERROR] Toolchain connection failed.'));
    } finally {
      setState(() => _isCompiling = false);
      _scrollToBottom();
    }
  }

  void _selectFile(String filename) {
    setState(() {
      _activeFile = filename;
      _controller.text = EXAMPLES[filename] ?? '';
      _terminal.add('> Opened $filename');
    });
    _scrollToBottom();
  }

  void _connectSentinelEye() {
    try {
      _eyeChannel = WebSocketChannel.connect(Uri.parse('ws://localhost:5001'));
      _eyeChannel!.stream.listen((message) {
        final data = json.decode(message);
        if (data['event'] == 'photonic_frame') {
          if (mounted) {
            setState(() {
              _latestFrameBase64 = data['data'];
            });
          }
        }
      }, onError: (e) {
        print("[EYE] Connection error: $e");
      });
    } catch (e) {
      print("[EYE] Failed to connect: $e");
    }
  }

  void _askCortex() async {
    final query = _chatController.text;
    if (query.isEmpty) return;

    setState(() {
      _chatLog.add('> You: $query');
      _chatController.clear();
      _isThinking = true;
    });

    try {
      final response = await http.post(
        Uri.parse('http://localhost:8083/api/ai/vision'),
        body: json.encode({'image': _latestFrameBase64, 'prompt': query}),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        setState(() {
          _chatLog.add('[CORTEX] ${data["insight"]}');
        });
      } else {
        setState(() => _chatLog.add('[CORTEX ERROR] Failed to analyze vision.'));
      }
    } catch (e) {
      setState(() => _chatLog.add('[CORTEX ERROR] Router disconnected.'));
    } finally {
      setState(() => _isThinking = false);
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    _chatController.dispose();
    _eyeChannel?.sink.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.black.withOpacity(0.85),
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          // LEFT: IDE Sidebar + Editor + Terminal
          Expanded(
            flex: 2,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(12),
              child: Scaffold(
                backgroundColor: const Color(0xFF0F1423),
                appBar: AppBar(
                  backgroundColor: const Color(0xFF1A2238),
                  title: const Text('ClojureV Forge', style: TextStyle(fontFamily: 'monospace')),
                  actions: [
                    IconButton(onPressed: _runSynthesis, icon: const Icon(Icons.play_arrow, color: Colors.green)),
                    IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
                  ],
                ),
                body: Row(
                  children: [
                    // Examples Sidebar
                    Container(
                      width: 200,
                      color: const Color(0xFF12182B),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Container(
                            padding: const EdgeInsets.all(10),
                            color: Colors.black,
                            width: double.infinity,
                            child: const Text("EXAMPLES", style: TextStyle(color: Colors.grey, fontSize: 10, letterSpacing: 1, fontWeight: FontWeight.bold)),
                          ),
                          Expanded(
                            child: ListView(
                              children: EXAMPLES.keys.map((filename) => InkWell(
                                onTap: () => _selectFile(filename),
                                child: Container(
                                  padding: const EdgeInsets.all(10),
                                  color: _activeFile == filename ? const Color(0xFF1A2238) : Colors.transparent,
                                  child: Text(filename, style: TextStyle(color: _activeFile == filename ? Colors.cyan : Colors.grey, fontSize: 12, fontFamily: 'monospace')),
                                ),
                              )).toList(),
                            ),
                          ),
                        ],
                      ),
                    ),
                    
                    // Editor & Terminal Column
                    Expanded(
                      child: Column(
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
                            height: 150,
                            width: double.infinity,
                            color: Colors.black,
                            padding: const EdgeInsets.all(10),
                            child: SelectionArea(
                              child: ListView.builder(
                                controller: _terminalScrollController,
                                itemCount: _terminal.length,
                                itemBuilder: (context, i) => Text(_terminal[i], style: const TextStyle(color: Colors.green, fontSize: 12, fontFamily: 'monospace')),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
          
          const SizedBox(width: 20),

          // RIGHT: Cortex AI & Sentinel Eye
          Expanded(
            flex: 1,
            child: Column(
              children: [
                // Sentinel Eye Feed
                Container(
                  height: 200,
                  width: double.infinity,
                  decoration: BoxDecoration(
                    color: Colors.black,
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: const Color(0xFF00FFCC)),
                  ),
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: _latestFrameBase64.isEmpty 
                      ? const Center(child: Text("NO SIGNAL\nWaiting for Sentinel Eye...", textAlign: TextAlign.center, style: TextStyle(color: Colors.red)))
                      : Image.memory(
                          base64Decode(_latestFrameBase64.split(',').last),
                          fit: BoxFit.cover,
                          gaplessPlayback: true,
                        ),
                  ),
                ),
                const SizedBox(height: 10),
                const Text("SENTINEL EYE: PHYSICAL OBSERVATION", style: TextStyle(color: Color(0xFF00FFCC), fontSize: 10, fontWeight: FontWeight.bold)),
                
                const SizedBox(height: 20),

                // Cortex Chat
                Expanded(
                  child: Container(
                    decoration: BoxDecoration(
                      color: const Color(0xFF1A2238),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      children: [
                        Container(
                          padding: const EdgeInsets.all(10),
                          width: double.infinity,
                          decoration: const BoxDecoration(
                            border: Border(bottom: BorderSide(color: Colors.black)),
                          ),
                          child: const Text("CORTEX AI: MULTIMODAL AUDIT", style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 12)),
                        ),
                        Expanded(
                          child: ListView.builder(
                            padding: const EdgeInsets.all(10),
                            itemCount: _chatLog.length,
                            itemBuilder: (context, i) {
                              bool isUser = _chatLog[i].startsWith('> You');
                              return Padding(
                                padding: const EdgeInsets.only(bottom: 8.0),
                                child: Text(_chatLog[i], style: TextStyle(color: isUser ? Colors.white : Colors.cyan, fontSize: 12)),
                              );
                            }
                          ),
                        ),
                        if (_isThinking) const Padding(
                          padding: EdgeInsets.all(8.0),
                          child: LinearProgressIndicator(color: Colors.cyan),
                        ),
                        Padding(
                          padding: const EdgeInsets.all(10.0),
                          child: Row(
                            children: [
                              Expanded(
                                child: TextField(
                                  controller: _chatController,
                                  style: const TextStyle(fontSize: 12),
                                  decoration: const InputDecoration(
                                    hintText: "Ask Cortex about the hardware...",
                                    filled: true,
                                    fillColor: Colors.black,
                                    border: OutlineInputBorder(borderSide: BorderSide.none),
                                    contentPadding: EdgeInsets.symmetric(horizontal: 10, vertical: 0)
                                  ),
                                  onSubmitted: (_) => _askCortex(),
                                ),
                              ),
                              IconButton(
                                icon: const Icon(Icons.send, color: Colors.cyan),
                                onPressed: _askCortex,
                              )
                            ],
                          ),
                        ),
                      ],
                    ),
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
