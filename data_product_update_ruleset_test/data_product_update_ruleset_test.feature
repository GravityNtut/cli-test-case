Feature: Data Product ruleset update

Scenario:
Given 已開啟服務 nats
Given 已開啟服務 dispatcher
#Scenario
	Scenario: 針對更新data product ruleset成功情境
	Given 已有 data product "'drink'" enabled "'[true]'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'" enabled "'<GivenRSEnabled>'"
	When 更新 dataproduct "'<ProductName>'" ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Then Cli 回傳更改成功
	And 使用 nats jetstream查詢 "'drink'" 的 "'drinkCreated'" 參數更改成功 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Examples:
    |  ID  | ProductName | Ruleset       | Event         | Method    | 		Schema          | 		Handler_script	   | Pk       | Desc          | Enabled  | GivenRSEnabled |
	| M(1) | drink       | drinkCreated  | drinkCreated  | create    | ./assets/schema.json |  ./assets/handler.js     | id       |  description  | [true]   |   [ignore]   |
	| M(2) | drink       | drinkCreated  | drinkCreated  | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    #單獨update method會跳Error: Invalid method
    | M(3) | drink       | drinkCreated  | [ignore]      | create    | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    | M(4) | drink       | drinkCreated  | [ignore]      | [ignore]  |./assets/schema.json  | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    | M(5) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 	./assets/handler.js    | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    | M(6) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | id       | [ignore]      | [ignore] |   [ignore]   |
	| M(7) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | "id, num"| [ignore]      | [ignore] |   [ignore]   |
    | M(8) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | ""       | [ignore]      | [ignore] |   [ignore]   |
    | M(9) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | description   | [ignore] |   [ignore]   |
    | M(10)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | ""            | [ignore] |   [ignore]   |
    | M(11)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | " "           | [ignore] |   [ignore]   |
    | M(12)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [true]   |   [ignore]   |
	| M(13)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [false]  |   [true]     |
	| M(14)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |

#Scenario
	Scenario: 針對更新data product ruleset失敗情境
	Given 已有 data product "'drink'" enabled "'[true]'"
    Given 已有 data product 的 ruleset "'drink'" "'drinkCreated'" enabled "'[ignore]'"
	Given 儲存 nats 現有 data product ruleset 副本 "'drink'" "'drinkCreated'" 
	When 更新 dataproduct "'<ProductName>'" ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Then Cli 回傳更改失敗
	# And 應有錯誤訊息 "'<Error_message>'"
	And 使用 nats jetstream 查詢 data product "'drink'" 的 "'drinkCreated'" 資料無任何改動
	Examples:
	|  ID   | ProductName | Ruleset       | Event         | Method    | 		Schema         	 	 | 		Handler_script	   | Pk       | Desc          | Enabled  | Error_message |
	| E1(1) | [null]	  | [null]		  | [ignore]	  | [ignore]  | 		[ignore]        	 | 		  [ignore]         | [ignore] | [ignore]      | [ignore] | 		         |
	| E1(2) | drink       | [null]		  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         | 
	| E1(3) | [null]      | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(4) | drink       | not_exist	  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(5) | not_exist   | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(6) | drink       | drinkCreated  | drinkCreated  | create    |		"not_exist.json" 	  	 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(7) | drink       | drinkCreated  | drinkCreated  | create    |		"" 						 | 	./assets/handler.js    | id       | description   | [true]   | 		         |	
	| E1(8) | drink       | drinkCreated  | drinkCreated  | create    |./assets/fail_schema.json	 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(9) | drink       | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	not_exist.js		   | id       | description   | [true]   | 		         |
	| E1(10)| drink       | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	""					   | id       | description   | [true]   | 		         |
	| E1(11)| drink  	  | drinkCreated  | ""			  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(12)| drink  	  | drinkCreated  | drinkCreated  | "" 	      |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(13)| drink  	  | drinkCreated  | " "			  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(14)| drink  	  | drinkCreated  | drinkCreated  | " "       |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |