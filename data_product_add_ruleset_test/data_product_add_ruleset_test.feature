Feature: Data Product ruleset add

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario
    Scenario: 針對data product加入ruleset 成功情境
    Given 已有data product "'drink'"
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'"
    Then ruleset 創建成功
    Then 使用gravity-cli 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在
    Then 使用nats jetstream 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在，且參數 method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'" 正確
    Examples:
        | ProductName | Ruleset       | Method       | Event         |   Pk          |  Desc             |          Handler_script       |           Schema               |
        | drink       | drinkCreated  | "create"     | "drinkCreated"|   "id"        | "description"     |     "./assets/handler.js"     |      "./assets/schema.json"    |
        | drink       | drinkUpdated  | update       | drinkUpdated  |   [ignore]    |   [ignore]        |            [ignore]           |           [ignore]             |
        | drink       | drinkDeleted  | delete       | drinkDeleted  |   [a]x32768   |   [a]x32768       |            [ignore]           |           [ignore]             |
        | drink       | [a]x32768     | delete       | drinkDeleted  |   [a]x32768   |   [a]x32768       |            [ignore]           |           [ignore]             |
        | drink       | drinkCreated  | drinkCreated | drinkCreated  |      id       | "drink_data_desc" |     "./assets/handler.js"     |      ./assets/schema.json      |
        | drink       | drinkCreated  | drinkCreated | drinkCreated  |      id,num   | "drink_data_desc" |     ./assets/handler.js       |      "./assets/schema.json"    |
        | drink       | drinkCreated  | 中文         | 中文           |     中文      |      ""           |     ./assets/handler.js       |      ./assets/schema.json      |
        | drink       | drinkCreated  | _-*=_?+@     | _-*=_?+@      |   _-*=_?+@    | "drink data desc" |     ./assets/handler.js       |      ./assets/schema.json      |
        | drink       | 中文          |  create      |   drinkCreated|   [ignore]    |    [ignore]       |             [ignore]          |            [ignore]            |
        | drink       | _-*=_?+@      |  create      |   drinkCreated|   [ignore]    |    [ignore]       |             [ignore]          |            [ignore]            |

#Scenario
    Scenario: 針對data product加入ruleset 重複建立情境
    Given 已有data product "'drink'"
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'"
    Then ruleset 創建成功
    Then 使用gravity-cli 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在
    Then 使用nats jetstream 查詢 "'<ProductName>'" 的 "'<Ruleset>'" 存在，且參數 method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'" 正確
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'"
    Then ruleset 創建失敗
    Examples:
        | ProductName | Ruleset       | Method  | Event         |   Pk          |  Desc             |         Handler_script        |             Schema                   |
        | drink       | drinkCreated  | create  | drinkCreated  |   id          |   "description"   |     ./assets/handler.js       |      ./assets/schema.json            |

#Scenario
    Scenario Outline: 針對data product加入ruleset 失敗情境
    Given 已有data product "'drink'"
    When "'<ProductName>'" 創建ruleset "'<Ruleset>'" method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'"
    Then ruleset 創建失敗
    And 應有錯誤訊息 "'<Error_message>'"
    Examples:
        | ProductName | Ruleset       | Method       | Event         |   Pk     |  Desc            |        Handler_script       |              Schema               |              Error_message             |
        | NotExists   |  drinkCreated |  create      |  drinkCreated | [ignore] |   [ignore]       |           [ignore]          |          [ignore]                 |                                        |
        |  [null]     |    [null]     | [ignore]     |   [ignore]    | [ignore] |   [ignore]       |           [ignore]          |           [ignore]                |                                        |
        | drink       |    [null]     | [ignore]     |   [ignore]    | [ignore] |   [ignore]       |           [ignore]          |           [ignore]                |                                        |
        | drink       | drinkCreated  | [ignore]     |  [ignore]     |   id     |   "description"  |     ./assets/handler.js     |     ./assets/schema.json          |                                        |
        | drink       | drinkCreated  | [ignore]     | drinkCreated  |   id     |   "description"  |     ./assets/handler.js     |     ./assets/schema.json          |                                        |
        | drink       | drinkCreated  | drinkCreated |  [ignore]     |   id     |   "description"  |     ./assets/handler.js     |     ./assets/schema.json          |                                        |
        | drink       | drinkCreated  | 中文         |  [ignore]     |   id     |   "description"  |     ./assets/handler.js     |      ./assets/schema.json         |                                        |
        | drink       | drinkCreated  | [ignore]     |  [ignore]     |   id     |   "description"  |     ./assets/handler.js     |      ./assets/schema.json         |                                        |  
        | drink       | drinkCreated  | create       | drinkCreated  |   id     |   "description"  |     ./assets/handler.js     |            not_exist.json         |                                        |  
        | drink       | drinkCreated  | create       | drinkCreated  |   id     |   "description"  |     ./assets/handler.js     |      ./assets/fail_schema.json    |                                        |  
        | drink       | drinkCreated  | create       | drinkCreated  |   id     |   "description"  |     ./assets/abc.js         |      ./assets/schema.json         |                                        |  
        | drink       | drinkCreated  | create       | drinkCreated  |   id     |   "description"  |                 ""          |      ./assets/schema.json         |                                        |  
