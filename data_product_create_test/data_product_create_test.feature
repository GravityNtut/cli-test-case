Feature: Data Product create
#Scenario
	Scenario: 使用者使用product create指令來建立data product
	When 創建一個data product "<ProductName>" 註解 "<Description>" schema檔案 "<Schema>" 
	Then 建立成功
	Then 使用gravity-cli查詢data product 列表 "<ProductName>" 存在
	Then 使用nats jetstream 查詢 data product 列表 "<ProductName>" 存在

Examples:
	| ProductName | Description  | Schema      |
	| drink       | description  | schema.json |
	| food        |   | schema.json |
