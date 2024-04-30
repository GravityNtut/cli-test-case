Feature: Data Product ruleset delete

Scenario:
    Given 已開啟服務 nats
    Given 已開啟服務 dispatcher

#Scenario
    Scenario: 針對刪除data product ruleset成功情境
    Given 已有 date product "'drink'" "'[ignore]'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'" "'[ignore]'"
    When 刪除 "'<ProductName>'" 的 ruleset "'<RulesetName>'"
    Then Cli 回傳刪除成功
    Then 使用 gravity-cli 查詢 "'<ProductName>'" 的 "'<RulesetName>'" 不存在
    Examples:
        | ID  | ProductName | RulesetName     |
        | M(1)| drink       | drinkCreated    |

#Scenario
    Scenario: 針對刪除data product ruleset失敗情境
    Given 已有 date product "'drink'" "'[ignore]'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'" "'[ignore]'"
    When 刪除 "'<ProductName>'" 的 ruleset "'<RulesetName>'"
    Then Cli 回傳刪除失敗
	# And 應有錯誤訊息 "'<Error_message>'"
    Examples:
        |  ID  | ProductName  | RulesetName  | Error_message |
        | E1(1)| drink        | NotExists    |               |
        | E1(2)| NotExists    | drinkCreated |               |
        | E1(3)|   [null]     |    [null]    |               |
