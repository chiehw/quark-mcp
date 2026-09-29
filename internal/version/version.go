package version

// Version is the quark-mcp release version. Release builds override it with
// -ldflags from the Git tag so the binary, MCP handshake, and MCPB manifest match.
var Version = "1.3.0"
