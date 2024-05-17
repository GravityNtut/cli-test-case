Feature: Data Product create

Scenario:
Given 已開啟服務 nats
Given 已開啟服務 dispatcher
#Scenario
	Scenario: 使用者使用product create指令來建立data product，成功情境
	When 創建 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Then Cli 回傳 "'<ProductName>'" 建立成功
	Then 使用 gravity-cli 查詢 "'<ProductName>'" 存在
	Then 使用 nats jetstream 查詢 "'<ProductName>'" 存在
	Examples:
	|  ID  | ProductName | Description  			| 		Schema         | Enabled |
	| M(1) | drink       | description				| ./assets/schema.json | [true]  |
	| M(2) |[a]x240      | description				| ./assets/schema.json | [true]  |
	| M(3) | drink       |     ""       		    | ./assets/schema.json | [true]  |
	| M(4) | drink       |     " " 	    		    | ./assets/schema.json | [true]  |
	| M(5) | drink       | [ignore] 			    | ./assets/schema.json | [true]  |
	| M(6) | drink       | [a]x32768    	  	    | ./assets/schema.json | [true]  |
	| M(7) | drink       | "drink data description" |"./assets/schema.json"| [true]  |
	| M(8) | drink       | description	 			|   	[ignore]	   | [true]  |
	| M(9) | drink       | description	 			| ./assets/schema.json | [false] |
	| M(10)| drink       | description	 			| ./assets/schema.json | [ignore]|


#Scenario
	Scenario: 使用者使用product create指令來建立data product，名稱重複
	When 創建 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Then Cli 回傳 "'<ProductName>'" 建立成功
	Then 使用 gravity-cli 查詢 "'<ProductName>'" 存在
	Then 使用 nats jetstream 查詢 "'<ProductName>'" 存在
	When 創建 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Then Cli 回傳建立失敗
	# And 應有錯誤訊息 "'<Error_message>'"
	Examples:
	|   ID  | ProductName | Description | 		Schema         | Enabled | Error_message  |
	| E1(1) | drink       | description | ./assets/schema.json | [true]  |			    |

#Scenario
	Scenario: 使用者使用product create指令來建立data product，失敗情境
	When 創建 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Then Cli 回傳建立失敗
	# And 應有錯誤訊息 "'<Error_message>'"
	Examples:
	|   ID  | ProductName   | Description  | 			Schema       			| Enabled | Error_message |
	| E2(1) | _-*\($\)?@    | description  | 		./assets/schema.json 		| [true]  | 			  |
	| E2(2) | 中文		 	 | description  | 	./assets/schema.json  			 | [true]  |               |
	| E2(3) | [null]        | description  |	    ./assets/schema.json  		| [true]  |				  |
	| E2(4) | drink         | description  | 		notExist.json 				| [true]  |				  |
	| E2(5) | drink         | description  | 		./assets/fail_schema.json   | [true]  |				  |
	| E2(6) | drink         | description  | 			""           			| [true]  |			      |
	
