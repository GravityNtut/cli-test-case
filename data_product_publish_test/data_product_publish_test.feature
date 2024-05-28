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
        |  M(1)  | drinkEvent | '{"id":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        #|  M(2)  | drinkEvent | '{"id":1,"name":"test"}'                       | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |

    Scenario: 與成功情境大致相同，但publish完畢後先更新ruleset為enabled，後更新data product為enabled
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    When 更新 data product "'drink'" 使用參數 enabled=true
    When 更新 data product "'drink'" 的 ruleset "'drinkEvent'" 使用參數 enabled=true
    Then 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  E1(1) | drinkEvent | '{"id":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [false]   | [false]   |
    
    Scenario: 針對data product 的 event publish(publish指令失敗情境)
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    Then Cli 回傳建立失敗
    # And 應有錯誤訊息 "'<Error_message>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  E2(1) | ""         | '{"id":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        # |  E2(2) | drinkEvent | ''                                             | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        # |  E2(3) | drinkEvent | ' '                                            | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        # |  E2(4) | drinkEvent | 'abc'                                          | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |

    Scenario: 指令執行成功但沒publish到指定的DP
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    Then 查詢 GVT_default_DP_"'drink'" 裡沒有 "'<Event>'"
    Then 使用 nats jetstream 查詢 GVT_default "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |    ID   |    Event     |                    Payload                    | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E3(1)  | drinkEvent   |               '{"name":"test"}'               |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |

    Scenario: 針對data product 的 event publish (同個event pub到多個data product中)
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given 創建 data product "'drink2'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    Given "'drink2'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    Then 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Then 查詢 GVT_default_DP_"'drink2'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |    ID   |    Event   |                    Payload                     | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E4(1)  | drinkEvent | '{"id":1,"name":"test","price":100,"kcal":50}' |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |

    Scenario: 針對data product 的 event publish (連續publish兩筆帶有相同PK值，但其他欄位數量與內容皆不相同的event)
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<Event>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload2>'"
    Then 查詢 GVT_default_DP_"'drink'" 裡是否有兩筆 "'<Event>'" 帶有 "'<Payload>'" 與 "'<Payload2>'"
    Examples:
        |    ID   |    Event   |                    Payload                     |         Payload2         | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E5(1)  | drinkEvent | '{"id":1,"name":"test","price":100,"kcal":50}' | '{"id":1,"name":"test"}' |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |