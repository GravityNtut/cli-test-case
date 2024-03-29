Feature: Data Product ruleset delete

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario
    Scenario: 針對刪除data product ruleset成功情境
    Given 已有date product "'drink'"
    Given 已有data product 的 ruleset "'drink'" "'drinkCreated'"
    When 刪除 "'<ProductName>'" 的 ruleset "'<RulesetName>'"
    Then 刪除成功
    Then 使用gravity-cli 查詢 "'<ProductName>'" 的 "'<RulesetName>'" 不存在
    Examples:
        | ProductName | RulesetName     |
        | drink       | drinkCreated    |

#Scenario
    Scenario: 針對刪除data product ruleset失敗情境
    Given 已有date product "'drink'"
    Given 已有data product 的 ruleset "'drink'" "'drinkCreated'"
    When 刪除 "'<ProductName>'" 的 ruleset "'<RulesetName>'"
    Then 刪除失敗
    Examples:
        | ProductName  | RulesetName  |
        | drink        | NotExists    |
        | NotExists    | drinkCreated |
        |   [null]     |    [null]    |
