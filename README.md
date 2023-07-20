# 👻 genie
![Build Status](https://github.com/hupe1980/genie/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/genie.svg)](https://pkg.go.dev/github.com/hupe1980/genie)
> Genie is a Proof of Concept (POC) source code generator that showcases the potential of utilizing Large Language Models (LLMs) for code generation. As a limited prototype, Genie provides a glimpse into the capabilities of LLM-based code generation tools. It allows users to experiment with generating source code based on simplified prompts or descriptions.

Genie is based on https://github.com/smol-ai/developer.

## How to use
```text
Usage:
  genie [flags]

Examples:
genie -p "Create a python hello world"
genie -p prompt.txt

Flags:
      --api-key string        openAI api key
  -h, --help                  help for genie
  -m, --model string          model to use (default "gpt-3.5-turbo")
  -o, --outdir string         outdir to use (default "dist")
  -p, --prompt string         prompt to use (required)
  -t, --temperature float32   temperature to use (default 0.4)
  -v, --version               version for genie
```
## Example
```bash
genie -p _examples/aws_cdk/prompt.md -o _examples/aws_cdk/dist
```

Outputs:
```text

 ██████  ███████ ███    ██ ██ ███████
██       ██      ████   ██ ██ ██
██   ███ █████   ██ ██  ██ ██ █████
██    ██ ██      ██  ██ ██ ██ ██
 ██████  ███████ ██   ████ ██ ███████

Create list of files:
• src/api/api.ts
• src/api/todoController.ts
• src/lambda/todoHandler.ts
• src/dynamoDB/todoTable.ts
• .env
• .gitignore
• README.md
• cdk.json
• jest.config.js
• package.json
• projenrc.json
• tsconfig.json

Reasoning:
• The app requires an API Gateway to expose the CRUD endpoints.
• The Lambda functions will handle the API requests and interact with DynamoDB, so we need to create a Lambda Backend.
• DynamoDB is needed to store the todo items.
• We'll use 'projen' to manage the AWS CDK project and its dependencies.
• To implement the CRUD operations, we'll need to create the necessary API endpoints.

Create list of shared Dependecies:
• {dataSchemas Shared data schemas for the request and response bodies [TodoItem CreateTodoRequest CreateTodoResponse GetAllTodosResponse GetTodoByIdResponse UpdateTodoRequest UpdateTodoResponse]}
• {environmentVariables Shared environment variables for configuring DynamoDB table name and settings [DYNAMODB_TABLE_NAME]}
• {errorHandlingUtilities Shared utilities for error handling [createErrorResponse]}
• {awsCdkLibraries Shared AWS CDK related libraries [aws_cdk aws_lambda aws_apigateway aws_dynamodb]}

Reasoning:
• To identify the shared dependencies between the generated files, we need to analyze the content of each file and look for common symbols or entities used across multiple files.
• Starting with the API endpoints, there might be shared data schemas for the request and response bodies.
• The todoController.ts file is responsible for handling the API requests and might have dependencies on the data schemas.
• The todoHandler.ts file implements the Lambda functions for CRUD operations and might have dependencies on the data schemas as well.
• The todoTable.ts file creates the DynamoDB table and might have dependencies on the data schemas, as it needs to define the primary key.
• Other shared dependencies could include environment variables, error handling utilities, and AWS CDK related libraries.
• We'll analyze the content of each file and identify the shared symbols to create the shared_dependencies list.

File src/lambda/todoHandler.ts created
File src/dynamoDB/todoTable.ts created
File .env created
File src/api/api.ts created
File .gitignore created
File src/api/todoController.ts created
File cdk.json created
File jest.config.js created
File package.json created
File tsconfig.json created
File projenrc.json created
File README.md created

Tokens Used: 28404
Prompt Tokens: 24452
Completion Tokens: 3952
Successful Requests: 14
Total Cost (USD): $0.00
```

## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request for any improvements or new features you would like to see.

## References
- https://github.com/smol-ai/developer
- https://github.com/hupe1980/golc

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.