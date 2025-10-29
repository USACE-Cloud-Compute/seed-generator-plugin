# Block Generation Action

## Description

The Block Generation Action is responsible for creating synthetic event data blocks based on specified parameters. It generates event data that can be used for testing, simulation, or demonstration purposes. The action supports both fixed-length and variable-length block generation modes, allowing flexibility in how event data is structured and organized.

## Implementation Details

The action implements the `cc.ActionRunner` interface and uses a `BlockGenerator` to create event data blocks. It supports two modes:
1. Standard block generation with variable block lengths
2. Fixed-length block generation with predetermined block sizes

The generated blocks are either stored in an event store or written to a data source, depending on the configuration.

## Process Flow

1. Initialize action with configuration parameters
2. Determine generation mode (fixed-length or standard)
3. Create `BlockGenerator` with specified parameters
4. Generate blocks using appropriate method
5. Marshal blocks to JSON format
6. Store blocks in configured output destination (eventstore or data source)
7. Return success or error status

## Configuration

### Environment

No specific environment variables required.

### Attributes

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `target_total_events` | int64 | Yes | Total number of events to generate |
| `blocks_per_realization` | int32 | Yes | Number of blocks per realization |
| `target_events_per_block` | int32 | Yes | Target number of events per block |
| `seed` | int64 | No | Random seed for reproducible results (default: 1234) |
| `outputdataset_name` | string | Yes | Name of the output dataset |
| `store_type` | string | Yes | Storage type ("eventstore" or other) |

### Action

```json
{
  "name": "block-generation",
  "attributes": {
    "target_total_events": 1000,
    "blocks_per_realization": 10,
    "target_events_per_block": 100,
    "seed": 1234,
    "outputdataset_name": "generated-events",
    "store_type": "eventstore"
  }
}
```

### Global

No global configuration required.

### Inputs

No explicit inputs required beyond the action attributes.

### Input Data Sources

No input data sources required.

### Outputs

The action generates event data blocks that are stored in the configured output destination.

### Output Data Sources

- Event Store: When `store_type` is "eventstore"
- Data Source: When `store_type` is any other value

## Configuration Examples

### Standard Block Generation

```json
{
  "name": "block-generation",
  "attributes": {
    "target_total_events": 5000,
    "blocks_per_realization": 5,
    "target_events_per_block": 1000,
    "seed": 42,
    "outputdataset_name": "standard-blocks",
    "store_type": "eventstore"
  }
}
```

### Fixed-Length Block Generation

```json
{
  "name": "block-generation-fixed-length",
  "attributes": {
    "target_total_events": 10000,
    "blocks_per_realization": 20,
    "target_events_per_block": 500,
    "seed": 999,
    "outputdataset_name": "fixed-length-blocks",
    "store_type": "filestore"
  }
}
```

## Outputs

### Format

JSON array of block objects containing event data.

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `BlockID` | string | Unique identifier for the block |
| `Events` | array | Array of event objects |
| `EventCount` | int32 | Number of events in the block |
| `Timestamp` | string | Block creation timestamp |

### Field Definitions

- `BlockID`: A unique identifier for each generated block
- `Events`: Array containing the actual event data objects
- `EventCount`: Number of events contained in the block
- `Timestamp`: Creation timestamp of the block in ISO format

## Error Handling

The action handles several error conditions:
1. JSON marshaling failures when encoding blocks
2. Event store creation failures
3. Data source write failures
4. Missing required attributes
5. Invalid configuration parameters

All errors are logged with descriptive messages and propagated up the call stack.

## Usage Notes

1. The action requires proper configuration of output destinations
2. For reproducible results, always specify a seed value
3. The `target_total_events` parameter should be divisible by `blocks_per_realization` for optimal results
4. When using eventstore, ensure the event store is properly configured and accessible
5. For large datasets, consider the memory implications of JSON marshaling

## Future Enhancements

1. Add support for different event data formats (CSV, XML)
2. Implement parallel block generation for better performance
3. Add data validation for configuration parameters
4. Support for custom event generators
5. Integration with more storage backends

## Patterns and Best Practices

1. Always provide default values for optional parameters
2. Use descriptive error messages for debugging
3. Implement proper logging for monitoring and debugging
4. Validate input parameters before processing
5. Follow the existing code structure and naming conventions
6. Ensure thread safety when dealing with shared resources
7. Provide comprehensive documentation for all configuration options