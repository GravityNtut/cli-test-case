Feature: Data Product puiblish 

Scenario:
Given 已開啟服務 nats
Given 已開啟服務 dispatcher
#Scenario
	Scenario: 使用product sub指令來接收已publish到該data product的所有資料成功情境
	Given 已有 data product "'drink'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'"
    Given 已 publish 10 筆 Event
    When 訂閱data product "'<ProductName>'" 使用參數 "'<SubName>'" "'<Partitions>'" "'<Seq>'"
    Then Cli 回傳 "'<Partitions>'" 內 "'<Seq>'" 後所有事件資料
    Examples:
    |  ID  | ProductName | SubName |   Partitions   |      Seq      |
    | M(1) |   drink     |         |      -1        |       1       |
    | M(2) |   drink     |         |      0         |       1       |
    | M(3) |   drink     |         |      200       |       1       |
    | M(4) |   drink     |         |  2147483647    |       1       |
    | M(5) |   drink     |         |  -2147483647   |       1       |
    | M(6) |   drink     |         |    [ignore]    |       1       |
    | M(7) |   drink     |         |   2147483648   |       1       |
    | M(8) |   drink     |         |  -2147483648   |       1       |
    | M(9) |   drink     |         |    131,200     |       1       |
    | M(10) |   drink    |         |      -1        |   4294967295  |
    | M(11)|   drink     |         |      -1        |    [ignore]   |
    | M(12)|   drink     |         |      -1        |   4294967296  |
    | M(13)|   drink     |         |      -1        |       5       |

#Scenario
	Scenario: 使用product sub指令來接收已publish到該data product的所有資料失敗情境
    Given 已有 data product "'drink'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'"
    Given 已 publish 10 筆 Event
    When 訂閱data product "'<ProductName>'" 使用參數 "'<SubName>'" "'<Partitions>'" "'<Seq>'"
    Then Cli 回傳指令失敗
    Examples:
    |   ID   | ProductName | SubName    |   Partitions   |      Seq      |
    |  E1(1) |   notExist  |            |      -1        |       1       |
    |  E1(3) |     drink   |            |    notNumber   |       1       |
    |  E1(4) |     drink   |            |      0.1       |       1       |
    |  E1(5) |     drink   |            |      -1        |       0       |  
    |  E1(6) |     drink   |            |      -1        |      -1       | 
    |  E1(7) |     drink   |            |      -1        |   notNumber   |
    |  E1(8) |     drink   |            |      -1        |       0.1     | 
