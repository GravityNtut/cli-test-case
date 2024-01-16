# 單個雙引號是參數，兩個雙引號是字串中的參數，三個雙引號是字串
Feature: Data Product CRUD
    Scenario: 使用者使用product指令可以管理並使用在jetstream與dispatcher裡product的資料
    When container "nats-jetstream" ready
    When container "gravity-dispatcher" ready

    # step 1
    When "nats-jetstream" 沒有Data_Product "drinks"
    Given 輸入cli指令 """gravity-cli product create drink --desc="drink data" --enabled --schema="drink_schema.json" -s nats-jetstream:32803"""
    Then cli回應 """Product "drink" was created"""
    Then "gravity-sdk" 的 "product" 中應該包含 "drink"
    Then "nats-jetstream" 的 "stream" 中應該包含 "GVT_default_DP_drink"

    # step 2
    Given 輸入cli指令 """gravity-cli product ruleset add drink drinkCreated --enabled --event=drinkCreated --method=create --handler="drink_handler_script.js" --schema="drink_schema.json" -s nats-jetstream:32803"""
    Then "gravity-sdk" 的 "drink" 中的 "ruleset" 應該包含 "drinkCreated"

    # step 3
    Given "data" : """{"id":1,"name":"water","price":10,"kcal":0}"""
    Given 輸入cli指令 """gravity-cli pub drinkCreated '""data""' -s nats-jetstream:32803"""
    Then "gravity-sdk" 的 "drink" 訂閱者獲取 "data"
    Then "nats-jetstream" 的 "GVT_default_DP_drink" 中 "Messages" 為"1"筆

    # step 4
    Given 輸入cli指令 """gravity-cli product purge drink -s nats-jetstream:32803"""
    Then "nats-jetstream" 的 "GVT_default_DP_drink" 中 "Messages" 為"0"筆

    #step 5
    Given 輸入cli指令 """gravity-cli product delete drink -s nats-jetstream:32803"""
    Then "gravity-sdk" 的 "product" 中不應該包含 "drink"
    Then "nats-jetstream" 的 "stream" 中不應該包含 "GVT_default_DP_drink"