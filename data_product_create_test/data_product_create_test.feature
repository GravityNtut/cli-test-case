Feature: Data Product create

Scenario:
Given 已開啟服務nats
Given 已開啟服務dispatcher
#Scenario
	Scenario: 使用者使用product create指令來建立data product，成功情境
	When 創建 data product "<ProductName>" 使用參數 "<Description>" "<Schema>"
	Then Cli 回傳 "<ProductName>" 建立成功
	Then 使用 gravity-cli 查詢 "<ProductName>" 存在
	Then 使用 nats jetstream 查詢 "<ProductName>" 存在
	Examples:
	|  ID  | ProductName | Description  | Schema      |
	| M(1) | drink       | description  | schema.json |
	| M(2) |[a]x240      | description  | schema.json |
	| M(3) | drink       |              | schema.json |
	| M(4) | drink       | [space] 	    | schema.json |
	| M(5) | drink       | [ignore] 	| schema.json |
	| M(6) | drink       | [a]x32768    | schema.json |
	| M(7) | drink       | description  | [ignore]    |

#Scenario
	Scenario: 使用者使用product create指令來建立data product，名稱重複
	When 創建 data product "<ProductName>" 使用參數 "<Description>" "<Schema>"
	Then Cli 回傳 "<ProductName>" 建立成功
	Then 使用 gravity-cli 查詢 "<ProductName>" 存在
	Then 使用 nats jetstream 查詢 "<ProductName>" 存在
	When 創建 data product "<ProductName>" 使用參數 "<Description>" "<Schema>"
	Then Cli 回傳建立失敗
	And 應有錯誤訊息 "<Error_message>"
	Examples:
	|   ID  | ProductName | Description  | Schema      | Error_message |
	| E1(1) | drink       | description  | schema.json |			    |

#Scenario
	Scenario: 使用者使用product create指令來建立data product，失敗情境
	When 創建 data product "<ProductName>" 使用參數 "<Description>" "<Schema>"
	Then Cli 回傳建立失敗
	And 應有錯誤訊息 "<Error_message>"
	Examples:
	|   ID  | ProductName   | Description  | Schema        | Error_message   |
	| E2(1) | _-*\($\)?@    | description  | schema.json   | 			     |
	| E2(2) | 中文		 	| description   | schema.json  |                 |
	| E2(3) | [null]        | description  | schema.json   |				 |
	| E2(4) |               | description  | schema.json   |				 |
	| E2(5) | drink         | description  | notExist.json |				 |
	| E2(6) | drink         | description  | abc.json      |				 |
	| E2(7) | drink         | description  | [null]        |				 |
	
