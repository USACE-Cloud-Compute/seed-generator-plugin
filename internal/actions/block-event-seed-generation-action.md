# Block Event Seed Generation Action

## Description
The Block Event Seed Generation Action is responsible for generating deterministic random seeds for event processing based on block data and initial seed values. This action ensures reproducible results by using block information and provided seed values to compute new seeds for event execution.

## Implementation Details
This action implements the `cc.ActionRunner` interface and is registered under the name `block-event-seed-generation`. It reads block data from an input data source, processes this data along with provided seed values, and outputs computed seeds to a designated output data source.

## Process Flow
1. Initialize the action
2. Extract required attributes (initial_event_seed, initial_realization_seed, plugins)
3. Read and parse block data from input data source named "blockfile"
5. Compute seeds using the BlockModel.Compute method
6. Marshal computed results to JSON
7. Write results to output data source named "seeds"

## Configuration

### Environment
- Requires access to input/output data sources
- Must have proper permissions for reading block data and writing seed data

### Attributes

### Action
| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `initial_event_seed` | int | Yes | Starting seed value for event-level randomization |
| `initial_realization_seed` | int | Yes | Initial seed value for realization generation   |
| `plugins` | string array | Yes | List of plugin names that will receive seeds |

### Global
- No global attributes required

### Input Cofiguration
- Block data from "blockfile" data source

### Input Data Sources
- **blockfile**: Contains block data for seed computation. The configuration requires a data source named `blockfile` and it should have a single **Path** with a key value of **default**

### Output Configuration

### Output Data Sources
- `seeds` (DataSourceName): Contains computed seed data and has a single **Path** with a key value of **default**

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
JSON array containing seed information for each block and plugin combination

### Data Structures
the seed information is structured as a map of seedsets.

### Seed Set Fields
- `event_seed`: Random seed value for event-level randomization
- `block_seed`: Random seed value for block-level randomization
- `realization_seed`: Random seed value for realization-level randomization

## Error Handling
Errors are logged to the compute environment and processing will stop on error