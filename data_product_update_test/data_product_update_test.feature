Feature: Data Product update

Scenario:
Given 已開啟服務 nats
Given 已開啟服務 dispatcher
#Scenario
	Scenario: 使用者使用product update指令更新data product，成功情境
	Given 已有 data product "'drink'" "'<ProductEnabled>'"
	When 更新 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'" "'<Schema>'"
	Then Cli 回傳更改成功
	And 使用 nats jetstream 查詢 "'drink'" 參數更改成功 "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Examples:
	|  ID  | ProductName | Description  |		 Schema     	| Enabled   | ProductEnabled |
	| M(1) | drink       | [ignore]     |		[ignore]   		| [ignore]  |    [ignore]    |
	| M(2) | drink       | description  |		[ignore] 		| [ignore]  |    [ignore]    |
	| M(3) | drink       | ""  			|		[ignore]   		| [ignore]  |    [ignore]    |
	| M(4) | drink       | " "    	    |	    [ignore]   	    | [ignore]  |    [ignore]    |
	| M(5) | drink       | [ignore]     | ./assets/schema.json  | [ignore]  |    [ignore]    |
	| M(6) | drink       | [ignore]     |		[ignore]		| [true]    |    [ignore]    |
	| M(7) | drink       | [ignore]     |		[ignore]		| [false]   |    [true]      |
	| M(8) | drink       |  description | ./assets/schema.json  | [true]    |    [ignore]    |

#Scenario
	Scenario: 使用者使用product update指令更新data product，失敗情境
	Given 已有 data product "'drink'" "'[ignore]'"
	Given 儲存 nats 現有 data product 副本 "'drink'"
	When 更新 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'" "'<Schema>'"
	Then Cli 回傳更改失敗
	# And 應有錯誤訊息 "'<Error_message>'"
	And 使用 nats jetstream 查詢 "'drink'" 參數無任何改動
	Examples:
	|  ID   | ProductName | Description  | 	    Schema         		 | Enabled   | Error_message   |
	| E1(1) | not_exist   | [ignore]  	 | 		[ignore] 			 | [ignore]  |                 |
	| E1(2) | [null]      | [ignore] 	 |		[ignore]  		     | [ignore]  | 			       |
	| E1(3) | drink		  | [ignore]     | 		""					 | [ignore]  |                 |
	| E1(4) | drink		  | [ignore]     | ./assets/fail_schema.json | [ignore]  |                 |
	| E1(5) | drink		  | [ignore]     |		not_exist.json		 | [ignore]  |                 |