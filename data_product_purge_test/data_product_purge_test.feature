Feature: Data Product purge

Scenario:
    Given Nats has been opened
    Given Dispatcher has been opened

#Scenario: 
Scenario: Success scenario for the deletion of an existing Data Product
    Given  Create data product with "'<ProductName>'" using parameters "'[true]'"
    Given Create data product ruleset with "'<ProductName>'", "'drink_ruleset'" using parameters "'drinkEvent'", "'[true]'" 
    Given Publish Event "'drinkEvent'" using parameters "''{"id":1,"uid":1,"name":"test","price":100,"kcal":50}''"
    Then Check data product "'<ProductName>'"'s Events amount is "'1'"
    Then Use NATS JetStream to query the Messages amount of the data product "'<ProductName>'" to be "'1'"
    When Purge data product "'<ProductName>'"
    Then Check data product "'<ProductName>'"'s Events amount is "'0'"
    Then Use NATS JetStream to query the Messages amount of the data product "'<ProductName>'" to be "'0'"
Examples:
    |  ID  | ProductName |
    | M(1) | drink       |

#Scenario: 
Scenario: The purge of a non-existent Data Product.
    When Purge data product "'<ProductName>'"
    Then CLI returns exit code 1
	And The error message should be "'<Error_message>'"
Examples:
    |  ID   | ProductName | Error_message     |
    | E1(1) | failProduct | product not found |
    | E2(2) |   [null]    | requires at least 1 arg(s), only received 0 |
    