Feature: Data Product ruleset delete

Background: Check the NATS and Dispatcher
    Given Nats has been opened
    Given Dispatcher has been opened

#Scenario
    @M
    Scenario Outline: Success scenario for the deletion of a data product ruleset
    Given Create data product with "'<ProductName>'" using parameters "'[true]'"
    Given Create data product ruleset with "'<ProductName>'", "'<RulesetName>'" using parameters "'[true]'"
    When Delete ruleset "'<RulesetName>'" for data product "'<ProductName>'"
    Then CLI returned successfully deleted
    Then Using gravity-cli to query that "'<RulesetName>'" does not exist for "'<ProductName>'"
    Examples:
        | ID  | ProductName | RulesetName     |
        | M(1)| drink       | drinkCreated    |

#Scenario
    @E1
    Scenario Outline: Fail scenario for the deletion of a non-existent Data Product ruleset
    Given Create data product with "'drink'" using parameters "'[true]'"
    Given Create data product ruleset with "'drink'", "'drinkCreated'" using parameters "'[true]'"
    When Delete ruleset "'<RulesetName>'" for data product "'<ProductName>'"
    Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
    Examples:
        |  ID  | ProductName  | RulesetName  | Error_message |
        | E1(1)| drink        | NotExists    |               |
        | E1(2)| NotExists    | drinkCreated |               |
        | E1(3)|   [null]     |    [null]    |               |
