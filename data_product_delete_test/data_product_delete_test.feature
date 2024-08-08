Feature: Data Product delete

Background: Check the NATS and Dispatcher
    Given Nats has been opened
    Given Dispatcher has been opened

#Scenario: 
    @M
    Scenario Outline: Success scenario for the deletion of an existing Data Product
    Given  Create data product with "'<ProductName>'" using parameters "'[true]'"
    When Delete data product "'<ProductName>'"
    Then The CLI returned the message: Product "'<ProductName>'" was deleted.
    Then Using gravity-cli to query "'<ProductName>'" shows it does not exist.
    Then Using NATS Jetstream to query "'<ProductName>'" shows it does not exist.
    Examples:
        |  ID  | ProductName |
        | M(1) | drink       |

#Scenario: 
    @E1
    Scenario Outline: Fail scenario for the deletion of a non-existent Data Product
    When Delete data product "'<ProductName>'"
    Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
    Examples:
        |  ID   | ProductName | Error_message |
        | E1(1) | failProduct |               |
        | E1(2) |   [null]    |               |
    