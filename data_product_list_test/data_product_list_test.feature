Feature: Data Product list

Scenario:
    Given 已開啟服務 nats
    Given 已開啟服務 dispatcher
#Scenario
    Scenario: 針對data product 的 event list，成功情境
    Given 創建 "'<ProductAmount>'" 個 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'"
    Given 對"'<ProductName>'" 創建 "'<RulesetAmount>'" 個 ruleset
    Given 對Event做 "'<EventAmount>'" 次 publish 
    When 使用gravity-cli 列出所有 data product
    Then Cli 回傳 "'<ProductAmount>'" 個 data product, 每個 data product 裡面的名字為 "'<ProductName>'", 描述內容為 "'<Description>'", Enabled 的狀態為 "'<Enabled>'", Ruleset 的數量為 "'<RulesetAmount>'" 個, 以及 Event 總共發布 "'<EventAmount>'" 個
    Examples:
        |  ID   | ProductName | Description | Enabled | RulesetAmount | EventAmount | ProductAmount |
        |  M(1) | [a]x240     | description | [false] | 0             | 0           | 1             | 
        |  M(2) | drink       | [a]x32768   | [true]  | 0             | 0           | 1             | 
        |  M(3) | drink       | " "         | [true]  | 1             | 500         | 1             | 
        |  M(4) | drink       | ""          | [true]  | 1             | 1           | 100           | 

#Scenario
    Scenario: 針對data product 的 event list，未建立任何data product情境
    When 使用gravity-cli 列出所有 data product
    Then  Cli 回傳建立失敗
    And 應有錯誤訊息 "'<No available products>'"

    Examples:
        |  ID   | ProductAmount |
        |  M(1) | 0             | 