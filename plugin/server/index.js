'use strict';

// Minimal Model Context Protocol server over stdio, with no dependencies.
// Speaks JSON-RPC 2.0 framed as newline-delimited JSON on stdin/stdout.
// stdout carries protocol messages only; all logging goes to stderr.

const { buildTools, loadConfig, ToolInputError } = require('./tools');
const { buildPrompts } = require('./prompts');

const SERVER_NAME = 'eurohop-flight-planner';
const SERVER_VERSION = '1.0.1';

// Protocol revisions this server knows how to speak. We echo the client's
// version when we support it, and otherwise offer our newest.
const SUPPORTED_PROTOCOL_VERSIONS = ['2025-06-18', '2025-03-26', '2024-11-05'];
const LATEST_PROTOCOL_VERSION = SUPPORTED_PROTOCOL_VERSIONS[0];

const INSTRUCTIONS = `Flight search across Ryanair's European point-to-point network, using public FareFinder data.

This is an unofficial, independent community tool. It is not affiliated with, authorised by, or endorsed by Ryanair DAC or any of its subsidiaries. It reads publicly accessible fare data and cannot book, change, or cancel anything. Always send the traveller to ryanair.com to book, and tell them fares change constantly, so a price shown here may differ by the time they book.

Picking a tool:
- plan_trip: the traveller named places, not airports ("Brussels to Andalusia"). It expands both ends into nearby airports and searches every pairing.
- find_nearby_airports: turn a place name into airport codes, or find cheaper secondary airports in range.
- search_flights / search_return_flights: both airport codes are already known.
- get_daily_flights: a route is chosen and the traveller needs actual departure times. The other tools return only the cheapest flight per day.

Worth asking before a broad search, since each materially changes the result: whether dates are flexible by a day or two (fares vary widely between adjacent days); whether the return must use the same airport as the outbound, or an open-jaw route is acceptable; and which currency to price in.

When reporting results, give the distance from each airport to where the traveller actually wants to be. A fare that is cheaper from an airport 200km away is often not the better trip.`;

function log(...parts) {
  process.stderr.write(`[${SERVER_NAME}] ${parts.join(' ')}\n`);
}

function send(message) {
  process.stdout.write(JSON.stringify(message) + '\n');
}

function sendResult(id, result) {
  send({ jsonrpc: '2.0', id, result });
}

function sendError(id, code, message) {
  send({ jsonrpc: '2.0', id, error: { code, message } });
}

// Read only the three named variables this server uses, each one the value of
// a userConfig option declared in plugin.json. The environment is never read
// as a whole, so nothing else on the machine can reach the network through us.
const config = loadConfig({
  currency: process.env.CLAUDE_PLUGIN_OPTION_DEFAULT_CURRENCY,
  origin: process.env.CLAUDE_PLUGIN_OPTION_DEFAULT_ORIGIN,
  radiusKm: process.env.CLAUDE_PLUGIN_OPTION_DEFAULT_RADIUS_KM,
});
const tools = buildTools(config);
const prompts = buildPrompts(config);

const toolsByName = new Map(tools.map((tool) => [tool.name, tool]));
const promptsByName = new Map(prompts.map((prompt) => [prompt.name, prompt]));

function listedTools() {
  return tools.map((tool) => ({
    name: tool.name,
    title: tool.title,
    description: tool.description,
    inputSchema: tool.inputSchema,
    annotations: tool.annotations,
  }));
}

function listedPrompts() {
  return prompts.map((prompt) => ({
    name: prompt.name,
    title: prompt.title,
    description: prompt.description,
    arguments: prompt.arguments,
  }));
}

async function callTool(params) {
  const name = params && params.name;
  const tool = toolsByName.get(name);

  if (!tool) {
    const known = Array.from(toolsByName.keys()).join(', ');
    return {
      content: [{ type: 'text', text: `Unknown tool "${name}". Available tools: ${known}.` }],
      isError: true,
    };
  }

  const args = (params && params.arguments) || {};
  if (typeof args !== 'object' || Array.isArray(args)) {
    return {
      content: [{ type: 'text', text: `Arguments for "${name}" must be a JSON object.` }],
      isError: true,
    };
  }

  try {
    const result = await tool.handler(args);
    return { content: [{ type: 'text', text: JSON.stringify(result, null, 2) }] };
  } catch (err) {
    const detail = err instanceof Error ? err.message : String(err);
    const prefix = err instanceof ToolInputError ? 'Invalid input' : `${name} failed`;
    if (!(err instanceof ToolInputError)) log(`${name} error:`, detail);
    return { content: [{ type: 'text', text: `${prefix}: ${detail}` }], isError: true };
  }
}

function getPrompt(params) {
  const name = params && params.name;
  const prompt = promptsByName.get(name);

  if (!prompt) {
    const known = Array.from(promptsByName.keys()).join(', ');
    throw new ToolInputError(`Unknown prompt "${name}". Available prompts: ${known}.`);
  }

  const args = (params && params.arguments) || {};
  for (const argument of prompt.arguments) {
    if (argument.required && !args[argument.name]) {
      throw new ToolInputError(`Prompt "${name}" requires the "${argument.name}" argument.`);
    }
  }

  const built = prompt.build(args);
  return {
    description: built.description,
    messages: [{ role: 'user', content: { type: 'text', text: built.text } }],
  };
}

async function handleRequest(method, params) {
  switch (method) {
    case 'initialize': {
      const requested = params && params.protocolVersion;
      return {
        protocolVersion: SUPPORTED_PROTOCOL_VERSIONS.includes(requested)
          ? requested
          : LATEST_PROTOCOL_VERSION,
        capabilities: { tools: {}, prompts: {} },
        serverInfo: { name: SERVER_NAME, title: 'EuroHop Flight Planner', version: SERVER_VERSION },
        instructions: INSTRUCTIONS,
      };
    }
    case 'ping':
      return {};
    case 'tools/list':
      return { tools: listedTools() };
    case 'tools/call':
      return await callTool(params);
    case 'prompts/list':
      return { prompts: listedPrompts() };
    case 'prompts/get':
      return getPrompt(params);
    case 'resources/list':
      return { resources: [] };
    case 'resources/templates/list':
      return { resourceTemplates: [] };
    default:
      return undefined; // Signals "method not found" to the caller.
  }
}

async function handleMessage(message) {
  if (!message || message.jsonrpc !== '2.0' || typeof message.method !== 'string') {
    if (message && message.id !== undefined && message.id !== null) {
      sendError(message.id, -32600, 'Invalid JSON-RPC 2.0 request.');
    }
    return;
  }

  const isNotification = message.id === undefined || message.id === null;

  // Responses to requests we never make, and notifications, need no reply.
  if (isNotification) return;

  try {
    const result = await handleRequest(message.method, message.params);
    if (result === undefined) {
      sendError(message.id, -32601, `Method not found: ${message.method}`);
      return;
    }
    sendResult(message.id, result);
  } catch (err) {
    const detail = err instanceof Error ? err.message : String(err);
    const code = err instanceof ToolInputError ? -32602 : -32603;
    if (!(err instanceof ToolInputError)) log(`${message.method} error:`, detail);
    sendError(message.id, code, detail);
  }
}

function main() {
  let buffer = '';
  let inFlight = 0;
  let stdinEnded = false;

  // Exit only once stdin has closed AND every request we accepted has replied,
  // so a client that closes the stream mid-call still receives its answers.
  function exitWhenDrained() {
    if (stdinEnded && inFlight === 0) process.exit(0);
  }

  function track(promise) {
    inFlight += 1;
    promise
      .catch((err) => log('unhandled error:', err && err.message))
      .finally(() => {
        inFlight -= 1;
        exitWhenDrained();
      });
  }

  process.stdin.setEncoding('utf8');

  process.stdin.on('data', (chunk) => {
    buffer += chunk;

    let newlineIndex = buffer.indexOf('\n');
    while (newlineIndex !== -1) {
      const line = buffer.slice(0, newlineIndex).trim();
      buffer = buffer.slice(newlineIndex + 1);

      if (line !== '') {
        let parsed;
        try {
          parsed = JSON.parse(line);
        } catch (err) {
          sendError(null, -32700, 'Parse error: each line must be one JSON-RPC 2.0 message.');
          parsed = undefined;
        }

        if (parsed !== undefined) {
          // Requests run concurrently; handleMessage writes its own reply when
          // it finishes, and JSON-RPC ids let the client match them up.
          track(handleMessage(parsed));
        }
      }

      newlineIndex = buffer.indexOf('\n');
    }
  });

  process.stdin.on('end', () => {
    stdinEnded = true;
    exitWhenDrained();
  });
  process.stdin.on('error', (err) => {
    log('stdin error:', err && err.message);
    process.exit(1);
  });

  log(`ready on stdio (default currency ${config.defaultCurrency}, radius ${config.defaultRadiusKm}km)`);
}

main();
