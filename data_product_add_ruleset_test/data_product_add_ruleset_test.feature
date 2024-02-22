Feature: Data Product ruleset add

#Scenario
    Scenario: Extension 7 event與method皆為null,space,ignore
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗
    And 應有event與method未填入的錯誤訊息
    Examples:
        | productName | ruleset       | method  | event         | pk |  desc           | handler_script     |      schema               |
        | drink       | drinkCreated  |         |               | id | description     |     handler.js     |      schema.json          |

#Scenario
    Scenario: Extension 8 schema檔案找不到
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗
    And 應有schema檔案找不到的錯誤訊息
    Examples:
        | productName | ruleset       | method  | event         | pk |  desc           | handler_script     |      schema               |
        | drink       | drinkCreated  | create  | drinkCreated  | id | description     |     handler.js     |   not_exist.json          |

#Scenario
    Scenario: Extension 9 schema錯誤json格式
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗
    And 應有schema格式錯誤的錯誤訊息
    Examples:
        | productName | ruleset       | method  | event         | pk |  desc           | handler_script     |      schema               |
        | drink       | drinkCreated  | create  | drinkCreated  | id | description     |     handler.js     | fail_schema.json          |

#Scenario
    Scenario: Extension 10 handler檔案找不到,space或null
    Given 已有data product "<productName>"
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗
    And 應有handler檔案找不到的錯誤訊息
    Examples:
        | productName | ruleset       | method  | event         | pk |  desc           | handler_script | schema               |
        | drink       | drinkCreated  | create  | drinkCreated  | id | description     |     abc.js     | schema.json          |
        | drink       | drinkCreated  | create  | drinkCreated  | id | description     |                | schema.json          |

# #Scenario
#     Scenario: Extension 11 desc中輸入null
#     Given 已有data product "<productName>"
#     When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
#     Then ruleset 創建失敗
#     And 應有desc為空的錯誤訊息
#     Examples:
#         | productName | ruleset       | method  | event         | pk | desc | handler_script | schema               |
#         | drink       | drinkCreated  | create  | drinkCreated  | id |      |     handler.js | schema.json          |

# #Scenario
#     Scenario: Extension 12 pk中輸入空白或null
#     Given 已有data product "<productName>"
#     When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
#     Then ruleset 創建失敗
#     And 應有pk為空的錯誤訊息
#     Examples:
#         | productName | ruleset       | method  | event         | pk | desc           | handler_script | schema               |
#         | drink       | drinkCreated  | create  | drinkCreated  |    |  This is desc  |     handler.js | schema.json          |

# #Scenario
#     Scenario: Extension 13 pk與desc皆為空白或null
#     Given 已有data product "<productName>"
#     When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
#     Then ruleset 創建失敗
#     And 應有pk 和 desc為空的錯誤訊息
#     Examples:
#         | productName | ruleset       | method  | event         | pk | desc | handler_script | schema               |
#         | drink       | drinkCreated  | create  | drinkCreated  |    |      |     handler.js | schema.json          |
