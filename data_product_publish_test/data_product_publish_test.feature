Feature: Data Product publish

Scenario:
    Given 已開啟服務 nats
    Given 已開啟服務 dispatcher

#Scenario
    Scenario: 針對data product 的 event publish，成功情境
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    #When publish 事件使用參數 "'<Event>'" "'<Payload>'"
    #When 更新 data product "'drink'" 使用參數 --enabled
    #When 更新 data product "'drink'" 的 ruleset "drinkCreated" 使用參數 --enabled
    #Then 使用 SDK 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |   ID   | Event      |                    Payload                   | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  M(1)  | drinkEvent | {"id":1,"name":"test","price":100,"kcal":50} | init     | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        #|  M(2)  | drinkEvent | '{"id":1,"name":"test"}'                       | init     | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |