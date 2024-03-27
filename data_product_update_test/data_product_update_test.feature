Feature: Data Product update

Scenario:
Given 已開啟服務nats
Given 已開啟服務dispatcher
#Scenario
	Scenario: 使用者使用product update指令更新data product，成功情境
	Given 已有data product "'<ProductName>'"
	When 更新data product "'<ProductName>'" 註解 "'<Description>'" "'<Enabled>'" schema檔案 "'<Schema>'"
	Then data product更改成功
	And 使用nats驗證data product "'<ProductName>'" description "'<Description>'" schema檔案 "'<Schema>'" enabled "'<Enabled>'" 更改成功
	Examples:
	| ProductName | Description  |		 Schema     	  | Enabled   |
	| drink       | [ignore]     | 		[ignore]   		  | [ignore]  |
	| drink       |     ""       | 		[ignore] 		  | [ignore]  |
	| drink       | description  | 		[ignore]   		  | [ignore]  |
	| drink       | [ignore]     | ./assets/schema.json   | [ignore]  |
	| drink       | [ignore]     | 		[ignore]    	  | [true]    |
	| drink       | "description"| "./assets/schema.json" | [true]    |

#Scenario
	Scenario: 使用者使用product update指令更新data product，失敗情境
	Given 已有data product "'<ProductName>'"
	Given 儲存nats現有data product "'<ProductName>'" 副本
	When 更新data product "'<ProductName>'" 註解 "'<Description>'" "'<Enabled>'" schema檔案 "'<Schema>'"
	Then data product更改失敗
	And 應有錯誤訊息 "'<Error_message>'"
	And 使用nats驗證data product "'<ProductName>'" description "'<Description>'" schema檔案 "'<Schema>'" enabled "'<Enabled>'" 資料無任何改動
	Examples:
	| ProductName   | Description  | 		  Schema           | Enabled   | Error_message   |
	| NotExist      |    ""        |    		 ""            | [ignore]  | 			     |
	| drink		 	| [ignore]     | "./assets/failSchema.json"| [ignore]  |                 |
	| drink		 	| [ignore]     | 		NotExist.json      | [true]    |                 |