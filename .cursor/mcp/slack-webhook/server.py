import json
import os
import urllib.error
import urllib.request

from mcp.server.fastmcp import FastMCP

mcp = FastMCP("Slack Webhook")


@mcp.tool()
def send_slack_message(text: str, blocks: str = "") -> str:
    """Send a message to a Slack channel via incoming webhook.

    Args:
        text: The message text (used as fallback and notification text).
        blocks: Optional JSON string of Slack Block Kit blocks for rich formatting.
    """
    webhook_url = os.environ.get("SLACK_WEBHOOK_URL")
    if not webhook_url:
        return "Error: SLACK_WEBHOOK_URL environment variable is not set"

    payload = {"text": text}
    if blocks:
        try:
            payload["blocks"] = json.loads(blocks)
        except json.JSONDecodeError:
            return "Error: blocks parameter is not valid JSON"

    req = urllib.request.Request(
        webhook_url,
        data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )

    try:
        with urllib.request.urlopen(req) as resp:
            return f"Message sent successfully (HTTP {resp.status})"
    except urllib.error.HTTPError as e:
        return f"Error sending message: HTTP {e.code} - {e.read().decode()}"
    except urllib.error.URLError as e:
        return f"Error sending message: {e.reason}"


if __name__ == "__main__":
    mcp.run(transport="stdio")
