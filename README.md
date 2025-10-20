### Seed Generator

For simulations that contain a nested loop for natural variability and knowledge uncertainty there is a need for seeds for realization and events. For simulations that also include multiple events for a synthetic year, each event is further associated with a "block" or synthetic year which has a conisistent seed. The realization seed is held constant across a realization, and the block or synthetic year seed is held constant across a block or synthetic year.

## Plugin
Seed generation is the second step in a simulation. It defines the reproducable seeds for each plugin for each realization, block and event. It is an action in the Seed-Generator plugin.

## Action Parameterization
Seed generation relies on an action in the seed-generator plugin named "block_all_seed_generation" all input parameters are defined in the action's attributes and input datasources. The action produces a single ouput specified at the action level named seeds.json. For FFRD the seeds are generated based on the structure of the blocks.json file to determine the count and structure of seeds, and the names of the plugin/model unit names to define how many events are generated for each event.
# Attributes
 - "block_dataset_name": the name of the dataset in the input datasources that should be used to define the number of events per block, blocks per realization, and total realization count.
 - "block_store_type": defines the store type for the blocks. For FFRD the SOP dictates json format.
 - "seed_dataset_name": the name of the datasource in the output datasources of the action that represents the output for seeds.
 - "seed_dataset_name": defines the store type for the seeds. For FFRD the SOP dictates json format.
 - "initial_event_seed": an int64 that defines the seed for a random number generator that generates event seeds.
 - "initial_block_seed": an int64 that defines the seed for a random number generator that generates block seeds.
 - "initial_realization_seed": an int64 that defines the seed for a random number generator that generates realization seeds.
 - "plugins": this is an array of string, each string representing a plugin that needs the event, block, or realization seeds for each event.

 ## Steps to compute on aws
 1. Ensure that the seed-generator plugin is loaded into your instance of ECR and the image and tag name in the seed-generator-plugin-manifest.json is correct.
 2. Set the scenario attribute in the compute-manifest.json payload attributes to be consistent with the conformance or production phase
    - Conformance: "scenario": "conformance"
    - Production: "scenario": "production"
 3. Ensure the store is set correctly
    - "store_type": "S3"
    - "params": {"root": "model-library/ffrd-trinity} (where the name of the HUC4 basin is correctly identified, and if the root is not model library it should match your environment and infrastructure)
 4. Ensure your compute-stream.json file points at your compute-manifest.json and the seed-generator-plugin-manifest.json correctly, also ensure that the stream_list.csv has only one event defined (this will trigger one job in aws, which will be for the entire simulation in the case of this plugin)
 5. run the manifestor: .\manifestor --envFile="your env file" run --computeFile="compute-stream.json"

  ## Steps to compute on docker-local
 1. Ensure that the seed-generator plugin image has been downloaded into your local docker repo the image and tag name in the seed-generator-plugin-manifest.json is correct.
 2. start docker-compose.yaml to initalize minio
 3. Set the scenario attribute in the compute-manifest.json payload attributes to be consistent with the conformance or production phase
    - Conformance: "scenario": "conformance"
    - Production: "scenario": "production"
 4. Ensure the store is set correctly
    - "store_type": "S3"
    - "params": {"root": "model-library/ffrd-trinity} (where the name of the HUC4 basin is correctly identified, and if the root is not model library it should match your environment and infrastructure)
 5. Ensure your compute-stream.json file points at your compute-manifest.json and the seed-generator-plugin-manifest.json correctly, also ensure that the stream_list.csv has only one event defined (this will trigger one job in aws, which will be for the entire simulation in the case of this plugin)
 6. run the manifestor: .\manifestor --envFile="your env file" run --computeFile="compute-stream.json"