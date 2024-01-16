Feature: Data Product CRUD
    Scenario: 使用者使用product指令可以管理並使用在jetstream與dispatcher裡product的資料
    When container "nats-jetstream" ready
    When container "gravity-dispatcher" ready
    Given 測試schema "drink_schema" : 
    """
    {
        "id": { "type": "uint" },
        "name": { "type": "string" },
        "price": { "type": "uint" },
        "kcal": { "type": "uint" }
    }
    """

    Given 測試Handle_script "drink_handler_script" :
    """
    return {
        id: source.id,
        name: source.name,
        price: source.price,
        kcal: source.kcal
    }
    """