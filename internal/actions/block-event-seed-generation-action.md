# Block Event Seed Generation Action

## Description
The Block Event Seed Generation Action is responsible for generating deterministic random seeds for event processing based on block data and initial seed values. This action ensures reproducible results by using block information and provided seed values to compute new seeds for event execution.

## Implementation Details
This action implements the `cc.ActionRunner` interface and is registered under the name "block-event-seed-generation". It reads block data from an input data source, processes this data along with provided seed values, and outputs computed seeds to a designated output data source.

## Process Flow
1. Initialize the action with registered name
2. Extract required attributes (initial_event_seed, initial_realization_seed, plugins)
3. Read block data from input data source named "blockfile"
4. Parse block data into structured format
5. Compute seeds using the BlockModel.Compute method
6. Marshal computed results to JSON
7. Write results to output data source named "seeds"

## Configuration

### Environment
- Requires access to input/output data sources
- Must have proper permissions for reading block data and writing seed data

### Attributes
- `initial_event_seed` (int64): Initial seed value for event generation
- `initial_realization_seed` (int64): Initial seed value for realization generation  
- `plugins` (string slice): List of plugin names to be used in seed generation

### Action
- Action name: "block-event-seed-generation"
- Required attributes: initial_event_seed, initial_realization_seed, plugins
- Input data source: "blockfile" with path key "default"
- Output data source: "seeds" with path key "default"

### Global
- None specified

### Inputs
- Block data from "blockfile" data source
- Event identifier from plugin manager

### Input Data Sources
- `blockfile` (DataSourceName): Contains block data for seed computation
- PathKey: "default"

### Outputs
- Computed seed data written to "seeds" data source

### Output Data Sources
- `seeds` (DataSourceName): Contains computed seed data
- PathKey: "default"

## Configuration Examples

```json
{
  "action": "block-event-seed-generation",
  "attributes": {
    "initial_event_seed": 12345,
    "initial_realization_seed": 67890,
    "plugins": ["plugin1", "plugin2", "plugin3"]
  },
  "inputs": {
    "blockfile": {
      "path_key": "default"
    }
  },
  "outputs": {
    "seeds": {
      "path_key": "default"
    }
  }
}
```

## Outputs

### Format
JSON formatted output containing computed seed values

### Fields
- `event_seed` (int64): Computed event seed value
- `realization_seed` (int64): Computed realization seed value
- `plugins` (string slice): List of plugins used in computation

### Field Definitions
- `event_seed`: Random seed value for event processing
- `realization_seed`: Random seed value for realization processing
- `plugins`: Array of plugin names used in the seed computation process

## Error Handling
The action handles several error conditions:
- Invalid JSON in block data
- Missing required attributes
- Data source access errors
- Invalid event identifier format
- Computation errors from BlockModel.Compute
- JSON marshaling errors

All errors are propagated up to the caller with descriptive messages.

## Usage Notes
- This action requires deterministic block data for reproducible results
- The event identifier from plugin manager should be properly formatted as an integer
- Ensure input data source "blockfile" contains valid JSON block data
- Output data source "seeds" must be writable and properly configured

## Future Enhancements
- Support for custom seed computation algorithms
- Integration with external seed generation services
- Enhanced logging and monitoring capabilities
- Support for multiple output formats (JSON, Protobuf, etc.)
- Parallel seed computation for large block sets

## Patterns and Best Practices
- Use deterministic computation for reproducible results
- Proper resource management with deferred Close() calls
- Validate input data before processing
- Handle errors gracefully with meaningful error messages
- Follow the existing plugin architecture patterns
- Maintain backward compatibility in output format