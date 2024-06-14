Feature: Data Product publish

Scenario:
    Given 已開啟服務 nats
    Given 已開啟服務 dispatcher

# Scenario
    Scenario: 針對data product 的 event publish，成功情境
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    When 更新 data product "'drink'" 使用參數 enabled=true
    When 更新 data product "'drink'" 的 ruleset "'drinkEvent'" 使用參數 enabled=true
    Then 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |   ID   | Event      |                    Payload                             | RSMethod |       RSHandler         |       RSSchema                    | RSPk     | RSEnabled | DPEnabled |
        |  M(1)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(2)  | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(3)  | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(4)  | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(5)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(6)  | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(7)  | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(8)  | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(9)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(10) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(11) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(12) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(13) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(14) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(15) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(16) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(17) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(18) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(19) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(20) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(21) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(22) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(23) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(24) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(25) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | ""       | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(26) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | " "      | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(27) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | _-*($)?@ | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(28) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(29) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/WrongType.js   | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(30) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/ExtraColumn.js | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(31) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   |      [ignore]           |          ./assets/schema.json     |    id    | [true]    | [true]    |
        |  M(32) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schemaWrongType.json     | id       | [true]    | [true]    |
        |  M(33) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(34) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,      | [true]    | [true]    |
        |  M(35) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [false]   |

    Scenario: 與成功情境大致相同，但publish完畢後先更新ruleset為enabled，後更新data product為enabled
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    When 更新 data product "'drink'" 使用參數 enabled=true
    When 更新 data product "'drink'" 的 ruleset "'drinkEvent'" 使用參數 enabled=true
    Then 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  E1(1) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [false]   | [false]   |
    
    Scenario: 針對data product 的 event publish(publish指令失敗情境)
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    Then Cli 回傳建立失敗
    # And 應有錯誤訊息 "'<Error_message>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  E2(1) | ""         | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        |  E2(2) | drinkEvent | ''                                             | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        |  E2(3) | drinkEvent | ' '                                            | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        |  E2(4) | drinkEvent | 'abc'                                          | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |

    Scenario: 指令執行成功但沒publish到指定的DP
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    Then 查詢 GVT_default_DP_"'drink'" 裡沒有 "'<Event>'"
    Then 使用 nats jetstream 查詢 GVT_default "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |    ID   |    Event     |                    Payload                     | RSMethod |         RSHandler                                  |               RSSchema              |    RSPk    | RSEnabled |    DPEnabled    |
        |  E3(1)  | drinkEvent   |               '{"name":"test"}'                |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(2)  | drinkEvent   |         '{"name":"test","not_exist":"Hi"}'     |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(3)  | drinkEvent   |             '{"not_exist":"Hi"}'               |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(4)  | drinkEvent   |                     '{}'                       |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(5)  | drinkEvent   |               '{"name":"test"}'                |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(6)  | drinkEvent   |         '{"name":"test","not_exist":"Hi"}'     |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(7)  | drinkEvent   |             '{"not_exist":"Hi"}'               |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(8)  | drinkEvent   |                     '{}'                       |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(9)  | NotExist     | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        #|  E3(10) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/empty.js                             |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(11) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/fail_handler.js                      |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(12) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/WrongFileExtension.txt               |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(13) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/WrongFileExtensionAndFormat.jpg      |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(14) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    ""      |  [true]   |      [true]     |
        |  E3(15) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    " "     |  [true]   |      [true]     |
        |  E3(16) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |  _-*($)?@  |  [true]   |      [true]     |
        |  E3(17) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [false]  |      [true]     |
        |  E3(18) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [false]    |
        |  E3(19) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schemaWithNoPK.json      | id       | [true]    | [true]    |
        |  E3(20) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schemaWithEmptyJson.json | id       | [true]    | [true]    |

    Scenario: 針對data product 的 event publish (同個event pub到多個data product中)
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given 創建 data product "'drink2'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    Given "'drink2'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    Then 查詢 GVT_default_DP_"'drink'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Then 查詢 GVT_default_DP_"'drink2'" 裡有 "'<Event>'" 帶有 "'<Payload>'"
    Examples:
        |    ID   |    Event   |                    Payload                     | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E4(1)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |

    Scenario: 針對data product 的 event publish (連續publish兩筆帶有相同PK值，但其他欄位數量與內容皆不相同的event)
    Given 創建 data product "'drink'" 使用參數 "'<DPEnabled>'"
    Given "'drink'" 創建 ruleset "'drinkEvent'" 使用參數 "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload>'"
    When publish Event "'<Event>'" 使用參數 "'<Payload2>'"
    Then 查詢 GVT_default_DP_"'drink'" 裡是否有兩筆 "'<Event>'" 帶有 "'<Payload>'" 與 "'<Payload2>'"
    Examples:
        |    ID   |    Event   |                    Payload                             |         Payload2                 | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E5(1)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | '{"id":1,"uid":222,"name":"replace"}' |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |