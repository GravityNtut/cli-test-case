Feature: Data Product create

Scenario:
Given 已開啟服務nats
Given 已開啟服務dispatcher
#Scenario
	Scenario: 使用者使用product create指令來建立data product，成功情境
	When 創建一個data product "<ProductName>" 註解 "<Description>" schema檔案 "<Schema>"
	Then Cli回傳"<ProductName>"建立成功
	Then 使用gravity-cli查詢data product 列表 "<ProductName>" 存在
	Then 使用nats jetstream 查詢 data product 列表 "<ProductName>" 存在
	Examples:
	| ProductName | Description  | Schema      |
	| drink       | description  | schema.json |
	# ProductName max 240 characters
	| [a]x240     | description  | schema.json |
	# desc 32768 characters
	| drink       | [a]x32768    | schema.json |
	| drink       |  			 | schema.json |
	| drink       | [ignore] 	 | schema.json |
	| drink       | description  | [ignore]    | 

#Scenario
	Scenario: 使用者使用product create指令來建立data product，名稱重複
	When 創建一個data product "<ProductName>" 註解 "<Description>" schema檔案 "<Schema>"
	Then Cli回傳"<ProductName>"建立成功
	Then 使用gravity-cli查詢data product 列表 "<ProductName>" 存在
	Then 使用nats jetstream 查詢 data product 列表 "<ProductName>" 存在
	When 創建一個data product "<ProductName2>" 註解 "<Description>" schema檔案 "<Schema>"
	Then Cli回傳建立失敗
	And 應有錯誤訊息 "<Error_message>"
	Examples:
	| ProductName | Description  | Schema      | ProductName2 | Error_message |
	| drink       | description  | schema.json | drink        |			      |

#Scenario
	Scenario: 使用者使用product create指令來建立data product，失敗情境
	When 創建一個data product "<ProductName>" 註解 "<Description>" schema檔案 "<Schema>"
	Then Cli回傳建立失敗
	And 應有錯誤訊息 "<Error_message>"
	Examples:
	| ProductName   | Description  | Schema        | Error_message   |
	| _-*\($\)?@      | description  | schema.json   | 			     |
	| 中文		 	| description   | schema.json  |                 |
	| [null]        | description  | schema.json   |				 |
	|               | description  | schema.json   |				 |
	| drink         | description  | notExist.json |				 |
	| drink         | description  | abc.json      |				 |
	# | drink         | [null]       | schema.json   |				 |
	# | drink         | [null]       | notExist.json |				 |
	# | drink         | [null]       | abc.json	   |				 |
