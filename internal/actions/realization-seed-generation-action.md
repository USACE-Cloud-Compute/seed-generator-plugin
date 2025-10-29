# Realization Seed Generation Action

## Description
The Realization Seed Generation Action is responsible for generating seed data based on a provided realization model. It reads input data from a designated data source, processes it using the realization model's compute function, and outputs the generated seed data to another data source.

## Implementation Details
This action implements the `cc.ActionRunner` interface and is registered under the name "realization-seed-generation". It performs the following key operations:
- Retrieves input data from the "seedgenerator" data source
- Decodes the JSON payload into a `RealizationModel`
- Uses the event identifier to compute the realization model
- Marshals the computed result into JSON
- Writes the output to the "seeds" data source

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
- Requires proper data source configuration in the plugin manager

### Attributes
- Action name: `realization-seed-generation`
- Plugin type: `ActionRunner`

### Action
```json
{
  "name": "realization-seed-generation",
  "type": "action",
  "config": {}
}
```

### Global
- No global configuration required

### Inputs
- `eventIdentifier`: String representing the event index for computation
- `outputs`: Array containing exactly one output definition

### Input Data Sources
- `seedgenerator`: Data source containing the realization model
- Path key: `default`

### Outputs
- Single output definition with path key `default`

### Output Data Sources
- `seeds`: Data source for storing generated seed data
- Path key: `default`

## Configuration Examples

### Basic Configuration
```json
{
  "action": "realization-seed-generation",
  "inputs": {
    "eventIdentifier": "42"
  },
  "outputs": [
    {
      "pathKey": "default",
      "dataSourceName": "seeds"
    }
  ]
}
```

### Complete Configuration
```json
{
  "action": "realization-seed-generation",
  "inputs": {
    "eventIdentifier": "100"
  },
  "outputs": [
    {
      "pathKey": "default",
      "dataSourceName": "seeds"
    }
  ],
  "dataSources": {
    "seedgenerator": {
      "type": "file",
      "path": "/data/seed-model.json"
    },
    "seeds": {
      "type": "file",
      "path": "/data/generated-seeds.json"
    }
  }
}
```

## Outputs

### Format
JSON-encoded data representing the computed seed values

### Fields
- `eventIndex`: Integer representing the event identifier used for computation
- `computedData`: Object containing the generated seed data

### Field Definitions
- `eventIndex`: The integer value of the event identifier used for computation
- `computedData`: The result of the realization model computation

## Error Handling
The action handles several error conditions:
- Returns error if more than one output is defined in the payload
- Returns error if input data cannot be read from the "seedgenerator" data source
- Returns error if JSON decoding fails
- Returns error if event identifier cannot be converted to integer
- Returns error if model computation fails
- Returns error if JSON marshaling fails
- Returns error if output cannot be written to the "seeds" data source

## Usage Notes
- The action requires exactly one output definition in the payload
- The event identifier must be convertible to an integer
- Input data must be valid JSON that can be decoded into a `RealizationModel`
- The action assumes the existence of properly configured data sources named "seedgenerator" and "seeds"

## Future Enhancements
- Support for multiple outputs
- Enhanced error reporting with more detailed context
- Configuration options for data source paths
- Support for different data formats beyond JSON
- Batch processing capabilities

## Patterns and Best Practices
- Always validate input parameters before processing
- Use proper error handling with descriptive error messages
- Close readers and writers appropriately
- Maintain consistent data source naming conventions
- Ensure JSON serialization/deserialization errors are handled gracefully
- Use appropriate data types for event identifiers