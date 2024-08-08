Feature: Data Product purge

Background: Check the NATS and Dispatcher
    Given Nats has been opened
    Given Dispatcher has been opened

#Scenario: 
    @M
    Scenario Outline: Success scenario for the deletion of an existing Data Product
    Given  Create data product with "'<ProductName>'" using parameters "'[true]'"
    Given Create data product ruleset with "'<ProductName>'", "'drink_ruleset'" using parameters "'[true]'" 
    Given Publish an Event
    When Purge data product "'<ProductName>'"
    Then Check purging data product success
    Then Check data product "'<ProductName>'"'s Events amount is "'0'"
    Then Use NATS JetStream to query the Messages amount of the data product "'<ProductName>'" to be "'0'"
    Examples:
        |  ID  | ProductName |
        | M(1) | drink       |

#Scenario: 
    @E1
    Scenario Outline: The purge of a non-existent Data Product
    Given  Create data product with "'drink'" using parameters "'[true]'"
    Given Create data product ruleset with "'drink'", "'drink_ruleset'" using parameters "'[true]'" 
    Given Publish an Event
    When Purge data product "'<ProductName>'"
    Then CLI returns exit code 1
	And The error message should be "'<Error_message>'"
    Examples:
        |  ID   | ProductName |   Error_message   |
        | E1(1) | failProduct | product not found |
        | E1(2) |     ""      | product not found |
    