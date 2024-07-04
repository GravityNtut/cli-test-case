Feature: Data Product update

Scenario:
Given 已開啟服務 nats
Given 已開啟服務 dispatcher
#Scenario
	Scenario: 使用者使用product update指令更新data product，成功情境
	Given 已有 data product "'drink'" enabled "'<GivenDPEnabled>'"
	When 更新 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'" "'<Schema>'"
	Then Cli 回傳更改成功
	And 使用 nats jetstream 查詢 "'drink'" 參數更改成功 "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Examples:
	|  ID  | ProductName | Description  |		 Schema     	| Enabled   | GivenDPEnabled |
	| M(1) | drink       | [ignore]     |		[ignore]   		| [ignore]  |   [false]    |
	| M(2) | drink       | description  |		[ignore] 		| [ignore]  |   [false]    |
	| M(3) | drink       | ""  			|		[ignore]   		| [ignore]  |   [false]    |
	| M(4) | drink       | " "    	    |	    [ignore]   	    | [ignore]  |   [false]    |
	| M(5) | drink       | [ignore]     | ./assets/schema.json  | [ignore]  |   [false]    |
	| M(6) | drink       | [ignore]     |		[ignore]		| [true]    |   [false]    |
	| M(7) | drink       | [ignore]     |		[ignore]		| [false]   |   [true]     |
	| M(8) | drink       |  description | ./assets/schema.json  | [true]    |   [false]    |

#Scenario
	Scenario: 使用者使用product update指令更新data product，失敗情境
	Given 已有 data product "'drink'" enabled "'[true]'"
	Given 儲存 nats 現有 data product 副本 "'drink'"
	When 更新 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'" "'<Schema>'"
	Then Cli 回傳更改失敗
	# And 應有錯誤訊息 "'<Error_message>'"
	And 使用 nats jetstream 查詢 "'drink'" 參數無任何改動
	Examples:
	|  ID   | ProductName | Description  | 	    Schema         		 | Enabled   | Error_message   |
	| E1(1) | not_exist   | [ignore]  	 | 		[ignore] 			 | [false]  |                 |
	| E1(2) | [null]      | [ignore] 	 |		[ignore]  		     | [false]  | 			       |
	| E1(3) | drink		  | [ignore]     | 		""					 | [false]  |                 |
	| E1(4) | drink		  | [ignore]     | ./assets/fail_schema.json | [false]  |                 |
	| E1(5) | drink		  | [ignore]     |		not_exist.json		 | [false]  |                 |