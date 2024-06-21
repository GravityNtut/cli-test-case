Feature: Data Product list

Scenario:
    Given NATS has been opened
    Given Dispatcher has been opened
    #Scenario
    Scenario Outline: Success scenario for list of data products
    Given Create "'<ProductAmount>'" data products with "'<ProductName>'" using parameters "'<Description>'" "'<Enabled>'"
    Given Create "'<RulesetAmount>'" rulesets for "'<ProductName>'"
    Given Publish the event "'<EventAmount>'" times
    When Listing all data products using gravity-cli
    Then The CLI returns "'<ProductAmount>'" data products. The final product has the name "'<ProductName>'", with "'<RulesetAmount>'" rulesets, and a total of "'<EventAmount>'" events published. Each data product has a description of "'<Description>'" and an enabled status of "'<Enabled>'"
    Examples:
        |  ID   | ProductName | Description | Enabled | RulesetAmount | EventAmount | ProductAmount |
        |  M(1) | [a]x240     | description | [false] | 0             | 0           | 1             | 
        |  M(2) | drink       |     ""      | [true]  | 1             | 1           | 100           |
        |  M(3) | drink       |     " "     | [true]  | 1             | 100         | 1             | 
        |  M(4) | drink       |  [a]x32768  | [true]  | 0             | 0           | 1             |

#Scenario
    Scenario Outline: Fail scenario for list of data products
    When Listing all data products using gravity-cli
    Then CLI returns exit code 1
    And The error message should be "'<Error_message>'"

    Examples:
        |  ID   | ProductAmount | Error_message         |
        |  M(1) | 0             | No available products |