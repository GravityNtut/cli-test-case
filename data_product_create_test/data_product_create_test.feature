Feature: Data Product create
#Scenario
	Scenario: 使用者使用product create指令來建立data product
	Given schema.json=
"""
{
	"id": { "type": "uint" },
	"name": { "type": "string" },
	"price": { "type": "uint" },
	"kcal": { "type": "uint" }
}
"""
	#step 1
	When 創建data product "drink" 註解 "drink data" "success"

	#step 2
	Then 使用gravity-cli 查詢 "drink" 成功

	#step 3
	Then 使用nats jetstream 查詢 "drink" 成功

#Scenario
	Scenario: 測試data product名稱使用特殊字元
	When 創建data product "-_*($)?@" 註解 "drink data" "fail"

#Scenario
	Scenario: 輸入相同名稱的data product

	#step 1
	When 創建data product "drink" 註解 "drink data" "success"

	#step 2
	Then 使用gravity-cli 查詢 "drink" 成功

	#step 3
	Then 使用nats jetstream 查詢 "drink" 成功

	#step 4
	When 創建data product "drink" 註解 "repeat" "fail"

#Scenario
	Scenario: 測試data product名稱使用中文
	When 創建data product "飲料" 註解 "飲料相關資料表" "fail"