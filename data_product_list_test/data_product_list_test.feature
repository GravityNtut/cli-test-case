Feature: Data Product list

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher
    Given 已有gravity cli 工具  
#Scenario
    Scenario: 針對data product 的 event list，成功情境
    When 創建 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'"
    Then Cli 回傳 "'<ProductName>'" 建立成功
    When "'<ProductName>'" 創建 "'<RulesetAmount>'" 個 ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'"
    Then ruleset 創建成功
    When 對 "'<Event>'" 做 "'<eventAmount>'" 次 publish "'<Payload>'"
    Then publish 成功
    When 使用gravity-cli 列出所有 data product
    Then 回傳 data product ProductName = "'<ProductName>'", Description = "'<Description>'", Enabled="'<Enabled>'", RulesetAmount="'<RulesetAmount>'", EventAmount="'<EventAmount>'"
    Examples:
        |  ID   | ProductName | Description   | Enabled | RulesetAmount | EventAmount | ProductAmount |
        | M(1)  | [a]x256     | description   | [true]  | 0             | 0           | 1             |
        | M(2)  | drink       | ""            | [true]  | 0             | 1           | 1             |
        | M(3)  | drink       | " "           | [true]  | 0             | 5           | 1             |
        | M(4)  | drink       | [a]x32768     | [true]  | 1             | 0           | 1             |
        | M(5)  | drink       | ""            | [true]  | 1             | 1           | 1             |
        | M(6)  | drink       | ""            | [true]  | 1             | 5           | 1             |
        | M(7)  | drink       | ""            | [true]  | 5             | 0           | 1             |
        | M(8)  | drink       | ""            | [true]  | 5             | 1           | 1             |
        | M(9)  | drink       | ""            | [true]  | 5             | 5           | 1             |
        | M(10) | drink       | ""            | [false] | 0             | 0           | 1             |
        | M(11) | drink       | ""            | [false] | 0             | 1           | 1             |
        | M(12) | drink       | ""            | [false] | 0             | 5           | 1             |
        | M(13) | drink       | ""            | [false] | 1             | 0           | 1             |
        | M(14) | drink       | ""            | [false] | 1             | 1           | 1             |
        | M(15) | drink       | ""            | [false] | 1             | 5           | 1             |
        | M(16) | drink       | ""            | [false] | 5             | 0           | 1             |
        | M(17) | drink       | ""            | [false] | 5             | 1           | 1             |
        | M(18) | drink       | ""            | [false] | 5             | 5           | 1             |