# MeiliJSONImporter

**MeiliJSONImporter** is a tool designed to facilitate the batch import of JSON data into a MeiliSearch instance. It allows for efficient indexing of large JSON datasets by managing data size and API key authentication.

## Features

- **Easy Configuration**: Set up with flags for host, master key, index name, and JSON file path.
- **Batch Importing**: Manages batch sizes automatically for optimal import performance.
- **Command Line Interface**: Fully functional CLI tool to easily integrate into any workflow.

## Installation

To get started with **MeiliJSONImporter**, clone this repository and build the application.

```bash
git clone https://github.com/yourusername/MeiliJSONImporter.git
cd MeiliJSONImporter
go build -o MeiliJSONImporter ./cmd/
```

## Usage

Run the importer using the following command:

```bash
./MeiliJSONImporter --host [meilisearch_host] --master-key [your_master_key] --index [index_name] --json [path_to_json_file] --batch-size [batch_size_in_bytes]
```

### Example

```bash
./MeiliJSONImporter --host http://127.0.0.1:7700 --master-key masterKey --index movies --json ./data/movies.json --batch-size 10485760
```

### Flags

- `--host`: MeiliSearch host URL. Default is `http://127.0.0.1:7700`.
- `--master-key`: Master key for MeiliSearch if your instance uses authentication.
- `--index`: Name of the MeiliSearch index to which the JSON data will be pushed.
- `--json`: Path to the JSON file containing the data to be imported.
- `--batch-size`: Maximum batch size in bytes for pushing documents. Default is 10MB.

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

Distributed under the MIT License. See `LICENSE` for more information.