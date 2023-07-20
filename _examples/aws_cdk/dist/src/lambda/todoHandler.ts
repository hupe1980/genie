import { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import * as AWS from 'aws-sdk';

const dynamoDB = new AWS.DynamoDB.DocumentClient();
const tableName = process.env.DYNAMODB_TABLE_NAME;

export const createTodoHandler = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const requestBody = JSON.parse(event.body);
    const { title, description } = requestBody;

    // Generate unique ID for the todo item
    const id = generateUniqueId();

    // Get current timestamp
    const timestamp = new Date().toISOString();

    // Create a new todo item
    const params: AWS.DynamoDB.DocumentClient.PutItemInput = {
      TableName: tableName,
      Item: {
        id,
        title,
        description,
        createdAt: timestamp,
        updatedAt: timestamp
      }
    };

    // Save the todo item to DynamoDB
    await dynamoDB.put(params).promise();

    // Return the created todo item
    return {
      statusCode: 201,
      body: JSON.stringify({
        id,
        title,
        description,
        createdAt: timestamp,
        updatedAt: timestamp
      })
    };
  } catch (error) {
    // Handle error and return error response
    return createErrorResponse(500, 'Failed to create the todo item.');
  }
};

const generateUniqueId = (): string => {
  return Date.now().toString(36) + Math.random().toString(36).substr(2, 5);
};

const createErrorResponse = (
  statusCode: number,
  message: string
): APIGatewayProxyResult => {
  return {
    statusCode,
    body: JSON.stringify({
      error: message
    })
  };
};

// Export additional CRUD functions and exports used by other files as well