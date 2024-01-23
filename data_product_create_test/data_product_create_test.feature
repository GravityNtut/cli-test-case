Feature: Data Product create
	Scenario: 使用者使用product create指令來建立data product
	Given 沒有data product
	Given schema.json=
	"""
	{
        "id": { "type": "uint" },
        "name": { "type": "string" },
        "price": { "type": "uint" },
        "kcal": { "type": "uint" }
	}
	"""
	When container "nats-jetstream" ready
    When container "gravity-dispatcher" ready

	#step 1
	When 創建data product "drink" 註解 "drink data"
	Then data product "drink" 創建成功

	#step 2
	Then 使用gravity-cli 查詢 "drink" 成功

	#step 3
	Then 使用nats jetstream 查詢 "drink" 成功

	Scenario: 測試data product名稱使用特殊字元
	Given 沒有data product
	When 創建data product "-_*($)?@" 註解 "drink data"
	Then data product "-_*($)?@" 創建失敗

	Scenario: 輸入相同名稱的data product
	Given 沒有data product

	#step 1
	When 創建data product "drink" 註解 "drink data"
	Then data product "drink" 創建成功

	#step 2
	Then 使用gravity-cli 查詢 "drink" 成功

	#step 3
	Then 使用nats jetstream 查詢 "drink" 成功

	#step 4
	When 創建data product "drink" 註解 "repeat"
	Then data product "drink" 創建失敗

	Scenario: 測試data product名稱使用中文
	Given 沒有data product
	When 創建data product "飲料" 註解 "飲料相關資料表"
	Then data product "飲料" 創建失敗