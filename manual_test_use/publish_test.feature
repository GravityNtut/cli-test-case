Feature: Data Product ruleset add

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario
    Scenario: 針對data product加入ruleset，成功情境
    Given 已有data product "'drink'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'"
    Then ruleset 創建成功
    Then 已 publish "'3'" 筆 Event
    Then 顯示 "'drink'" 資料

    Examples:
        |  ID   | ProductName | Ruleset       |
        | M(1)  | drink       | drinkCreated  |

