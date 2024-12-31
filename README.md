# SuDeploy
easy Systemd Unit Deploy tools

# Install

`go install github.com/zjyl1994/sudeploy@latest`

# Use

1. Copy `sudeploy.example.json` in your work dir.
2. Rename `sudeploy.example.json` to `sudeploy.json`.
3. Change content in `sudeploy.json`.
4. Run `sudeploy`,it will deploy you binary to remote server.
5. If remote service is running, it only change binary.

# FAQ
1. If you want wait service start, set `wait_seconds` in config file.
2. Upload file only run in service install.
3. Upgrade binary only change binary file, systemd unit will not change.