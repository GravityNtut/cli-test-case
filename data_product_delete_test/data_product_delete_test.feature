Feature: Data Product delete

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher

#Scenario: 
Scenario: 針對已存在的Data Product進行刪除 成功情境
    Given 已有data product "drink"
    When 刪除data product "<ProductName>"
    Then data product 刪除成功
    Then 使用gravity-cli查詢data product 列表 "<ProductName>" 不存在
    Then 使用nats jetstream 查詢 data product 列表 "<ProductName>" 不存在
Examples:
    | ProductName |
    | drink       |

#Scenario: 
Scenario: 針對不存在的Data Product進行刪除 失敗情境
    When 刪除data product "<ProductName>"
    Then data product 刪除失敗
Examples:
    | ProductName |
    | failProduct |
    