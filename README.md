# README - High-Level Scripts for AWS CLI

## Description
This project provides a set of high-level scripts to facilitate interaction with Amazon Web Services (AWS) through the command line. The scripts allow performing common cloud infrastructure management tasks without needing to use the AWS web console. Available functionalities include managing EC2 instances, RDS and DynamoDB databases, S3 storage, Auto Scaling and CloudWatch.

## Technologies Used
- **Programming Language:** Go
- **SDK:** AWS SDK for Go
- **CLI Framework:** Cobra

## List of Features
The scripts implement various functionalities for different AWS services:

### Auto Scaling
- List, create, update, and delete Auto Scaling groups.
- Retrieve instances associated with an Auto Scaling group.

### CloudWatch
- Create and list alarms, metrics, and logs.
- Retrieve log events.
- Delete alarms, metrics, and log groups.

### DynamoDB
- List, describe, create, and delete tables.
- Insert, retrieve, delete, and query items.

### EC2
- Create, list, start, stop, restart, and terminate instances.

### RDS
- List, create, delete, and start/stop database instances.
- Create and manage snapshots.

### S3
- List, create, and delete buckets.
- Manage objects within buckets (list, copy, delete, filter by extension).

## Usage Tutorial
Below is a guide on how to use the scripts to manage AWS services:

1. **Install Go:**
   - Download and install Go from [golang.org](https://golang.org/)

2. **Configure AWS credentials:**
   - Create or edit the `~/.aws/credentials` file with the following format:
     ```ini
     [default]
     aws_access_key_id = YOUR_ACCESS_KEY
     aws_secret_access_key = YOUR_SECRET_KEY
     ```
   - Alternatively, configure environment variables:
     ```sh
     export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
     export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
     ```

3. **Clone the repository:**
   ```sh
   git clone <REPOSITORY_URL>
   cd <REPOSITORY_NAME>
   ```

4. **Build the project:**
   ```sh
   go build -o aws-cli-tool
   ```

5. **Execute commands:**
   ```sh
   ./aws-cli-tool s3 list
   ./aws-cli-tool ec2 list --all
   ```

## References
- [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
