# Seed Generator

For simulations that contain a nested loop for natural variability and knowledge uncertainty there is a need for seeds for realization and events. For simulations that also include multiple events for a synthetic year, each event is further associated with a "block" or synthetic year which has a conisistent seed. The realization seed is held constant across a realization, and the block or synthetic year seed is held constant across a block or synthetic year.

The seed generator plugin defines the reproducable seeds for stochasic simulation realizations and blocks 

## Seed Generator Image
The current image is in the usace cloud compute ghcr at: `ghcr.io/usace-cloud-compute/seed-generator-plugin:latest`.

## Resource Requirements
The seed generate requires little in the way of resources:
  - vCPU: 1 or 2
  - Memory: >1GB

## Actions and Action documentation
  - [Block All Seed Generation](./internal/actions/block-all-seed-generation-action.md): generate seeds for compute runs the follow a realization/block structure
  - [Block Event Seed Generation](./internal/actions/block-event-seed-generation-action.md): generate seeds for an events structure
  - [Block Generation](./internal/actions/block-generation-action.md): generate blocks of events
  - [Realization Seed Generation](./internal/actions/realization-seed-generation-action.md): generate seeds for compute realizations

## Running
To run the seed generator, the you will need to perform the following:
  - Create a `plugin-manifest` this is used to `register` the seed-generator with your compute provider.  Note that if you have already registered the seed-generator in your environment, you will not need to do so again.
    - sample seed generator plugin manifests can be referenced from the following:
      - AWS Batch...
      - Local Docker...    

## Plugin Manifests
Detailed descriptions for the plugin manifest configuration are provided below.  The following fields are required:
  - `name`: the required name represents the name that this plugin will be registered under.  In AWS batch is is effectevely the `job-description` name.
  - `image_and_tag`: the required container repository referecne for the image you would like to run.  Currently that is: `ghcr.io/usace-cloud-compute/seed-generator-plugin:latest`
  - `description`: an optional description for the registered plugin
  - `compute_environment`: the compute resource requirements
    - `vcpu`: 1 vcpu is adequate for this plugin
    - `memory`: 1GB or less is ok for this plugin
    - `extraHosts`: if you are running using the Docker provider on a local system, you might need to set an extra host.  This is not necessary on windows or osx, but might be on a local linux system
  - `environment`: environment variables necessary to run the image.  Fro the seed generator, you will need two sets of environment variables:
    - `cloud compute storage`: environment variables for the cloud compute sdk to retrieve payloads from cloud compute storage
    - `data storage`: environment variabled necessary to connect, read, and write to the data store
  - `credentials`:
    - sets of credentials for the clouyd compute storagte and data storege.  Credentials are managed in the compute provider credential management system, for example AWS Secrets manager for Amazon.

  - #### Seed Generator Local Docker Plugin Manifest using local minio to emulate S3 storage
    ```json
    {
        "name":"FFRD-SEED-GENERATOR",
        "image_and_tag":"ghcr.io/usace-cloud-compute/seed-generator-plugin:latest",
        "description":"Seed Generator",
        "compute_environment":{
            "vcpu":"1",
            "memory":"1000"
        },
        "environment":[
            {
                "name":  "CC_AWS_DEFAULT_REGION",
                "value": "us-east-1"
            },
            {
                "name":  "CC_AWS_S3_BUCKET",
                "value": "ccstore"
            },
            {
                "name":  "CC_AWS_ENDPOINT",
                "value": "http://host.docker.internal:9000"
            },
            {
                "name":  "FFRD_AWS_DEFAULT_REGION",
                "value": "us-east-1"
            },
            {
                "name":  "FFRD_AWS_S3_BUCKET",
                "value": "project-data"
            },
            {
                "name":  "FFRD_AWS_ENDPOINT",
                "value": "http://host.docker.internal:9000"
            }
        ],
        "credentials":[
            {
                "name":  "CC_AWS_ACCESS_KEY_ID",
                "value": "secretsmanager:AWS_ACCESS_KEY_ID::"
            },
            {
                "name":  "CC_AWS_SECRET_ACCESS_KEY",
                "value": "secretsmanager:AWS_SECRET_ACCESS_KEY::"
            },
            {
                "name":  "FFRD_AWS_ACCESS_KEY_ID",
                "value": "secretsmanager:AWS_ACCESS_KEY_ID::"
            },
            {
                "name":  "FFRD_AWS_SECRET_ACCESS_KEY",
                "value": "secretsmanager:AWS_SECRET_ACCESS_KEY::"
            }
        ]
    }
    ```

    - #### Seed Generator AWS Batch Plugin Manifest
    ```json
    {
        "name":"FFRD-SEED-GENERATOR",
        "image_and_tag":"ghcr.io/usace-cloud-compute/seed-generator-plugin:latest",
        "description":"Seed Generator",
        "compute_environment":{
            "vcpu":"1",
            "memory":"1000",
            "extraHosts": ["host.docker.internal:host-gateway"]
        },
        "environment":[
            {
                "name":  "CC_AWS_DEFAULT_REGION",
                "value": "us-east-1"
            },
            {
                "name":  "CC_AWS_S3_BUCKET",
                "value": "ccstore"
            },
            {
                "name":  "FFRD_AWS_DEFAULT_REGION",
                "value": "us-east-1"
            },
            {
                "name":  "FFRD_AWS_S3_BUCKET",
                "value": "project-data"
            }
        ],
        "credentials":[
            {
                "name":  "CC_AWS_ACCESS_KEY_ID",
                "value": "arn:aws-us:secretsmanager:us-east-1:012345687:secret:prod/CloudCompute/S3-gggg:AWS_ACCESS_KEY_ID::"
            },
            {
                "name":  "CC_AWS_SECRET_ACCESS_KEY",
                "value": "arn:aws-us:secretsmanager:us-east-1:012345687:secret:prod/CloudCompute/S3-gggg:AWS_ACCESS_KEY_ID::"
            },
            {
                "name":  "FFRD_AWS_ACCESS_KEY_ID",
                "value": "arn:aws-us:secretsmanager:us-east-1:012345687:secret:prod/CloudCompute/S3-gggg:AWS_ACCESS_KEY_ID::"
            },
            {
                "name":  "FFRD_AWS_SECRET_ACCESS_KEY",
                "value": "arn:aws-us:secretsmanager:us-east-1:012345687:secret:prod/CloudCompute/S3-gggg:AWS_ACCESS_KEY_ID::"
            }
        ]
    }
    ```