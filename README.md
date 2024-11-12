![workflow](https://github.com/GravityNtut/cli-test-case/actions/workflows/main.yml/badge.svg)
![Coverage](https://byob.yarr.is/GravityNtut/cli-test-case/coverage)

# Gravity CLI Test case

This repository contains test cases for the Gravity CLI tool.

## File Structure
The file structure follows a pattern where each folder ending with `_test` represents a test case.
Within each test case folder, you'll find:
- A `.feature` file containing the test scenarios and detailed test descriptions
- Supporting files for test execution

For example, in `data_product_add_ruleset_test/`:
- `data_product_add_ruleset_test.feature` describes the test scenarios
- `assets/` directory containing files needed for testing
- `data_product_add_ruleset_test.go` implements the test logic

You can examine the `.feature` files to understand what each test case is designed to verify.
```
.
├── config
│   └── config.json
├── data_product_add_ruleset_test
│   ├── assets
│   ├── command.sh
│   ├── data_product_add_ruleset_test.feature
│   └── data_product_add_ruleset_test.go
├── ...
├── docker-compose.yaml
├── Earthfile
├── go.mod
├── go.sum
├── gravity-cli
```
## Usage

Clone this repository and navigate to the cli-test-case folder.
```shell
git clone git@github.com:GravityNtut/cli-test-case.git
cd cli-test-case
```
Build gravity-cli using Earthly:
```shell
earthly -P +build-cli
```

There are two ways to run cli-test-case:
- Through Earthly.
    ```sh
    earthly -P +ci
    ```
- Through docker compose and go test
    
    1. Start the dependency services:
        ```sh
        docker compose up -d
        ```
    2. Execute the tests:
        ```sh
        go test -p 1 ./...
        ```

    You can also execute specific test cases using go test:
    ```shell
    go test -p 1 ./[test_case_folder_name] -v
    ```

    Optional flags:
    - `-v`: Show verbose output
    - `-t=[tag]`: Run specific scenarios (see [tags documentation](https://github.com/cucumber/godog?tab=readme-ov-file#tags))



