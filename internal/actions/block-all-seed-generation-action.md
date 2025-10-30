
# Block All Seed Generation Action

## Description
The Block All Seed Generation Action is designed to generate seeds for compute runs the follow a realization/block structure. It creates deterministic seeds for event, block, and realization levels, which can be used for reproducible simulations or experiments. The action supports multiple storage backends for both input blocks and output seeds.

## Implementation Details
This action implements the `cc.ActionRunner` interface and is registered in the action registry under the name `block-all-seed-generation`. It reads block data from either an event store or a JSON data source, computes seeds for each block using the provided plugins, and stores the results in either an event store or JSON format.

## Process Flow
1. Initialize action with provided attributes
2. Read block data from specified input source (eventstore or JSON)
3. Compute seeds for each block using the provided plugins
4. Store computed seeds in the specified output storage type (eventstore or JSON)

## Configuration

### Environment

- Requires access to configured data stores (eventstore or file-based)
- Must have proper plugin manager setup
- Requires appropriate permissions for reading/writing to configured data sources
---

### Attributes

### Action
| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `initial_event_seed` | int64 | Yes | Starting seed value for event-level randomization |
| `plugins` | string array | Yes | List of plugin names that will receive seeds |
| `block_dataset_name` | string | Yes | Name of the input dataset containing blocks |
| `block_store_type` | string | Yes | Storage type for input blocks ("eventstore" or "file") |
| `seed_dataset_name` | string | Yes | Name of the output dataset for seeds |
| `seed_store_type` | string | Yes | Storage type for output seeds ("eventstore" or "file") |

### Global
- No global attributes required
---

### Input Configuration
- Block data source (either eventstore or JSON file)

### Input Data Sources
- **blocks**: Blocks are read either directly from an **event** store or from a **Data Source** file resource.  The `block_dataset_name` is used as the resource identifier for both stores. 
  - Data Source:
    - when configured for a Data Source, the configuration requires a data source named with the same value as the `blocks_dataset_name` attribute.  The datasource should have a single **Path** with a key value of **default**

## Output Configuration
- **seeds**: Seeds are written to an **event** store or a **Data Source** file resource depending ont he configuration. 
  - **event store**: When event store is configured the seeds will be written to the store as a multi-dimensional dense array (refer to event store documentation for additional details) using the `seed_dataset_name` for the array name.
  - **Data Source**: When a data source is configured for the output, a Data Source configuration must be included in the configuration with the name of the data source being the same as `seed_dataset_name` and including a path with a key value of **default**   
---
## Configuration Examples

### Using Event Store for Both Input and Output
```json
{
  "action": "block-all-seed-generation",
  "attributes": {
    "initial_event_seed": 1000,
    "initial_realization_seed": 2000,
    "plugins": ["plugin1", "plugin2", "plugin3"],
    "block_dataset_name": "my-blocks",
    "block_store_type": "eventstore",
    "seed_dataset_name": "my-seeds",
    "seed_store_type": "eventstore"
  }
}
```

### Using File Store for Input and Output
```json
{
  "action": "block-all-seed-generation",
  "attributes": {
    "initial_event_seed": 1000,
    "initial_realization_seed": 2000,
    "plugins": ["plugin1", "plugin2"],
    "block_dataset_name": "blocks.json",
    "block_store_type": "file",
    "seed_dataset_name": "seeds.json",
    "seed_store_type": "eventstore"
  }
  "inputs": [
    {
      "name": "blocks.json",
      "paths": {
        "default": "conformance/simulations/blocks.json"
      },
      "store_name": "FFRD"
    }
  ],
  "outputs": [
    {
      "name": "seeds.json",
      "paths": {
        "default": "conformance/simulations/seeds.json"
      },
      "store_name": "FFRD"
    }
  ]
}
```

## Outputs

### Format
- **eventstore**: Multidimensional array with dimensions for Events, Plugins, and seed values
- **file**: JSON array containing seed information for each block and plugin combination

### Data Structures
the seed information is structured as a map of seedsets.

### Seed Set Fields
- `event_seed`: Random seed value for event-level randomization
- `block_seed`: Random seed value for block-level randomization
- `realization_seed`: Random seed value for realization-level randomization

## Error Handling
Errors are logged to the compute environment and processing will stop on error

## Usage Notes
- For large datasets, consider using eventstore for better performance and memory management