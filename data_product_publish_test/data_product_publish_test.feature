Feature: Data Product publish

Scenario:
    Given 已開啟服務 nats
    Given 已開啟服務 dispatcher

#Scenario
    Scenario: 針對data product 的 event publish，成功情境
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    When 更新 data product "'drink'" 使用參數 enabled=true
    When 更新 data product "'drink'" 的 ruleset "'drinkEvent'" 使用參數 enabled=true
    Then 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  M(1)  | drinkEvent | '{"id":1,"name":"test","price":100,"kcal":50}' | init     | ./assets/handler.js | ./assets/schema.json | id   | [false]    | [false]    |
        #|  M(2)  | drinkEvent | '{"name":"test"}'                       | init     | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |

    Scenario: 指令執行成功但沒publish到指定的DP
        Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
        Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
        When publish Event "'<Event>'" 使用參數 "'<Payload>'"
        Then 查詢 GVT_default_DP_"'drink'" 裡沒有 "'<Event>'" 帶有 "'<Payload>'"
        Then 使用 nats jetstream 查詢 GVT_default "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |    ID   |    Event     |                    Payload                    | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E3(1)  | drinkEvent   |            '{"name":"test"}'           |   init   |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |

