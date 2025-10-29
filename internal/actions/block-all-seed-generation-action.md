
# Block All Seed Generation Action

## Description
The Block All Seed Generation Action is designed to generate seeds for all blocks within a dataset. It creates deterministic seeds for event, block, and realization levels, which can be used for reproducible simulations or experiments. The action supports multiple storage backends for both input blocks and output seeds.

## Implementation Details
This action implements the `cc.ActionRunner` interface and is registered in the action registry under the name "block-all-seed-generation". It reads block data from either an event store or a JSON data source, computes seeds for each block using the provided plugins, and stores the results in either an event store or JSON format.

## Process Flow
1. Initialize action with provided attributes
2. Read block data from specified input source (eventstore or JSON)
3. Compute seeds for each block using the provided plugins
4. Store computed seeds in the specified output storage type (eventstore or JSON)
5. Handle metadata storage for plugin information when using eventstore

## Configuration

### Environment
- Requires access to configured data stores (eventstore or file-based)
- Must have proper plugin manager setup
- Requires appropriate permissions for reading/writing to configured data sources

### Attributes
- `initial_event_seed` (int64): Starting seed value for event-level randomization
- `initial_realization_seed` (int64): Starting seed value for realization-level randomization
- `plugins` (string slice): List of plugin names that will receive seeds
- `block_dataset_name` (string): Name of the input dataset containing blocks
- `block_store_type` (string): Storage type for input blocks ("eventstore" or "file")
- `seed_dataset_name` (string): Name of the output dataset for seeds
- `seed_store_type` (string): Storage type for output seeds ("eventstore" or "file")

### Action
- Action name: "block-all-seed-generation"
- Required attributes: initial_event_seed, initial_realization_seed, plugins, block_dataset_name, block_store_type, seed_dataset_name, seed_store_type

### Global
- No global configuration required

### Inputs
- Block data source (either eventstore or JSON file)
- Plugin configuration
- Seed initialization values

### Input Data Sources
- **eventstore**: Read blocks from an event store using the provided block_dataset_name
- **file**: Read blocks from a JSON file using the provided block_dataset_name

### Outputs
- Seed data for all blocks and plugins
- Metadata about plugin assignments (when using eventstore)

### Output Data Sources
- **eventstore**: Store seeds in a multidimensional array format with dimensions for Events, Plugins, and seed values
- **file**: Store seeds in JSON format

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

### Using File Store for Input and Event Store for Output
```json
{
  "action": "block-all-seed-generation",
  "attributes": {
    "initial_event_seed": 1000,
    "initial_realization_seed": 2000,
    "plugins": ["plugin1", "plugin2"],
    "block_dataset_name": "blocks.json",
    "block_store_type": "file",
    "seed_dataset_name": "seeds",
    "seed_store_type": "eventstore"
  }
}
```

## Outputs

### Format
- **eventstore**: Multidimensional array with dimensions for Events, Plugins, and seed values
- **file**: JSON array containing seed information for each block and plugin combination

### Fields
- `event_seed`: Random seed value for event-level randomization
- `block_seed`: Random seed value for block-level randomization
- `realization_seed`: Random seed value for realization-level randomization

### Field Definitions
- `event_seed`: Used for generating random numbers at the event level
- `block_seed`: Used for generating random numbers at the block level
- `realization_seed`: Used for generating random numbers at the realization level

## Error Handling
The action handles several error conditions:
- Invalid input attributes
- Failed reading from input data sources
- Failed writing to output data sources
- Incompatible store types (non-multidimensional array stores when using eventstore)
- JSON marshaling/unmarshaling errors
- File I/O errors

## Usage Notes
- The action requires that all plugins be specified in the plugins attribute
- When using eventstore for output, the action creates a multidimensional array with proper dimensions
- The action supports both deterministic and non-deterministic seed generation based on input parameters
- For large datasets, consider using eventstore for better performance and memory management

## Future Enhancements
- Support for more storage backends
- Additional seed generation algorithms
- Parallel processing capabilities for large datasets
- Integration with external seed generation services
- Enhanced metadata storage and retrieval capabilities

## Patterns and Best Practices
- Always specify all required plugins in the plugins attribute
- Use appropriate seed values to ensure reproducible results
- Consider using eventstore for large datasets to optimize memory usage
- Validate input data sources before processing
- Handle errors gracefully and provide meaningful error messages
- Use consistent naming conventions for data sources and datasets