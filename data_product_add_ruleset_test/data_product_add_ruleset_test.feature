Feature: Data Product ruleset add

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario
    Scenario: 針對data product加入ruleset，成功情境
    Given 已有data product "'drink'" enabled "'[true]'"
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then ruleset 創建成功
    Then 使用gravity-cli 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在
    Then 使用nats jetstream 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在，且參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'" 正確
    Examples:
        |  ID   | ProductName | Ruleset       | Method       | Event         |   Pk          |  Desc             |          Handler_script       |           Schema               | Enabled |
        | M(1)  | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(2)  | drink       | _-*=_?+@      | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(3)  | drink       | 中文          |  create      | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(4)  | drink       | [a]x32768     | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(5)  | drink       | drinkCreated  | " "          | " "           |   " "         |  " "              |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(6)  | drink       | drinkCreated  | _-*=_?+@     | _-*=_?+@      |   _-*=_?+@    |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(7)  | drink       | drinkCreated  | 中文         | 中文           |     中文      |   description     |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(8)  | drink       | drinkCreated  | [a]x32768    | [a]x32768     |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(9)  | drink       | drinkCreated  | "create"     | "drinkCreated"|   "id"        | "drink data desc" |     "./assets/handler.js"     |      "./assets/schema.json"    | [true]  |
        | M(10) | drink       | drinkCreated  | create       | drinkCreated  |      id,id2   | description       |     ./assets/handler.js       |      "./assets/schema.json"    | [true]  |
        | M(11) | drink       | drinkCreated  | create       | drinkCreated  |   ""          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(12) | drink       | drinkCreated  | create       | drinkCreated  |   [ignore]    |    [ignore]       |             [ignore]          |            [ignore]            | [true]  |
        | M(13) | drink       | drinkCreated  | create       | drinkCreated  |   [a]x32768   |   [a]x32768       |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(14) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  ""               |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(15) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [false] |
        | M(16) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [ignore]|

#Scenario
    Scenario: 針對data product加入ruleset，重複建立情境
    Given 已有data product "'drink'" enabled "'[true]'"
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then ruleset 創建成功
    Then 使用gravity-cli 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在
    Then 使用nats jetstream 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在，且參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'" 正確
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then ruleset 創建失敗
    # And 應有錯誤訊息 "'<Error_message>'"
    Examples:
        | ID   | ProductName | Ruleset       | Method  | Event         |   Pk  | Desc             |   Handler_script   |       Schema         | Enabled |  Error_message             |
        | E1(1)| drink       | drinkCreated  | create  | drinkCreated  |   id  |  "description"   |./assets/handler.js | ./assets/schema.json | [true]  |                            |

#Scenario
    Scenario Outline: 針對data product加入ruleset，失敗情境
    Given 已有data product "'drink'" enabled "'[true]'"
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then ruleset 創建失敗
    # And 應有錯誤訊息 "'<Error_message>'"
    Examples:
        | ID     | ProductName | Ruleset       | Method       | Event         |   Pk     |  Desc            |        Handler_script       |              Schema               | Enabled | Error_message |
        | E2(1)  | NotExists   |  drinkCreated |  create      | drinkCreated  |   id     |   description    |     ./assets/handler.js     |      ./assets/schema.json         | [true]  |               |
        | E2(2)  |  [null]     |    [null]     | [ignore]     |   [ignore]    | [ignore] |   [ignore]       |           [ignore]          |           [ignore]                | [true]  |               |
        | E2(3)  | drink       |    [null]     | [ignore]     |   [ignore]    | [ignore] |   [ignore]       |           [ignore]          |           [ignore]                | [true]  |               |
        | E2(4)  | drink       | drinkCreated  | ""           | ""            |   ""     |  ""              |     ./assets/handler.js     |      ./assets/schema.json         | [true]  |               |
        | E2(5)  | drink       | drinkCreated  | 中文         |  [ignore]     |   id     |   description    |     ./assets/handler.js     |      ./assets/schema.json         | [true]  |                |
        | E2(6)  | drink       | drinkCreated  | [ignore]     | drinkCreated  |   id     |   description    |     ./assets/handler.js     |     ./assets/schema.json          | [true]  |               |
        | E2(7)  | drink       | drinkCreated  | [ignore]     |  [ignore]     |   id     |   description    |     ./assets/handler.js     |     ./assets/schema.json          | [true]  |               |
        | E2(8)  | drink       | drinkCreated  | create       |  [ignore]     |   id     |   description    |     ./assets/handler.js     |     ./assets/schema.json          | [true]  |               |
        | E2(9)  | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/not_exist.js   |      ./assets/schema.json         | [true]  |               |  
        | E2(10) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |                 ""          |      ./assets/schema.json         | [true]  |               |  
        | E2(11) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ""                      |      ""                           | [true]  |               |
        | E2(12) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/handler.js     |            not_exist.json         | [true]  |               |  
        | E2(13) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/handler.js     |      ./assets/fail_schema.json    | [true]  |               |