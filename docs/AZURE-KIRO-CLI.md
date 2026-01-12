# Kiro CLI Integration

Use Kiro CLI with wetwire-azure for AI-assisted infrastructure design in corporate Azure environments.

## Prerequisites

- Go 1.23+ installed
- Kiro CLI installed ([installation guide](https://kiro.dev/docs/cli/installation/))
- AWS Builder ID or GitHub/Google account (for Kiro authentication)

---

## Step 1: Install wetwire-azure

### Option A: Using Go (recommended)

```bash
go install github.com/lex00/wetwire-azure-go/cmd/wetwire-azure@latest
```

### Option B: Pre-built binaries

Download from [GitHub Releases](https://github.com/lex00/wetwire-azure-go/releases):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/lex00/wetwire-azure-go/releases/latest/download/wetwire-azure-darwin-arm64
chmod +x wetwire-azure-darwin-arm64
sudo mv wetwire-azure-darwin-arm64 /usr/local/bin/wetwire-azure

# macOS (Intel)
curl -LO https://github.com/lex00/wetwire-azure-go/releases/latest/download/wetwire-azure-darwin-amd64
chmod +x wetwire-azure-darwin-amd64
sudo mv wetwire-azure-darwin-amd64 /usr/local/bin/wetwire-azure

# Linux (x86-64)
curl -LO https://github.com/lex00/wetwire-azure-go/releases/latest/download/wetwire-azure-linux-amd64
chmod +x wetwire-azure-linux-amd64
sudo mv wetwire-azure-linux-amd64 /usr/local/bin/wetwire-azure
```

### Verify installation

```bash
wetwire-azure --version
```

---

## Step 2: Install Kiro CLI

```bash
# Install Kiro CLI
curl -fsSL https://cli.kiro.dev/install | bash

# Verify installation
kiro-cli --version

# Sign in (opens browser)
kiro-cli login
```

---

## Step 3: Configure Kiro for wetwire-azure

Run the design command with `--provider kiro` to auto-configure:

```bash
# Create a project directory
mkdir my-infra && cd my-infra

# Initialize Go module
go mod init my-infra

# Run design with Kiro provider (auto-installs configs on first run)
wetwire-azure design --provider kiro "Create a storage account"
```

This automatically installs:

| File | Purpose |
|------|---------|
| `~/.kiro/agents/wetwire-azure-runner.json` | Kiro agent configuration |
| `.kiro/mcp.json` | Project MCP server configuration |

### Manual configuration (optional)

The MCP server is provided as a subcommand `wetwire-azure mcp`. If you prefer to configure manually:

**~/.kiro/agents/wetwire-azure-runner.json:**
```json
{
  "name": "wetwire-azure-runner",
  "description": "Infrastructure code generator using wetwire-azure",
  "prompt": "You are an infrastructure design assistant...",
  "model": "claude-sonnet-4",
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-azure",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  },
  "tools": ["*"]
}
```

**.kiro/mcp.json:**
```json
{
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-azure",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  }
}
```

> **Note:** The `cwd` field ensures MCP tools resolve paths correctly in your project directory. When using `wetwire-azure design --provider kiro`, this is configured automatically.

---

## Step 4: Run Kiro with wetwire design

### Using the wetwire-azure CLI

```bash
# Start Kiro design session
wetwire-azure design --provider kiro "Create a serverless API with Azure Functions and Cosmos DB"
```

This launches Kiro CLI with the wetwire-azure-runner agent and your prompt.

### Using Kiro CLI directly

```bash
# Start chat with wetwire-azure-runner agent
kiro-cli chat --agent wetwire-azure-runner

# Or with an initial prompt
kiro-cli chat --agent wetwire-azure-runner "Create a storage account with geo-redundant storage"
```

---

## Available MCP Tools

The wetwire-azure MCP server exposes three tools to Kiro:

| Tool | Description | Example |
|------|-------------|---------|
| `build` | Build ARM template from Go resource definitions | `build(path="./myapp")` |
| `lint` | Lint code for issues | `lint(path="./infra/...")` |
| `import` | Convert ARM template JSON to Go code | `import(file="template.json", package="infra")` |

---

## Example Session

```
$ wetwire-azure design --provider kiro "Create a storage account with geo-redundant storage and encryption"

Installed Kiro agent config: ~/.kiro/agents/wetwire-azure-runner.json
Installed project MCP config: .kiro/mcp.json
Starting Kiro CLI design session...

> I'll help you create a storage account with geo-redundant storage and encryption enabled.

Let me initialize the project and create the infrastructure code.

[Calling build...]
[Calling lint...]
[Calling build...]

I've created the following files:
- infra/storage.go

The storage account includes:
- Geo-redundant storage (Standard_GRS)
- Encryption enabled
- Secure transfer required

Would you like me to add any additional configurations?
```

---

## Workflow

The Kiro agent follows this workflow:

1. **Explore** - Understand your requirements
2. **Plan** - Design the infrastructure architecture
3. **Implement** - Generate Go code using wetwire-azure patterns
4. **Lint** - Run `lint` to check for issues
5. **Build** - Run `build` to generate ARM template

---

## Deploying Generated Templates

After Kiro generates your infrastructure code:

```bash
# Build the ARM template
wetwire-azure build ./infra > template.json

# Deploy with Azure CLI
az deployment group create \
  --resource-group my-resource-group \
  --template-file template.json

# Or deploy to subscription level
az deployment sub create \
  --location eastus \
  --template-file template.json
```

---

## Troubleshooting

### MCP server not found

```
Mcp error: -32002: No such file or directory
```

**Solution:** Ensure `wetwire-azure` is in your PATH:

```bash
which wetwire-azure

# If not found, add to PATH or reinstall
go install github.com/lex00/wetwire-azure-go/cmd/wetwire-azure@latest
```

### Kiro CLI not found

```
kiro-cli not found in PATH
```

**Solution:** Install Kiro CLI:

```bash
curl -fsSL https://cli.kiro.dev/install | bash
```

### Authentication issues

```
Error: Not authenticated
```

**Solution:** Sign in to Kiro:

```bash
kiro-cli login
```

---

## Known Limitations

### Automated Testing

When using `wetwire-azure test --provider kiro`, tests run in non-interactive mode (`--no-interactive`). This means:

- The agent runs autonomously without waiting for user input
- Persona simulation is limited - all personas behave similarly
- The agent won't ask clarifying questions

For true persona simulation with multi-turn conversations, use the Anthropic provider:

```bash
wetwire-azure test --provider anthropic --persona expert "Create a storage account"
```

### Interactive Design Mode

Interactive design mode (`wetwire-azure design --provider kiro`) works fully as expected:

- Real-time conversation with the agent
- Agent can ask clarifying questions
- Lint loop executes as specified in the agent prompt

---

## See Also

- [CLI Reference](CLI.md) - Full wetwire-azure CLI documentation
- [Quick Start](QUICK_START.md) - Getting started with wetwire-azure
- [Kiro CLI Installation](https://kiro.dev/docs/cli/installation/) - Official installation guide
- [Kiro CLI Docs](https://kiro.dev/docs/cli/) - Official Kiro documentation
