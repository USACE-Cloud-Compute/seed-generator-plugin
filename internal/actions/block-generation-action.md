# Block Generation Action

## Description

The Block Generation Action is responsible for creating synthetic event data blocks. It generates event data that can be used for testing, simulation, or demonstration purposes. The action supports both fixed-length and variable-length block generation modes, allowing flexibility in how event data is structured and organized.

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

## Configuration

### Environment

No specific environment variables required.

---

### Attributes

### Action

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `target_total_events` | int64 | Yes | Total number of events to generate |
| `blocks_per_realization` | int32 | Yes | Number of blocks per realization |
| `target_events_per_block` | int32 | Yes | Target number of events per block |
| `seed` | int64 | No | Random seed for reproducible results (default: 1234) |
| `outputdataset_name` | string | Yes | Name of the output dataset |
| `store_type` | string | Yes | Storage type ("eventstore" or other) |


### Global

No global attributes are required.

---
## Input Configuration

No explicit inputs required beyond the action attributes.

### Input Data Sources

No input data sources required.

## Output Configuration
- **blocks**: Seeds are written to an **event** store or a **Data Source** file resource depending on the configuration. 
  - **event store**: When event store is configured the seeds will be written to the store as a multi-dimensional dense array (refer to event store documentation for additional details) using the `outputdataset_name` for the array name.
  - **Data Source**: When a data source is configured for the output, a Data Source configuration must be included in the configuration with the name of the data source being the same as `outputdataset_name` and including a path with a key value of **default**   

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
  "outputs": [
    {
      "name": "blocks.json",
      "paths": {
        "default": "conformance/simulations/blocks.json"
      },
      "store_name": "FFRD"
    }
  ]
}
```

## Outputs

### Format

JSON array of block objects.

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `realization_index` | int | Realization index which is 1-indexed |
| `block_index` | int | Block index which is 1-indexed |
| `block_event_count` | int | Number of events in the block |
| `block_event_start` | int | Starting event number for the block |
| `block_event_end` | int | Ending event number for the block |


## Error Handling
Errors are logged to the compute environment and processing will stop on error

## Usage Notes
- For large datasets, consider using eventstore for better performance and memory management
- Event Stores are considered experimental and might change in the future