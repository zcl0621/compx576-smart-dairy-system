"""HTTP transport — POST one metric to /agent/api/metric."""
import urequests
import json
import gc
import time


def post_metric(base_url, cow_id, token, metric_type, value, unit):
    """Send one metric. Returns (status_code, response_text)."""
    gc.collect()
    url = base_url.rstrip("/") + "/api/metric"
    payload = {
        "cow_id": cow_id,
        "source": "cow_agent",
        "metric_type": metric_type,
        "metric_value": float(value),
        "unit": unit,
    }
    headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + token,
    }
    r = None
    try:
        r = urequests.post(url, data=json.dumps(payload), headers=headers)
        code = r.status_code
        body = r.text
        return code, body
    finally:
        if r is not None:
            r.close()
        gc.collect()


def send_report(base_url, cow_id, token, report):
    """Send all non-null fields from a 30s aggregate report."""
    gc.collect()
    sent = 0
    failed = 0

    fields = []
    if report.get("temp_c") is not None:
        fields.append(("temperature", report["temp_c"], "celsius"))
    if report.get("bpm") is not None:
        fields.append(("heart_rate", report["bpm"], "bpm"))
    if report.get("spo2") is not None:
        fields.append(("blood_oxygen", report["spo2"], "percent"))
    if report.get("lat") is not None:
        fields.append(("latitude", report["lat"], "degrees"))
    if report.get("lon") is not None:
        fields.append(("longitude", report["lon"], "degrees"))

    for metric_type, value, unit in fields:
        try:
            code, body = post_metric(base_url, cow_id, token, metric_type, value, unit)
            if 200 <= code < 300:
                sent += 1
            else:
                failed += 1
                print("post {} -> {} {}".format(metric_type, code, body))
        except Exception as e:
            failed += 1
            print("post {} error: {} args={}".format(
                metric_type, type(e).__name__, e.args))
        # free TLS buffers between requests
        gc.collect()
        time.sleep_ms(200)

    return sent, failed
