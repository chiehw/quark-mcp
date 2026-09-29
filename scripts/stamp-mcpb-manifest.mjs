import fs from "node:fs";

const [src, dest, version, platform, binName] = process.argv.slice(2);
if (!src || !dest || !version || !platform || !binName) {
  console.error("usage: stamp-mcpb-manifest.mjs <src> <dest> <version> <platform> <binName>");
  process.exit(1);
}

const manifest = JSON.parse(fs.readFileSync(src, "utf8"));
manifest.version = version;
manifest.compatibility = manifest.compatibility || {};
manifest.compatibility.platforms = [platform];
manifest.server.entry_point = `server/${binName}`;
manifest.server.mcp_config = manifest.server.mcp_config || {};
manifest.server.mcp_config.command = `\${__dirname}/server/${binName}`;
fs.writeFileSync(dest, `${JSON.stringify(manifest, null, 2)}\n`);
