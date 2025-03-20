# dns-rewrite-go

A lightweight DNS rewriting tool built in Go. It intercepts and rewrites DNS queries based on your custom rules—ideal for testing, custom network setups, or experimenting with DNS behavior.

## Features

- **Custom Rules:** Easily define rewrite rules.
- **Fast & Efficient:** Leverages Go’s concurrency.
- **Simple Setup:** Minimal configuration for quick deployment.

## Installation

1. **Clone the repo:**
   ```bash
   git clone https://github.com/Morl0ck/dns-rewrite-go.git
   cd dns-rewrite-go
   ```

2. **Build the binary:**
   ```bash
   go build -o dns-rewrite-go
   ```

## Usage

Run the tool with a configuration file:

```bash
./dns-rewrite-go
```

Sample config.json

```json
{
  "rewrite_entries": {
    "example.com.": "192.168.1.100",
    "www.example.com.": "192.168.1.100",
    "test.local.": "192.168.1.200"
  },
  "listen_address": "0.0.0.0:53",
  "upstream_dns": "8.8.8.8:53"
}
```

## Perspectives & Solutions

- **Network Engineers:** Tailor DNS responses to fit custom network requirements.
- **Developers:** Explore Go’s networking capabilities and extend functionality.
- **Hobbyists:** Experiment with DNS rewriting for learning or fun modifications.

## Contributing

Contributions are welcome! Open an issue or submit a pull request via the GitHub issues page.

## License

Released under the MIT License. See the [LICENSE](./LICENSE) file for details.
