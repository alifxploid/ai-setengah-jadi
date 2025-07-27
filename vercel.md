base url : https://ai-gateway.vercel.sh/v1

Supported endpoints :

The AI Gateway currently supports the following OpenAI-compatible endpoints:

-GET /v1/models - List available models
-GET /v1/models/{model} - Retrieve a specific model
-POST /v1/chat/completions - Create chat completions with support for streaming, attachments, and tool calls

Streaming chat completion

Endpoint : POST /v1/chat/completions

Example request :

import OpenAI from 'openai';
 
const apiKey = process.env.AI_GATEWAY_API_KEY || process.env.VERCEL_OIDC_TOKEN;
 
const openai = new OpenAI({
  apiKey,
  baseURL: 'https://ai-gateway.vercel.sh/v1',
});
 
const stream = await openai.chat.completions.create({
  model: 'anthropic/claude-sonnet-4',
  messages: [
    {
      role: 'user',
      content: 'Write a one-sentence bedtime story about a unicorn.',
    },
  ],
  stream: true,
});
 
for await (const chunk of stream) {
  const content = chunk.choices[0]?.delta?.content;
  if (content) {
    process.stdout.write(content);
  }
}


response format :

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"anthropic/claude-sonnet-4","choices":[{"index":0,"delta":{"content":"Once"},"finish_reason":null}]}
 
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"anthropic/claude-sonnet-4","choices":[{"index":0,"delta":{"content":" upon"},"finish_reason":null}]}
 
data: [DONE]

Image attachments :

Example format : 

import fs from 'node:fs';
import OpenAI from 'openai';
 
const apiKey = process.env.AI_GATEWAY_API_KEY || process.env.VERCEL_OIDC_TOKEN;
 
const openai = new OpenAI({
  apiKey,
  baseURL: 'https://ai-gateway.vercel.sh/v1',
});
 
// Read the image file as base64
const imageBuffer = fs.readFileSync('./path/to/image.png');
const imageBase64 = imageBuffer.toString('base64');
 
const completion = await openai.chat.completions.create({
  model: 'anthropic/claude-sonnet-4',
  messages: [
    {
      role: 'user',
      content: [
        { type: 'text', text: 'Describe this image in detail.' },
        {
          type: 'image_url',
          image_url: {
            url: `data:image/png;base64,${imageBase64}`,
            detail: 'auto',
          },
        },
      ],
    },
  ],
  stream: false,
});
 
console.log('Assistant:', completion.choices[0].message.content);
console.log('Tokens used:', completion.usage);

PDF attachments :

example format :

import fs from 'node:fs';
import OpenAI from 'openai';
 
const apiKey = process.env.AI_GATEWAY_API_KEY || process.env.VERCEL_OIDC_TOKEN;
 
const openai = new OpenAI({
  apiKey,
  baseURL: 'https://ai-gateway.vercel.sh/v1',
});
 
// Read the PDF file as base64
const pdfBuffer = fs.readFileSync('./path/to/document.pdf');
const pdfBase64 = pdfBuffer.toString('base64');
 
const completion = await openai.chat.completions.create({
  model: 'anthropic/claude-sonnet-4',
  messages: [
    {
      role: 'user',
      content: [
        {
          type: 'text',
          text: 'What is the main topic of this document? Please summarize the key points.',
        },
        {
          type: 'file',
          file: {
            data: pdfBase64,
            media_type: 'application/pdf',
            filename: 'document.pdf',
          },
        },
      ],
    },
  ],
  stream: false,
});
 
console.log('Assistant:', completion.choices[0].message.content);
console.log('Tokens used:', completion.usage);

Tool calls :

Basic tool calls :

import OpenAI from 'openai';
 
const apiKey = process.env.AI_GATEWAY_API_KEY || process.env.VERCEL_OIDC_TOKEN;
 
const openai = new OpenAI({
  apiKey,
  baseURL: 'https://ai-gateway.vercel.sh/v1',
});
 
const tools: OpenAI.Chat.Completions.ChatCompletionTool[] = [
  {
    type: 'function',
    function: {
      name: 'get_weather',
      description: 'Get the current weather in a given location',
      parameters: {
        type: 'object',
        properties: {
          location: {
            type: 'string',
            description: 'The city and state, e.g. San Francisco, CA',
          },
          unit: {
            type: 'string',
            enum: ['celsius', 'fahrenheit'],
            description: 'The unit for temperature',
          },
        },
        required: ['location'],
      },
    },
  },
];
 
const completion = await openai.chat.completions.create({
  model: 'anthropic/claude-sonnet-4',
  messages: [
    {
      role: 'user',
      content: 'What is the weather like in San Francisco?',
    },
  ],
  tools: tools,
  tool_choice: 'auto',
  stream: false,
});
 
console.log('Assistant:', completion.choices[0].message.content);
console.log('Tool calls:', completion.choices[0].message.tool_calls);

Tool call response format :

{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "anthropic/claude-sonnet-4",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": null,
        "tool_calls": [
          {
            "id": "call_123",
            "type": "function",
            "function": {
              "name": "get_weather",
              "arguments": "{\"location\": \"San Francisco, CA\", \"unit\": \"celsius\"}"
            }
          }
        ]
      },
      "finish_reason": "tool_calls"
    }
  ],
  "usage": {
    "prompt_tokens": 82,
    "completion_tokens": 18,
    "total_tokens": 100
  }
}

Required parameters :

- model (string): The model to use for the completion (e.g., anthropic/claude-sonnet-4)
- messages (array): Array of message objects with role and content fields

Optional parameters :

- stream (boolean): Whether to stream the response. Defaults to false
- temperature (number): Controls randomness in the output. Range: 0-2
- max_tokens (integer): Maximum number of tokens to generate
- top_p (number): Nucleus sampling parameter. Range: 0-1
- frequency_penalty (number): Penalty for frequent tokens. Range: -2 to 2
- presence_penalty (number): Penalty for present tokens. Range: -2 to 2
- stop (string or array): Stop sequences for the generation
- tools (array): Array of tool definitions for function calling
- tool_choice (string or object): Controls which tools are called (auto, none, or specific function)

Message format :

Text messages :

{
  "role": "user",
  "content": "Hello, how are you?"
}

multimodal messages :

{
  "role": "user",
  "content": [
    { "type": "text", "text": "What's in this image?" },
    {
      "type": "image_url",
      "image_url": {
        "url": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD..."
      }
    }
  ]
}

File messages :

{
  "role": "user",
  "content": [
    { "type": "text", "text": "Summarize this document" },
    {
      "type": "file",
      "file": {
        "data": "JVBERi0xLjQKJcfsj6IKNSAwIG9iago8PAovVHlwZSAvUGFnZQo...",
        "media_type": "application/pdf",
        "filename": "document.pdf"
      }
    }
  ]
}

Error handling :
Common error codes
400 Bad Request: Invalid request parameters
401 Unauthorized: Invalid or missing authentication
403 Forbidden: Insufficient permissions
404 Not Found: Model or endpoint not found
429 Too Many Requests: Rate limit exceeded
500 Internal Server Error: Server error

error response format :

{
  "error": {
    "message": "Invalid request: missing required parameter 'model'",
    "type": "invalid_request_error",
    "param": "model",
    "code": "missing_parameter"
  }
}


Image analysis :

import fs from 'node:fs';
 
// Read the image file as base64
const imageBuffer = fs.readFileSync('./path/to/image.png');
const imageBase64 = imageBuffer.toString('base64');
 
const response = await fetch(
  'https://ai-gateway.vercel.sh/v1/chat/completions',
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${process.env.AI_GATEWAY_API_KEY}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: 'anthropic/claude-sonnet-4',
      messages: [
        {
          role: 'user',
          content: [
            { type: 'text', text: 'Describe this image in detail.' },
            {
              type: 'image_url',
              image_url: {
                url: `data:image/png;base64,${imageBase64}`,
                detail: 'auto',
              },
            },
          ],
        },
      ],
      stream: false,
    }),
  },
);
 
const result = await response.json();
console.log(result);


NOTED : TOOLS CALLS BISA DIGUNAKAN, DAN TOLONG KAMU PIKIRKAN DENGAN BIJAK TOOLS CALLS BISA UNTUK APA AJA, JIKA BAGUS LANGSUNG IMPLEMENTASIKAN. DAN INI DOCS LENGKAP DARI VERCEL AI GATEWAY. JADI APAKAH INI AKAN DI MASUKKAN SEMUANYA ATAU KAMU MEMILIH PERLU APA SAJA AGAR SEPERTI AI PADA UMUMNYA DI LUAR SANA.