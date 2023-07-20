import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';

export class TodoTable extends cdk.Construct {
  constructor(scope: cdk.Construct, id: string, tableName: string) {
    super(scope, id);

    const table = new dynamodb.Table(this, 'TodoTable', {
      tableName: tableName,
      partitionKey: {
        name: 'id',
        type: dynamodb.AttributeType.STRING,
      },
    });

    table.addGlobalSecondaryIndex({
      indexName: 'CreatedAtIndex',
      partitionKey: {
        name: 'createdAt',
        type: dynamodb.AttributeType.STRING,
      },
      projectionType: dynamodb.ProjectionType.ALL,
    });

    table.addGlobalSecondaryIndex({
      indexName: 'UpdatedAtIndex',
      partitionKey: {
        name: 'updatedAt',
        type: dynamodb.AttributeType.STRING,
      },
      projectionType: dynamodb.ProjectionType.ALL,
    });
  }
}