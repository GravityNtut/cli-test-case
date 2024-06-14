Feature: Data Product delete

Scenario:
    Given Nats has been opened
    Given Dispatcher has been opened

#Scenario: 
Scenario: Success scenario for the deletion of an existing Data Product
    Given  Create data product with "'<ProductName>'" using parameters "'[true]'"
    When Delete data product "'<ProductName>'"
    Then The CLI returned the message: Product "'<ProductName>'" was deleted.
    Then Using gravity-cli to query "'<ProductName>'" shows it does not exist.
    Then Using NATS Jetstream to query "'<ProductName>'" shows it does not exist.
Examples:
    |  ID  | ProductName |
    | M(1) | drink       |

#Scenario: 
Scenario: Fail scenario for the deletion of a non-existent Data Product.
    When Delete data product "'<ProductName>'"
    Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
Examples:
    |  ID   | ProductName | Error_message |
    | E1(1) | failProduct |               |
    | E2(2) |   [null]    |               |
    