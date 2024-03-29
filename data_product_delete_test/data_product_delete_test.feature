Feature: Data Product delete

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario: 
Scenario: 針對已存在的Data Product進行刪除 成功情境
    Given 已有 data product "'drink'"
    When 刪除 data product "'<ProductName>'"
    Then Cli 回傳 "'<ProductName>'" 刪除成功
    Then 使用 gravity-cli 查詢 "'<ProductName>'" 不存在
    Then 使用 nats jetstream 查詢 "'<ProductName>'" 不存在
Examples:
    |  ID  | ProductName |
    | M(1) | drink       |

#Scenario: 
Scenario: 針對不存在的Data Product進行刪除 失敗情境
    When 刪除 data product "'<ProductName>'"
    Then Cli 回傳刪除失敗
Examples:
    |  ID   | ProductName |
    | E1(1) | failProduct |
    | E2(2) |   [null]    |
    