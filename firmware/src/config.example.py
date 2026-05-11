"""Device-side config template. Copy to config.py and fill in real values.

Get cow_id and agent_token from the admin UI or directly from the DB:
    select id, name, tag, agent_token from cows where tag = 'DC-2406';
"""

# --- WiFi ---
WIFI_SSID = "your-wifi-ssid"
WIFI_PASSWORD = "your-wifi-password"

# --- Backend ---
# Ingress maps /agent/* on the host to agent-server:8081
AGENT_BASE_URL = "https://smart-farm.app/agent"
COW_ID = "your-cow-id"
AGENT_TOKEN = "your-agent-token"
