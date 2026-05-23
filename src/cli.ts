const command = process.argv[2] ?? "help";

const commands = new Set([
  "crawl:daily",
  "crawl:weekly",
  "report:daily",
  "report:weekly",
]);

if (command === "help" || command === "--help" || command === "-h") {
  console.log(`Usage: pnpm dev -- <command>

Commands:
  crawl:daily
  crawl:weekly
  report:daily
  report:weekly`);
  process.exit(0);
}

if (!commands.has(command)) {
  console.error(`Unknown command: ${command}`);
  process.exit(1);
}

console.log(`${command} is not implemented yet.`);
