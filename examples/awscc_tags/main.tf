terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    awscc = {
      source = "hashicorp/awscc"
    }
  }
}

provider "aws" {
  region = "us-west-2"

  # AWS provider supports default_tags
  default_tags {
    tags = {
      Owner   = "DevOps"
      Project = "Terratags"
    }
  }
}

provider "awscc" {
  region = "us-west-2"
  # Note: AWSCC provider doesn't support default_tags
}

# AWSCC resource with tags in list-of-maps format
resource "awscc_apigateway_rest_api" "example" {
  name = "example-api"

  # AWSCC provider uses list of maps with key/value pairs for tags
  tags = [
    {
      key   = "Name"
      value = "Example API"
    },
    {
      key   = "Environment"
      value = "Test"
    },
    {
      key   = "Owner"
      value = "API Team"
    },
    {
      key   = "Project"
      value = "Terratags"
    }
  ]
}


# AWSCC resource with non compliant tag format

resource "awscc_msk_cluster" "example" {
  cluster_name           = "example-msk-cluster"
  kafka_version          = "2.8.1"
  number_of_broker_nodes = 2

  broker_node_group_info = {
    instance_type = "kafka.t3.small"
    client_subnets = [
      awscc_ec2_subnet.msk_subnet_1.id,
      awscc_ec2_subnet.msk_subnet_2.id
    ]
    security_groups = [awscc_ec2_security_group.msk_sg.id]
    storage_info = {
      ebs_storage_info = {
        volume_size = 100
      }
    }
  }

  encryption_info = {
    encryption_in_transit = {
      client_broker = "TLS"
      in_cluster    = true
    }
  }

  enhanced_monitoring = "DEFAULT"

  logging_info = {
    broker_logs = {
      cloudwatch_logs = {
        enabled = true
      }
    }
  }

  tags = {
    Name        = "testcluster"
    Environment = "Production"
    Owner       = "Network"
    Project     = "Infrastructure"
    CostCenter  = "CC123"
  }
}


# AWS resource with tags in map format
resource "aws_s3_bucket" "example" {
  bucket = "example-bucket"

  # AWS provider uses map format for tags
  tags = {
    Name        = "Example Bucket"
    Environment = "Test"
    # Owner and Project come from default_tags
  }
}

# AWSCC resource with tags in list-of-maps format
resource "awscc_apigateway_rest_api" "example" {
  name = "example-api"

  # AWSCC provider uses list of maps with key/value pairs for tags
  tags = [
    {
      key   = "Name"
      value = "Example API"
    },
    {
      key   = "Environment"
      value = "Test"
    },
    {
      key   = "Owner"
      value = "API Team"
    },
    {
      key   = "Project"
      value = "Terratags"
    }
  ]
}


# AWSCC resource with non compliant tag format

resource "awscc_msk_cluster" "example" {
  cluster_name           = "example-msk-cluster"
  kafka_version          = "2.8.1"
  number_of_broker_nodes = 2

  broker_node_group_info = {
    instance_type = "kafka.t3.small"
    client_subnets = [
      awscc_ec2_subnet.msk_subnet_1.id,
      awscc_ec2_subnet.msk_subnet_2.id
    ]
    security_groups = [awscc_ec2_security_group.msk_sg.id]
    storage_info = {
      ebs_storage_info = {
        volume_size = 100
      }
    }
  }

  encryption_info = {
    encryption_in_transit = {
      client_broker = "TLS"
      in_cluster    = true
    }
  }



  enhanced_monitoring = "DEFAULT"

  logging_info = {
    broker_logs = {
      cloudwatch_logs = {
        enabled = true
      }
    }
  }

  tags = {
    Name        = "testcluster"
    Environment = "Production"
    Owner       = "Network"
    Project     = "Infrastructure"
    CostCenter  = "CC123"
  }
}


resource "awscc_bedrock_agent" "example" {
  agent_name              = "example-agent"
  description             = "Example agent configuration"
  agent_resource_role_arn = var.agent_role_arn
  foundation_model        = "anthropic.claude-v2:1"
  instruction             = "You are an office assistant in an insurance agency. You are friendly and polite. You help with managing insurance claims and coordinating pending paperwork."
  knowledge_bases = [{
    description          = "example knowledge base"
    knowledge_base_id    = var.knowledge_base_id
    knowledge_base_state = "ENABLED"
  }]

  customer_encryption_key_arn = var.kms_key_arn
  idle_session_ttl_in_seconds = 600
  auto_prepare                = true

  action_groups = [{
    action_group_name = "example-action-group"
    description       = "Example action group"
    api_schema = {
      s3 = {
        s3_bucket_name = var.bucket_name
        s3_object_key  = var.bucket_object_key
      }
    }
    action_group_executor = {
      lambda = var.lambda_arn
    }

  }]

  tags = {
    "Modified By" = "AWSCC"
  }

}
