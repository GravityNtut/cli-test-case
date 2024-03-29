Feature: Data Product ruleset update

Scenario:
Given 已開啟服務nats
Given 已開啟服務dispatcher
#Scenario
	Scenario: 針對更新data product ruleset成功情境
	Given 已有data product "'drink'"
    Given 已有data product 的 ruleset "'drink'" "'drinkCreated'" 
	When "'<ProductName>'" 更新ruleset "'<Ruleset>'" 參數 method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'" enabled "'<Enabled>'"
	Then ruleset 更改成功
	And 使用nats驗證 data product "'drink'" 的 ruleset "'drinkCreated'" 更改成功 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Examples:
    | ProductName | Ruleset       | Method    | Event         | Pk       | Desc          | 		Handler_script	   | 		Schema          | Enabled  |
	| drink       | drinkCreated  | "create"  | "drinkCreated"| "id"     | "description" | "./assets/handler.js"   | "./assets/schema.json" | [true]   |
	| drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		  [ignore]         | 		[ignore]        | [ignore] |
    #單獨update method會跳Error: Invalid method
    | drink       | drinkCreated  | create    | [ignore]      | [ignore] | [ignore]      | 		[ignore]         | 		[ignore]      	  | [ignore] |
    | drink       | drinkCreated  | [ignore]  | drinkCreated  | [ignore] | [ignore]      | 		  [ignore]         | 		[ignore]        | [ignore] |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | id       | [ignore]      |		  [ignore]         | 		[ignore]        | [ignore] |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | id,num   | [ignore]      | 		  [ignore]         | 		[ignore]        | [ignore] |
	| drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | 		abc      | 		  [ignore]         | 		[ignore]        | [true]   |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | "description" | 		  [ignore]         | 		[ignore]        | [ignore] |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] |       ""      | 		  [ignore]         | 		[ignore]        | [ignore] |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | ./assets/handler.js     | 		[ignore]        | [ignore] |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		  [ignore]         | ./assets/schema.json   | [ignore] |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		  [ignore]         | 		[ignore]        | [true]   |
#Scenario
	Scenario: 針對更新data product ruleset失敗情境
	Given 已有data product "'drink'"
    Given 已有data product 的 ruleset "'drink'" "'drinkCreated'" 
	Given 儲存nats現有data product ruleset 副本 "'drink'" "'drinkCreated'" 
	When "'<ProductName>'" 更新ruleset "'<Ruleset>'" 參數 method "'<Method>'" event "'<Event>'" pk "'<Pk>'" desc "'<Desc>'" handler "'<Handler_script>'" schema "'<Schema>'" enabled "'<Enabled>'"
	Then ruleset 更改失敗
	And 應有錯誤訊息 "'<Error_message>'"
	And 使用nats驗證 data product "'drink'" 的 ruleset "'drinkCreated'" 資料無任何改動 "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Examples:
    | ProductName | Ruleset       | Method    | Event         | Pk       | Desc          | 		Handler_script 		 | 		   Schema            | Enabled  | Error_message |
	|   [null]    |   [null]      | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		[ignore]	         | 		   [ignore]          | [ignore] |               |
	| drink       |	  [null]      | create    | drinkCreated  | id       | "description" | ./assets/handler.js       | ./assets/schema.json      | [true]   |               |
	| 	[null]    |	drinkCreated  | create    | drinkCreated  | id       | "description" | ./assets/handler.js       | ./assets/schema.json      | [true]   |               |
    | NotExists   | drinkCreated  | create    | drinkCreated  | id       | "description" | ./assets/handler.js       | ./assets/schema.json      | [true]   |               |
	| drink       | NotExists     | create    | drinkCreated  | id       | "description" | ./assets/handler.js       | ./assets/schema.json      | [true]   |               |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		not_exist.js  	     | 		   [ignore]          | [ignore] |               | 
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		[ignore]         	 | 		  not_exist.json     | [ignore] |               |
    | drink       | drinkCreated  | [ignore]  | [ignore]      | [ignore] | [ignore]      | 		[ignore]      	     | ./assets/fail_schena.json | [ignore] |               |