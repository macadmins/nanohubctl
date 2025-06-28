# nanohubctl
nanohubctl is a tool for interacting with nanohub.

## Usage
You can pass in the API key and URL every time you run the command or set env vars for them instead.

To set th env vars:
```bash
# This should be the base url to nanohub, Example: https://nanohub.example.com/
export DDM_URL="https://$YOUR_SERVER_HERE"
export DDM_API_KEY="$YOUR_API_KEY"
# Optional, but handy if you are only working with a single client
export DDM_CLIENT_ID="$TEST_CLIENT_ID"
```
