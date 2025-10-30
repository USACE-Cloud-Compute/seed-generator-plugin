# Realization Seed Generation Action

## Description
The Realization Seed Generation Action is responsible for generating seed data based on a provided realization model. It reads input data from a designated data source, processes it using the realization model's compute function, and outputs the generated seed data to another data source.

## Implementation Details
This action implements the `cc.ActionRunner` interface and is registered under the name "realization-seed-generation".

## Process Flow
1. Action is invoked with a payload containing a single output
2. Input data is retrieved from "seedgenerator" data source
3. Input data is decoded into a `RealizationModel`
4. Event identifier is converted to an integer for indexing
5. `RealizationModel.Compute()` is called with the event index
6. Computed result is marshaled to JSON
7. Output is written to "seeds" data source

## Configuration

### Environment
- No specific environment variables required

---
### Attributes
### Action
- No action attribute configuration required
### Global
- No global attribute configuration required


## Input Configuration

### Input Data Sources
- **seedgenerator**: Contains realization model data for seed computation. The configuration requires a data source named `seedgenerator` and it should have a single **Path** with a key value of **default**

## Output Configuration
- Single output definition with path key `default`

### Output Data Sources
- `seeds`: Data source for storing generated seed data and has a single path key: **default**

## Configuration Examples

### Basic Configuration
```json
{
  "action": "realization-seed-generation",
  "inputs": [
    {
      "name": "seedgenerator",
      "paths": {
        "default": "conformance/simulations/seed-model.json"
      },
      "store_name": "FFRD"
    }
  ],
  
  "outputs": [
    {
      "name": "seeds",
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
JSON array containing seed information for each block and plugin combination

### Data Structures
the seed information is structured as a map of seedsets.

### Seed Set Fields
- `event_seed`: Random seed value for event-level randomization
- `block_seed`: Random seed value for block-level randomization
- `realization_seed`: Random seed value for realization-level randomization

## Error Handling
Errors are logged to the compute environment and processing will stop on error