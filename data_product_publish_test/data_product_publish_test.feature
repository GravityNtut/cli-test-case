Feature: Data Product publish

Scenario:
    Given NATS has been opened
    Given Dispatcher has been opened

# Scenario
    Scenario Outline: Publish for data product of the event (Success scenario)
    Given Create data product "'drink'" using parameters "'<DPEnabled>'"
    Given "'drink'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When Publish Event "'<Event>'" using parameters "'<Payload>'"
    Then Query Jetstream GVT_default message increase
    Then Wait "'1'" second
    When Update data product "'drink'" using parameters enabled=true
    # Then Check data product info "'drink'" is enabled
    Then Wait "'1'" second
    When Update data product "'drink'" ruleset "'drinkEvent'" using parameters enabled=true
    # When Publish Event "'<Event>'" using parameters "''{"id":222,"uid":1,"name":"test","price":100,"kcal":50}''"
    Then Query GVT_default_DP_"'drink'" has a event with payload "'<Payload>'" and type is "'<Event>'" and its match with "'<RSSchema>'"
    Examples:
        |   ID   | Event      |                    Payload                             | RSMethod |       RSHandler         |       RSSchema                    | RSPk     | RSEnabled | DPEnabled |
        |  M(1)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(2)  | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(3)  | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(4)  | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(5)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(6)  | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(7)  | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(8)  | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(9)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(10) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(11) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(12) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [true]    |
        |  M(13) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(14) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(15) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(16) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [false]   | [true]    |
        |  M(17) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(18) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(19) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(20) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [false]   |
        |  M(21) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(22) | drinkEvent | '{"id":1,"uid":1,"name":"test"}'                       | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(23) | drinkEvent | '{"id":1,"uid":1,"name":"test","not_exist":"Hi"}'      | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(24) | drinkEvent | '{"id":123,"uid":1,"id":1213}'                         | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [false]   |
        |  M(25) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | ""       | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(26) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | " "      | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(27) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | _-*($)?@ | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(28) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(29) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/WrongType.js   | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(30) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/ExtraColumn.js | ./assets/schema.json              | id       | [true]    | [true]    |
        |  M(31) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   |      [ignore]           |          ./assets/schema.json     |    id    | [true]    | [true]    |
        |  M(32) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schemaWrongType.json     | id       | [true]    | [true]    |
        |  M(33) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,uid   | [true]    | [true]    |
        |  M(34) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id,      | [true]    | [true]    |
        |  M(35) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schema.json              | id       | [false]   | [false]   |

    Scenario Outline: Similar to the successful scenario, but after publishing, update the ruleset to enabled first, then update the data product to enabled.
    Given Create data product "'drink'" using parameters "'<DPEnabled>'"
    Given "'drink'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When Publish Event "'<Event>'" using parameters "'<Payload>'"
    When Update data product "'drink'" ruleset "'drinkEvent'" using parameters enabled=true
    When Update data product "'drink'" using parameters enabled=true
    Then Query GVT_default_DP_"'drink'" has a event with payload "'<Payload>'" and type is "'<Event>'" and its match with "'<RSSchema>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  E1(1) | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [false]   | [false]   |
    
    Scenario Outline: publish for data product of the event (failure scenario for the publish command).
    Given Create data product "'drink'" using parameters "'<DPEnabled>'"
    Given "'drink'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When Publish Event "'<Event>'" using parameters "'<Payload>'"
    Then CLI returns create failed
    # And The error message should be "'<Error_message>'"
    Examples:
        |   ID   | Event      |                    Payload                     | RSMethod |       RSHandler     |       RSSchema       | RSPk | RSEnabled | DPEnabled |
        |  E2(1) | ""         | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        |  E2(2) | drinkEvent | ''                                             | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        |  E2(3) | drinkEvent | ' '                                            | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |
        |  E2(4) | drinkEvent | 'abc'                                          | create   | ./assets/handler.js | ./assets/schema.json | id   | [true]    | [true]    |

    Scenario Outline: The command executes successfully but does not publish to the specified DP.
    Given Create data product "'drink'" using parameters "'<DPEnabled>'"
    Given "'drink'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When Publish Event "'<Event>'" using parameters "'<Payload>'"
    Then Query GVT_default_DP_"'drink'" has no "'<Event>'"
    Then Using NATS Jetstream to query GVT_default has a event with payload "'<Payload>'" and type is "'<Event>'"
    Examples:
        |    ID   |    Event     |                    Payload                     | RSMethod |         RSHandler                                  |               RSSchema              |    RSPk    | RSEnabled |    DPEnabled    |
        |  E3(1)  | drinkEvent   |               '{"name":"test"}'                |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(2)  | drinkEvent   |         '{"name":"test","not_exist":"Hi"}'     |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(3)  | drinkEvent   |             '{"not_exist":"Hi"}'               |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(4)  | drinkEvent   |                     '{}'                       |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(5)  | drinkEvent   |               '{"name":"test"}'                |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(6)  | drinkEvent   |         '{"name":"test","not_exist":"Hi"}'     |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(7)  | drinkEvent   |             '{"not_exist":"Hi"}'               |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(8)  | drinkEvent   |                     '{}'                       |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id,uid  |  [true]   |      [true]     |
        |  E3(9)  | NotExist     | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        # |  E3(10) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/empty.js                             |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        # E3(10) will cause dispatcher crash
        |  E3(11) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/fail_handler.js                      |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(12) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/WrongFileExtension.txt               |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(13) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/WrongFileExtensionAndFormat.jpg      |          ./assets/schema.json       |    id      |  [true]   |      [true]     |
        |  E3(14) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    ""      |  [true]   |      [true]     |
        |  E3(15) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    " "     |  [true]   |      [true]     |
        |  E3(16) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |  _-*($)?@  |  [true]   |      [true]     |
        |  E3(17) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [false]  |      [true]     |
        |  E3(18) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      ./assets/handler.js                           |          ./assets/schema.json       |    id      |  [true]   |      [false]    |
        |  E3(19) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schemaWithNoPK.json      | id       | [true]    | [true]    |
        |  E3(20) | drinkEvent   | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | create   | ./assets/handler.js     | ./assets/schemaWithEmptyJson.json | id       | [true]    | [true]    |

    Scenario Outline: publish for data product of the event (The same event is published to multiple data products)
    Given Create data product "'drink'" using parameters "'<DPEnabled>'"
    Given Create data product "'drink2'" using parameters "'<DPEnabled>'"
    Given "'drink'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    Given "'drink2'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When Publish Event "'<Event>'" using parameters "'<Payload>'"
    Then Query GVT_default_DP_"'drink'" has a event with payload "'<Payload>'" and type is "'<Event>'" and its match with "'<RSSchema>'"
    Then Query GVT_default_DP_"'drink2'" has a event with payload "'<Payload>'" and type is "'<Event>'" and its match with "'<RSSchema>'"
    Examples:
        |    ID   |    Event   |                    Payload                     | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E4(1)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |

    Scenario Outline: publish for data product of the event (Continuously publish two events with the same PK value, but the number and content of other fields are different.)
    Given Create data product "'drink'" using parameters "'<DPEnabled>'"
    Given "'drink'" create ruleset "'drinkEvent'" using parameters "'<RSMethod>'" "'<RSPk>'" "'<RSHandler>'" "'<RSSchema>'" "'<RSEnabled>'"
    When Publish Event "'<Event>'" using parameters "'<Payload>'"
    When Publish Event "'<Event>'" using parameters "'<Payload2>'"
    Then Query GVT_default_DP_"'drink'" has two events with payload "'<Payload>'" and "'<Payload2>'" and type is "'<Event>'"
    Examples:
        |    ID   |    Event   |                    Payload                             |         Payload2                 | RSMethod |         RSHandler           |               RSSchema            |    RSPk    | RSEnabled |    DPEnabled    |
        |  E5(1)  | drinkEvent | '{"id":1,"uid":1,"name":"test","price":100,"kcal":50}' | '{"id":1,"uid":222,"name":"replace"}' |  create  |      assets/handler.js      |          assets/schema.json       |    id      |  [true]   |      [true]     |