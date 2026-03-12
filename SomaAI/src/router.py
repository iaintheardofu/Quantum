import asyncio
import json
import base64
from fastapi import FastAPI, WebSocket, Request
import uvicorn
import vertexai
from vertexai.generative_models import GenerativeModel, Part
from contextlib import asynccontextmanager

# We can also expose this as an MCP server using fastmcp
from mcp.server.fastmcp import FastMCP

# Initialize Vertex AI with the authenticated project
PROJECT_ID = "mabel-488406"
LOCATION = "us-central1"
vertexai.init(project=PROJECT_ID, location=LOCATION)

model = GenerativeModel("gemini-1.5-pro-preview-0409")

# Create the MCP Server instance
mcp = FastMCP("SomaAI_Cortex")

# Create the FastAPI app
app = FastAPI(title="SomaAI Cortex Router")

@app.get("/")
async def root():
    return {"status": "SomaAI Router is Online. Vertex AI Connected."}

# --- MCP Tools & Endpoints ---

@mcp.tool()
async def analyze_telemetry_log(log_text: str) -> str:
    """Sends a chunk of hardware logs to Gemini for analysis."""
    prompt = f"""
    You are the SomaAI Cortex, analyzing real-time telemetry from the MABEL x8C quantum virtualizer.
    Review the following system logs. 
    1. Identify any signs of thermal decoherence (rapid temperature spikes).
    2. Confirm if the routing mode synthesis was successful.
    3. Look for any 'Silence Protocol' transmission errors.
    
    Keep your analysis under 3 sentences. Be highly technical.
    
    LOGS:
    {log_text}
    """
    try:
        response = await model.generate_content_async(prompt)
        return response.text
    except Exception as e:
        return f"Vertex AI Error: {e}"

@mcp.tool()
async def analyze_hardware_audio(base64_audio: str) -> str:
    """Analyzes the thermodynamic resonance (coil whine) of the ALINX board."""
    try:
        audio_bytes = base64.b64decode(base64_audio)
        audio_part = Part.from_data(data=audio_bytes, mime_type="audio/wav")
        prompt = "Listen to this audio recording of the FPGA's SPHY Engine. Do you hear any high-pitch variations or stuttering that might indicate a failure in the stochastic compensation loop?"
        response = await model.generate_content_async([audio_part, prompt])
        return response.text
    except Exception as e:
        return f"Vertex AI Audio Analysis Error: {e}"

@mcp.tool()
async def analyze_hardware_video(base64_image: str) -> str:
    """Analyzes a webcam frame of the ALINX board diagnostic LEDs."""
    try:
        image_bytes = base64.b64decode(base64_image)
        image_part = Part.from_data(data=image_bytes, mime_type="image/jpeg")
        prompt = "Look at this image of the ALINX 7020 board. Are the 'DONE' and 'INIT_B' LEDs illuminated, indicating a successful PL configuration?"
        response = await model.generate_content_async([image_part, prompt])
        return response.text
    except Exception as e:
        return f"Vertex AI Video Analysis Error: {e}"

@mcp.tool()
async def generate_clojurev(intent_description: str) -> str:
    """Generates ClojureV topological code based on a user's prompt."""
    prompt = f"""
    You are an expert in SomaOS and the ClojureV hardware description language.
    Write a ClojureV function (`defn-ai` or `defn-fractal`) to fulfill the following intent.
    Output ONLY the valid ClojureV code.

    Intent: {intent_description}
    """
    try:
        response = await model.generate_content_async(prompt)
        return response.text
    except Exception as e:
        return f";; Error generating code: {e}"


# --- FastAPI Routes (For the Flutter UI) ---

@app.post("/api/ai/telemetry")
async def api_telemetry(request: Request):
    data = await request.json()
    logs = data.get("logs", "")
    analysis = await analyze_telemetry_log(logs)
    return {"insight": analysis}

@app.post("/api/ai/generate_code")
async def api_generate_code(request: Request):
    data = await request.json()
    intent = data.get("intent", "")
    code = await generate_clojurev(intent)
    return {"code": code}

@app.post("/api/ai/vision")
async def api_vision(request: Request):
    data = await request.json()
    base64_image = data.get("image", "")
    prompt_text = data.get("prompt", "Analyze this image of the ALINX 7020 board. What is the status of the diagnostic LEDs?")
    
    if "base64," in base64_image:
        base64_image = base64_image.split("base64,")[1]

    try:
        image_bytes = base64.b64decode(base64_image)
        image_part = Part.from_data(data=image_bytes, mime_type="image/jpeg")
        response = await model.generate_content_async([image_part, prompt_text])
        return {"insight": response.text}
    except Exception as e:
        return {"insight": f"Vertex AI Vision Error: {e}"}

@app.websocket("/ws/telemetry_analysis")
async def websocket_telemetry_analysis(websocket: WebSocket):
    await websocket.accept()
    log_buffer = []
    
    try:
        while True:
            data = await websocket.receive_text()
            log_buffer.append(data)
            
            # Analyze every 20 log lines
            if len(log_buffer) >= 20:
                chunk = "\n".join(log_buffer)
                analysis = await analyze_telemetry_log(chunk)
                await websocket.send_json({"type": "ai_insight", "data": analysis})
                log_buffer = [] # Clear buffer
    except Exception as e:
        print(f"WebSocket Error: {e}")

# Note: In a real production deployment, FastMCP and FastAPI can be bridged.
# For this script, we launch the FastAPI server which the Flutter UI relies on.
# The MCP server can be run separately via `mcp run SomaAI/src/router.py:mcp`

if __name__ == "__main__":
    print(">> Starting SomaAI Cortex Router on port 8083...")
    uvicorn.run(app, host="0.0.0.0", port=8083)
