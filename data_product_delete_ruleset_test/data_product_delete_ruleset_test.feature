Feature: Data Product ruleset delete

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario
    Scenario: 針對刪除data product ruleset成功情境
    Given 已有date product "<productName>"
    Given 已有data product 的 ruleset "<productName>" "<rulesetName>"
    When 刪除 "<productName>" 的 ruleset "<rulesetName>"
    Then 刪除成功
    Then 使用gravity-cli 查詢 "<productName>" 的 "<rulesetName>" 不存在
    Examples:
        | productName | rulesetName     |
        | drink       | drinkCreated    |

#Scenario
    Scenario: 針對刪除data product ruleset失敗情境
    Given 已有date product "<productName>"
    Given 已有data product 的 ruleset "<productName>" "<rulesetName>"
    When 刪除 "<productName2>" 的 ruleset "<rulesetName2>"
    Then 刪除失敗
    Examples:
        | productName | rulesetName     | productName2 | rulesetName2 |
        | drink       | drinkCreated    | drink        | NotExists    |
        | drink       | drinkCreated    | NotExists    | drinkCreated |
