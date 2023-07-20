import * as aws_cdk from 'aws-cdk-lib';
import * as aws_apigateway from 'aws-cdk-lib/aws-apigateway';
import * as aws_lambda from 'aws-cdk-lib/aws-lambda';
import { DYNAMODB_TABLE_NAME } from '../environmentVariables';

export class TodoApi {
  constructor(scope: aws_cdk.Construct) {
    // Create Lambda function for CreateTodo
    const createTodoLambda = new aws_lambda.Function(scope, 'CreateTodoLambda', {
      runtime: aws_lambda.Runtime.NODEJS_14_X,
      handler: 'index.handler',
      code: aws_lambda.Code.fromAsset('lambda'),
      environment: {
        DYNAMODB_TABLE_NAME: DYNAMODB_TABLE_NAME,
      },
    });

    // Create Lambda function for GetAllTodos
    const getAllTodosLambda = new aws_lambda.Function(scope, 'GetAllTodosLambda', {
      runtime: aws_lambda.Runtime.NODEJS_14_X,
      handler: 'index.handler',
      code: aws_lambda.Code.fromAsset('lambda'),
      environment: {
        DYNAMODB_TABLE_NAME: DYNAMODB_TABLE_NAME,
      },
    });

    // Create Lambda function for GetTodoById
    const getTodoByIdLambda = new aws_lambda.Function(scope, 'GetTodoByIdLambda', {
      runtime: aws_lambda.Runtime.NODEJS_14_X,
      handler: 'index.handler',
      code: aws_lambda.Code.fromAsset('lambda'),
      environment: {
        DYNAMODB_TABLE_NAME: DYNAMODB_TABLE_NAME,
      },
    });

    // Create Lambda function for UpdateTodo
    const updateTodoLambda = new aws_lambda.Function(scope, 'UpdateTodoLambda', {
      runtime: aws_lambda.Runtime.NODEJS_14_X,
      handler: 'index.handler',
      code: aws_lambda.Code.fromAsset('lambda'),
      environment: {
        DYNAMODB_TABLE_NAME: DYNAMODB_TABLE_NAME,
      },
    });

    // Create Lambda function for DeleteTodo
    const deleteTodoLambda = new aws_lambda.Function(scope, 'DeleteTodoLambda', {
      runtime: aws_lambda.Runtime.NODEJS_14_X,
      handler: 'index.handler',
      code: aws_lambda.Code.fromAsset('lambda'),
      environment: {
        DYNAMODB_TABLE_NAME: DYNAMODB_TABLE_NAME,
      },
    });

    // Create API Gateway REST API
    const api = new aws_apigateway.RestApi(scope, 'TodoApi', {
      restApiName: 'Todo API',
    });

    // Create API Gateway resources
    const todosResource = api.root.addResource('todos');
    const todoResource = todosResource.addResource('{id}');

    // Create API Gateway methods
    todosResource.addMethod('POST', new aws_apigateway.LambdaIntegration(createTodoLambda));
    todosResource.addMethod('GET', new aws_apigateway.LambdaIntegration(getAllTodosLambda));
    todoResource.addMethod('GET', new aws_apigateway.LambdaIntegration(getTodoByIdLambda));
    todoResource.addMethod('PUT', new aws_apigateway.LambdaIntegration(updateTodoLambda));
    todoResource.addMethod('DELETE', new aws_apigateway.LambdaIntegration(deleteTodoLambda));
  }
}