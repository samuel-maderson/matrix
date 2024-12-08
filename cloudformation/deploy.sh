#!/bin/bash

aws cloudformation create-stack --stack-name my-dynamodb-stack --template-body file://dynamodb.yaml