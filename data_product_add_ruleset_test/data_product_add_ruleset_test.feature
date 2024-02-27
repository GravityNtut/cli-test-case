Feature: Data Product ruleset add

#Scenario
    Scenario: 針對data product加入ruleset 成功情境
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建成功
    Then 使用gravity-cli 查詢 "<productName>" 的 "<ruleset>" 成功
    Examples:
        | productName | ruleset       | method  | event         |   pk          |  desc           | handler_script       |      schema                 |
        | drink       | drinkCreated  | create  | drinkCreated  |   id          | description     |     handler.js       |      schema.json            |
        | drink       | drinkUpdated  | update  | drinkUpdated  |   [ignore]    | [ignore]        |     [ignore]         |      [ignore]               |
        | drink       | drinkDeleted  | delete  | drinkDeleted  |   [a]x256     | [a]x4000        |     [ignore]         |      [ignore]               |

#Scenario
    Scenario: 針對data product加入ruleset 重複建立情境
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建成功
    Then 使用gravity-cli 查詢 "<productName>" 的 "<ruleset>" 成功
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗
    Examples:
        | productName | ruleset       | method  | event         |   pk          |  desc           | handler_script       |      schema                 |
        | drink       | drinkCreated  | create  | drinkCreated  |   id          | description     |     handler.js       |      schema.json            |

#Scenario
    Scenario Outline: 針對data product加入ruleset 失敗情境
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗
    And 應有錯誤訊息 "<error_message>"
    Examples:
        | productName | ruleset       | method  | event         |   pk     |  desc           | handler_script     |      schema               |              error_message             |
        | drink       | drinkCreated  |         |               |   id     | description     |     handler.js     |      schema.json          |                                        |  
        | drink       | drinkCreated  | create  | drinkCreated  |   id     | description     |     handler.js     |   not_exist.json          |                                        |  
        | drink       | drinkCreated  | create  | drinkCreated  |   id     | description     |     handler.js     | fail_schema.json          |                                        |  
        | drink       | drinkCreated  | create  | drinkCreated  |   id     | description     |     abc.js         | schema.json               |                                        |  
        | drink       | drinkCreated  | create  | drinkCreated  |   id     | description     |                    | schema.json               |                                        |  
        # | drink       | drinkCreated  | create  | drinkCreated  |          |                 |     handler.js     | schema.json               |                                        |  
        # | drink       | drinkCreated  | create  | drinkCreated  |          |  This is descvv |     handler.js     | schema.json               |                                        |  
        # | drink       | drinkCreated  | create  | drinkCreated  |          |    [null]       |     handler.js     | schema.json               |                                        |  
