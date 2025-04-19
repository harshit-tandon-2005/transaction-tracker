# transaction-tracker

transaction tracker in Go designed to fetch transaction history based for a specified Ethereum wallet address and exports it to a structured CSV file with relevant transaction details

## Prerequisites

*   [Go](https://golang.org/dl/) installed (version 1.16 or later recommended).

## Setup

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <your-repository-url>
    cd <repository-directory>
    ```

2.  **Install Dependencies:**
    To let Go manage dependencies based on imports, run:
    ```bash
    go mod tidy
    ```
    This command ensures `go.mod` file matches the packages used in the code and downloads them.

3.  **(Optional) Vendor Dependencies:**
    To include dependencies directly in your project repository for offline builds, run:
    ```bash
    go mod vendor
    ```
    This will create a `vendor` directory containing all necessary packages.

4.  **Configure:**
    Create a `config.yml` from `sample_config.yml` file and set the desired configuration values.

## Running the Script

To run the main script, execute:

```bash
go run main.go
```

If you vendored dependencies (Step 3), you might need to build or run using the `-mod=vendor` flag, although `go run` often detects the vendor directory automatically:

```bash
go run -mod=vendor main.go
# or build first
go build -mod=vendor
./<executable-name> # e.g., ./transaction-tracker
```

Once the script is executed the reports would be generated in the project folder under the directory `files/reports`

The naming of the csv files would be as follows:
1. `{{walletAddress}}_external_report.csv`
2. `{{walletAddress}}_internal_report.csv`
1. `{{walletAddress}}_erc-20_report.csv`
1. `{{walletAddress}}_erc-721_report.csv`