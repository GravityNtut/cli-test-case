Feature: Data Product ruleset add
#Scenario
    Scenario: Extension13 pk與desc皆為空白或null
    
    When "<productName>" 創建ruleset "<ruleset>" method "<method>" event "<event>" pk "<pk>" desc "<desc>" handler "<handler_script>" schema "<schema>"
    Then ruleset 創建失敗

    Examples:
        | productName | ruleset       | method  | event         | pk | desc | handler_script | schema               |
        | drink       | drinkCreated  | create  | drinkCreated  |    |      |     handler.js | schema.json          |
